package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func getKubeconfig() (string, error) {
	// TODO: Allow customizing kubeconfig path
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".kube", "config"), nil
	} else {
		return "", err
	}
}

func newClientset() (*kubernetes.Clientset, error) {
	kubeconfig, err := getKubeconfig()
	if err != nil {
		return nil, err
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func GetCurrentContext() (string, error) {
	kubeconfig, err := getKubeconfig()
	if err != nil {
		return "", err
	}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		nil,
	).RawConfig()
	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}
