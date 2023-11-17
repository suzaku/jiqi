package k8s

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"strings"
	"time"
)

type Node struct {
	Name           string            `json:"name"`
	Labels         map[string]string `json:"labels"`
	ConsolePageURL string            `json:"consolePageURL"`
	DashboardURL   string            `json:"dashboardURL"`
	InstanceType   string            `json:"instanceType"`
	Usage          NodeUsage         `json:"usage"`
	Capacity       NodeCapacity      `json:"capacity"`
}

type NodeUsage struct {
	Cpu    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

type NodeCapacity struct {
	Cpu    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
}

func (nm NodesManager) getNodesMetric(ctx context.Context) (map[string]NodeUsage, error) {
	config, err := clientcmd.BuildConfigFromFlags("", nm.kubeconfig)
	if err != nil {
		return nil, err
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	nodeMetricsList, err := mc.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	usages := make(map[string]NodeUsage, len(nodeMetricsList.Items))
	for _, metric := range nodeMetricsList.Items {
		usages[metric.Name] = NodeUsage{
			Cpu:    metric.Usage.Cpu().MilliValue(),
			Memory: metric.Usage.Memory().Value(),
		}
	}
	return usages, nil
}

func (nm NodesManager) GetNodes(ctx context.Context, shouldClearCache bool) ([]Node, error) {
	currentContext, err := nm.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	if cache, ok := nm.nodesCache[currentContext]; ok && !shouldClearCache {
		return cache.nodes, nil
	}

	clientset, err := nm.NewClientset()
	if err != nil {
		return nil, err
	}

	chUsageResult := make(chan map[string]NodeUsage)
	go func() {
		defer close(chUsageResult)
		result, err := nm.getNodesMetric(ctx)
		if err != nil {
			panic(err)
		}
		chUsageResult <- result
	}()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var usages map[string]NodeUsage
	var ok bool
	select {
	case usages, ok = <-chUsageResult:
		if !ok {
			usages = make(map[string]NodeUsage)
		}
	}

	myNodes := make([]Node, len(nodes.Items))
	for i, node := range nodes.Items {
		// TODO: Handle error
		consolePageURL, err := getConsolePageURL(node)
		if err != nil {
			panic(err)
		}
		dashboardURL, err := getDashboardURL(node)
		if err != nil {
			panic(err)
		}
		instanceType, ok := getInstanceType(node)
		if !ok {
			instanceType = "unknown"
		}
		capacity := NodeCapacity{
			Cpu:    node.Status.Capacity.Cpu().MilliValue(),
			Memory: node.Status.Capacity.Memory().Value(),
		}
		myNodes[i] = Node{
			Name:           node.Name,
			Labels:         node.Labels,
			ConsolePageURL: consolePageURL,
			DashboardURL:   dashboardURL,
			InstanceType:   instanceType,
			Usage:          usages[node.Name],
			Capacity:       capacity,
		}
	}

	nm.nodesCache[currentContext] = CachedNodes{
		createdAt: time.Now(),
		nodes:     myNodes,
	}

	return myNodes, nil
}

func getInstanceType(node v1.Node) (string, bool) {
	instanceType, ok := node.Labels["node.kubernetes.io/instance-type"]
	return instanceType, ok
}

func getDashboardURL(node v1.Node) (string, error) {
	var instanceId string
	var err error
	if instanceId, err = parseInstanceId(node); err != nil {
		return "", err
	}
	return "https://app.datadoghq.com/dash/host_name/" + instanceId, nil
}

func getConsolePageURL(node v1.Node) (string, error) {
	var instanceId string
	var err error
	if instanceId, err = parseInstanceId(node); err != nil {
		return "", err
	}
	var region string
	if region, err = getRegion(node); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"https://%s.console.aws.amazon.com/ec2/home?region=%s#InstanceDetails:instanceId=%s",
		region,
		region,
		instanceId,
	), nil
}

func parseInstanceId(node v1.Node) (string, error) {
	providerId := node.Spec.ProviderID
	parts := strings.Split(providerId, "/")
	if len(parts) < 2 {
		return "", errors.New("invalid provider id: " + providerId)
	}
	return parts[len(parts)-1], nil
}

func getRegion(node v1.Node) (string, error) {
	if region, ok := node.ObjectMeta.Labels["topology.kubernetes.io/region"]; ok {
		return region, nil
	}
	return "", errors.New("no region info found")
}
