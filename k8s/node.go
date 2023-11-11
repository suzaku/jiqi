package k8s

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Node struct {
	Name           string `json:"name"`
	ConsolePageURL string `json:"consolePageURL"`
	DashboardURL   string `json:"dashboardURL"`
}

func GetNodes() ([]Node, error) {
	clientset, err := newClientset()
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
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
		myNodes[i] = Node{
			Name:           node.Name,
			ConsolePageURL: consolePageURL,
			DashboardURL:   dashboardURL,
		}
	}
	return myNodes, nil
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
