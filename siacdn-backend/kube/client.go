package kube

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterType string

var ClusterTypeInvalid ClusterType = ClusterType("invalid")
var ClusterTypeInternal ClusterType = ClusterType("internal")
var ClusterTypeExternal ClusterType = ClusterType("external")

func KubeClient() (*kubernetes.Clientset, ClusterType, error) {
	home := homeDir()
	clusterType := ClusterTypeInvalid
	kubeConfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err == nil {
		clusterType = ClusterTypeExternal
	} else {
		log.Println("Problem with default kube config: " + err.Error())
		log.Println("Trying in-cluster config...")
		config, err = rest.InClusterConfig()
		if err == nil {
			clusterType = ClusterTypeInternal
		} else {
			log.Println("Could not get kube cluster config to work: " + err.Error())
			return nil, clusterType, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Could not create clientset from config: " + err.Error())
		return nil, ClusterTypeInvalid, err
	}
	return clientset, clusterType, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
