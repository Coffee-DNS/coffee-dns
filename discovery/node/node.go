package node

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Nodes returns a list of kubernetes nodes
func Nodes() ([]v1.Node, error) {
	clientset, err := getClient("")
	if err != nil {
		return nil, err
	}

	nodes, err := nodes(clientset)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func nodes(clientset kubernetes.Interface) ([]v1.Node, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
}

func getClient(pathToCfg string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if pathToCfg == "" {
		fmt.Println("Using in cluster config")
		config, err = rest.InClusterConfig()
		// in cluster access
	} else {
		fmt.Println("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", pathToCfg)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
