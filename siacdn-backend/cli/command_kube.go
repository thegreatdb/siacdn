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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeNamespace = "sia"

var kubeDefaultStorage = resource.MustParse("100Gi")

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
		log.Println("Problem with default kube config: " + err.Error())
		log.Println("Trying in-cluster config...")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Println("Could not get kube cluster config to work: " + err.Error())
			return err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Could not create clientset from config: " + err.Error())
		return err
	}
	for {
		if err = PerformKubeRun(clientset); err != nil {
			return err
		}
		time.Sleep(time.Second / 2)
	}
	return nil
}

func getPendingSiaNodes() ([]*models.SiaNode, error) {
	url := fmt.Sprintf("%s/sianodes/pending/all?secret=%s", URLRoot, SiaCDNSecretKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Could not create request GET " + url)
		return nil, err
	}

	res, err := cliClient.Do(req)
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

func PerformKubeRun(clientset *kubernetes.Clientset) error {
	siaNodes, err := getPendingSiaNodes()
	if err != nil {
		return err
	}
	for _, siaNode := range siaNodes {
		if err = pollKube(clientset, siaNode); err != nil {
			return err
		}
	}
	return nil
}

func pollKube(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	switch siaNode.Status {
	case models.SIANODE_STATUS_CREATED:
		return pollKubeCreated(clientset, siaNode)
	case models.SIANODE_STATUS_DEPLOYED:
		return pollKubeDeployed(clientset, siaNode)
	case models.SIANODE_STATUS_INSTANCED:
		return pollKubeInstanced(clientset, siaNode)
	case models.SIANODE_STATUS_SNAPSHOTTED:
		return pollKubeSnapshotted(clientset, siaNode)
	case models.SIANODE_STATUS_SYNCHRONIZED:
		return pollKubeSynchronized(clientset, siaNode)
	case models.SIANODE_STATUS_INITIALIZED:
		return pollKubeInitialized(clientset, siaNode)
	case models.SIANODE_STATUS_FUNDED:
		return pollKubeFunded(clientset, siaNode)
	case models.SIANODE_STATUS_CONFIGURED:
		return pollKubeConfigured(clientset, siaNode)
	default:
		log.Println("Unknown status: " + siaNode.Status)
	}
	return nil
}

func pollKubeCreated(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeCreated: " + siaNode.Shortcode)

	volumeClaims := clientset.PersistentVolumeClaims(kubeNamespace)
	deployments := clientset.AppsV1beta1Client.Deployments(kubeNamespace)
	services := clientset.Services(kubeNamespace)

	// First check for volume claim
	claim, err := volumeClaims.Get(siaNode.KubeNameVol(), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting volume claim from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if claim == nil || errors.IsNotFound(err) {
		spec := v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName("storage"): kubeDefaultStorage,
				},
			},
		}
		storageClass := "standard"
		spec.StorageClassName = &storageClass

		claim = &v1.PersistentVolumeClaim{}
		claim.Name = siaNode.KubeNameVol()
		claim.Namespace = kubeNamespace
		claim.Spec = spec

		log.Println("Creating volume claim " + siaNode.KubeNameVol())
		claim, err = volumeClaims.Create(claim)
		if err != nil {
			log.Println("Error creating volume claim: " + err.Error())
			return err
		}
	} else {
		log.Println("Found volume claim " + siaNode.KubeNameVol())
	}

	// Then check for service
	service, err := services.Get(siaNode.KubeNameSer(), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting service from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if service == nil || errors.IsNotFound(err) {
		service = &v1.Service{}
		service.Name = siaNode.KubeNameSer()
		service.Namespace = kubeNamespace
		service.Spec = v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
				v1.ServicePort{Name: "p1", Port: 9980, TargetPort: intstr.FromInt(9980), Protocol: v1.ProtocolTCP},
				v1.ServicePort{Name: "p2", Port: 9980, TargetPort: intstr.FromInt(9980), Protocol: v1.ProtocolUDP},
				v1.ServicePort{Name: "p3", Port: 9981, TargetPort: intstr.FromInt(9981), Protocol: v1.ProtocolTCP},
				v1.ServicePort{Name: "p4", Port: 9981, TargetPort: intstr.FromInt(9981), Protocol: v1.ProtocolUDP},
				v1.ServicePort{Name: "p5", Port: 9982, TargetPort: intstr.FromInt(9982), Protocol: v1.ProtocolTCP},
				v1.ServicePort{Name: "p6", Port: 9982, TargetPort: intstr.FromInt(9982), Protocol: v1.ProtocolUDP},
			},
			Selector: map[string]string{"app": siaNode.KubeNameApp()},
		}
		log.Println("Creating service " + siaNode.KubeNameSer())
		service, err = services.Create(service)
		if err != nil {
			log.Println("Error creating service: " + err.Error())
			return err
		}
	} else {
		log.Println("Found service " + siaNode.KubeNameSer())
	}

	// Finally, check for deployment
	deployment, err := deployments.Get(siaNode.KubeNameDep(), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting deployment from kubernetes: " + err.Error())
		return err
	}
	// If deployment doesn't exist, create it
	if deployment == nil || errors.IsNotFound(err) {
		deployment := &v1beta1.Deployment{}
		deployment.Name = siaNode.KubeNameDep()
		deployment.Namespace = kubeNamespace
		deployment.Spec = v1beta1.DeploymentSpec{Template: v1.PodTemplateSpec{}}
		deployment.Spec.Template.Labels = map[string]string{"app": siaNode.KubeNameApp()}
		deployment.Spec.Template.Spec = v1.PodSpec{
			Volumes: []v1.Volume{
				v1.Volume{
					Name: "sia-volume",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: siaNode.KubeNameVol(),
						},
					},
				},
			},
			Containers: []v1.Container{
				v1.Container{
					Name:            "sia-docker",
					Image:           "gcr.io/gradientzoo-1233/sia-docker:latest",
					ImagePullPolicy: v1.PullAlways,
					Ports: []v1.ContainerPort{
						v1.ContainerPort{ContainerPort: 9980},
						v1.ContainerPort{ContainerPort: 9981},
						v1.ContainerPort{ContainerPort: 9982},
					},
					VolumeMounts: []v1.VolumeMount{
						v1.VolumeMount{Name: "sia-volume", MountPath: "/sia"},
					},
				},
			},
		}
		log.Println("Creating deployment " + siaNode.KubeNameDep())
		deployment, err = deployments.Create(deployment)
		if err != nil {
			log.Println("Error creating deployment: " + err.Error())
			return err
		}
	} else {
		log.Println("Found deployment " + siaNode.KubeNameDep())
	}

	log.Println("TODO: Change state from created to deployed")
	// TODO: Change the state of the siaNode to deployed

	return nil
}

func pollKubeDeployed(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeDeployed: " + siaNode.Shortcode)
	return nil
}

func pollKubeInstanced(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeInstanced: " + siaNode.Shortcode)
	return nil
}

func pollKubeSnapshotted(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeSnapshotted: " + siaNode.Shortcode)
	return nil
}

func pollKubeSynchronized(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeSynchronized: " + siaNode.Shortcode)
	return nil
}

func pollKubeInitialized(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeInitialized: " + siaNode.Shortcode)
	return nil
}

func pollKubeFunded(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeFunded: " + siaNode.Shortcode)
	return nil
}

func pollKubeConfigured(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeConfigured: " + siaNode.Shortcode)
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
