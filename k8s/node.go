package k8s

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

type Node struct {
	Name           string `json:"name"`
	ConsolePageURL string `json:"consolePageURL"`
	DashboardURL   string `json:"dashboardURL"`
	InstanceType   string `json:"instanceType"`
}

func (nm NodesManager) GetNodes(ctx context.Context) ([]Node, error) {
	currentContext, err := nm.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	var nodes *v1.NodeList
	if cache, ok := nm.nodesCache[currentContext]; !ok {
		clientset, err := nm.NewClientset()
		if err != nil {
			return nil, err
		}

		nodes, err = clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		nm.nodesCache[currentContext] = snapshot{
			createdAt: time.Now(),
			nodes:     nodes,
		}
	} else {
		nodes = cache.nodes
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
		myNodes[i] = Node{
			Name:           node.Name,
			ConsolePageURL: consolePageURL,
			DashboardURL:   dashboardURL,
			InstanceType:   instanceType,
		}
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
