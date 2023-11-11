package k8s

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

type snapshot struct {
	createdAt time.Time
	nodes     *v1.NodeList
}

type NodesManager struct {
	kubeconfig string
	nodesCache map[string]snapshot
}

func NewNodesManager() (NodesManager, error) {
	kubeconfig, err := getKubeconfig()
	if err != nil {
		return NodesManager{}, err
	}
	return NodesManager{
		kubeconfig: kubeconfig,
		nodesCache: make(map[string]snapshot),
	}, nil
}

func (nm NodesManager) GetCurrentContext() (string, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: nm.kubeconfig},
		nil,
	).RawConfig()
	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}

func (nm NodesManager) NewClientset() (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", nm.kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func getKubeconfig() (string, error) {
	// TODO: Allow customizing kubeconfig path
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".kube", "config"), nil
	} else {
		return "", err
	}
}
