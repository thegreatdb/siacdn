package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/thegreatdb/siacdn/siacdn-backend/models"
	urfavecli "github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func KubeCommand() urfavecli.Command {
	return urfavecli.Command{
		Name:    "kube",
		Aliases: []string{"k"},
		Usage:   "Communicate with a local SiaCDN backend and coordinate changes with a kube server",
		Action:  kubeCommand,
	}
}

func kubeCommand(c *urfavecli.Context) error {
	home := homeDir()
	kubeConfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Pod not found\n")
		}
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for {
		if err = PerformKubeRun(); err != nil {
			return err
		}
		time.Sleep(time.Second / 2)
	}
	return nil
}

func GetPendingSiaNodes() ([]*models.SiaNode, error) {
	url := "http://localhost:9095/sianodes/pending/all?secret=" + SiaCDNSecretKey

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Could not create request GET " + url)
		return nil, err
	}

	client := http.Client{Timeout: time.Second * 6}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making GetPendingSiaNodes request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read GetPendingSiaNodes response: " + err.Error())
		return nil, err
	}

	var nodes struct {
		SiaNodes []*models.SiaNode `json:"sianodes"`
	}
	if err = json.Unmarshal(body, &nodes); err != nil {
		log.Println("Could not decode response: " + err.Error())
		return nil, err
	}

	return nodes.SiaNodes, nil
}

func PerformKubeRun() error {
	siaNodes, err := GetPendingSiaNodes()
	if err != nil {
		return err
	}
	for _, siaNode := range siaNodes {
		if err = pollKube(siaNode); err != nil {
			return err
		}
	}
	return nil
}

func pollKube(siaNode *models.SiaNode) error {
	switch siaNode.Status {
	case models.SIANODE_STATUS_CREATED:
		return pollKubeCreated(siaNode)
	case models.SIANODE_STATUS_DEPLOYED:
		return pollKubeDeployed(siaNode)
	case models.SIANODE_STATUS_INSTANCED:
		return pollKubeInstanced(siaNode)
	case models.SIANODE_STATUS_SNAPSHOTTED:
		return pollKubeSnapshotted(siaNode)
	case models.SIANODE_STATUS_SYNCHRONIZED:
		return pollKubeSynchronized(siaNode)
	case models.SIANODE_STATUS_INITIALIZED:
		return pollKubeInitialized(siaNode)
	case models.SIANODE_STATUS_FUNDED:
		return pollKubeFunded(siaNode)
	case models.SIANODE_STATUS_CONFIGURED:
		return pollKubeConfigured(siaNode)
	default:
		log.Println("Unknown status: " + siaNode.Status)
	}
	return nil
}

func pollKubeCreated(siaNode *models.SiaNode) error {
	log.Println("PollKubeCreated: " + siaNode.Shortcode)
	return nil
}

func pollKubeDeployed(siaNode *models.SiaNode) error {
	log.Println("PollKubeDeployed: " + siaNode.Shortcode)
	return nil
}

func pollKubeInstanced(siaNode *models.SiaNode) error {
	log.Println("PollKubeInstanced: " + siaNode.Shortcode)
	return nil
}

func pollKubeSnapshotted(siaNode *models.SiaNode) error {
	log.Println("PollKubeSnapshotted: " + siaNode.Shortcode)
	return nil
}

func pollKubeSynchronized(siaNode *models.SiaNode) error {
	log.Println("PollKubeSynchronized: " + siaNode.Shortcode)
	return nil
}

func pollKubeInitialized(siaNode *models.SiaNode) error {
	log.Println("PollKubeInitialized: " + siaNode.Shortcode)
	return nil
}

func pollKubeFunded(siaNode *models.SiaNode) error {
	log.Println("PollKubeFunded: " + siaNode.Shortcode)
	return nil
}

func pollKubeConfigured(siaNode *models.SiaNode) error {
	log.Println("PollKubeConfigured: " + siaNode.Shortcode)
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
