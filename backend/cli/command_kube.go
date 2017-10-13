package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/NebulousLabs/Sia/api"
	"github.com/NebulousLabs/Sia/types"
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/backend/kube"
	"github.com/thegreatdb/siacdn/backend/models"
	"github.com/thegreatdb/siacdn/backend/prime"
	urfavecli "github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const kubeNamespace = "sia"

//var kubeStorageClass = "standard"
var kubeStorageClass = "fast"
var kubeDefaultStorage = resource.MustParse("30Gi")
var minioStorageClass = "standard"
var nfsStorageClass = ""
var minioDefaultStorage = resource.MustParse("100Gi")

const siaHosts = 40
const siaNeededContracts = 30  // Once we have enough hosts, we can confirm it
const siaContractPeriod = 4380 // Number of 10m intervals in 1 month
const siaRenewWindow = 400
const siaFlightPrefix = "sia-"
const siaMinerFees = 12
const minioFlightPrefix = "minio-"

var securityContextPrivileged bool = true

var kubeMu sync.Mutex
var kubeInFlight map[string]bool = map[string]bool{}

func StartFlight(prefix string, siaNode *models.SiaNode) {
	kubeMu.Lock()
	kubeInFlight[prefix+siaNode.ID.String()] = true
	kubeMu.Unlock()
}

func StopFlight(prefix string, siaNode *models.SiaNode) {
	kubeMu.Lock()
	delete(kubeInFlight, prefix+siaNode.ID.String())
	kubeMu.Unlock()
}

func KubeCommand() urfavecli.Command {
	return urfavecli.Command{
		Name:    "kube",
		Aliases: []string{"k"},
		Usage:   "Communicate with a local SiaCDN backend and coordinate changes with a kube server",
		Action:  kubeCommand,
	}
}

func kubeCommand(c *urfavecli.Context) error {
	clientset, _, err := kube.KubeClient()
	for {
		if err = PerformKubeRun(clientset); err != nil {
			return err
		}
		time.Sleep(time.Second / 2)
	}
	return nil
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
	siaNodes, err = getReadyOrphanedSiaNodes()
	if err != nil {
		return err
	}
	for _, siaNode := range siaNodes {
		for i := 0; i < siaNode.MinioInstancesRequested; i++ {
			if err = deployMinio(clientset, siaNode, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func pollKube(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	kubeMu.Lock()
	_, inFlight := kubeInFlight[siaFlightPrefix+siaNode.ID.String()]
	kubeMu.Unlock()
	if inFlight {
		log.Println("Skipping " + siaNode.ID.String() + " (sia) because it is in-flight")
		return nil
	}
	// TODO: Also check length and skip if kubeInFlight has more than 10 going too
	switch siaNode.Status {
	case models.SIANODE_STATUS_CREATED:
		go pollKubeCreated(clientset, siaNode)
	case models.SIANODE_STATUS_DEPLOYED:
		go pollKubeDeployed(clientset, siaNode)
	case models.SIANODE_STATUS_INSTANCED:
		go pollKubeInstanced(clientset, siaNode)
	case models.SIANODE_STATUS_SNAPSHOTTED:
		go pollKubeSnapshotted(clientset, siaNode)
	case models.SIANODE_STATUS_SYNCHRONIZED:
		go pollKubeSynchronized(clientset, siaNode)
	case models.SIANODE_STATUS_INITIALIZED:
		go pollKubeInitialized(clientset, siaNode)
	case models.SIANODE_STATUS_UNLOCKED:
		go pollKubeUnlocked(clientset, siaNode)
	case models.SIANODE_STATUS_FUNDED:
		go pollKubeFunded(clientset, siaNode)
	case models.SIANODE_STATUS_CONFIRMED:
		go pollKubeConfirmed(clientset, siaNode)
	case models.SIANODE_STATUS_CONFIGURED:
		go pollKubeConfigured(clientset, siaNode)
	case models.SIANODE_STATUS_STOPPING:
		go pollKubeStopping(clientset, siaNode)
	default:
		log.Println("Unknown status: " + siaNode.Status)
	}
	return nil
}

func pollKubeCreated(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeCreated: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	volumeClaims := clientset.PersistentVolumeClaims(kubeNamespace)
	deployments := clientset.AppsV1beta1Client.Deployments(kubeNamespace)
	services := clientset.CoreV1Client.Services(kubeNamespace)
	secrets := clientset.Secrets(kubeNamespace)

	// First check for volume claim
	claim, err := volumeClaims.Get(siaNode.KubeNameVol(), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting volume claim from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if claim == nil || errors.IsNotFound(err) {
		spec := v1.PersistentVolumeClaimSpec{
			StorageClassName: &kubeStorageClass,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName("storage"): kubeDefaultStorage,
				},
			},
		}

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

	// Check on the off chance they already have a seed chosen
	siaWalletPassword := []byte{}
	if s, err := getWalletSeed(siaNode.ID); err == nil && s != nil && s.Words != "" {
		siaWalletPassword = []byte(s.Words)
	}

	// Third, check for secret
	secret, err := secrets.Get(siaNode.KubeNameSec(), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting secret from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if service == nil || errors.IsNotFound(err) {
		secret = &v1.Secret{}
		secret.Name = siaNode.KubeNameSec()
		secret.Namespace = kubeNamespace
		secret.Type = v1.SecretTypeOpaque
		secret.Data = map[string][]byte{
			"siawalletpassword": siaWalletPassword,
		}
		log.Println("Creating secret " + siaNode.KubeNameSec())
		secret, err = secrets.Create(secret)
		if err != nil {
			log.Println("Error creating secret: " + err.Error())
			return err
		}
	} else {
		log.Println("Found secret " + siaNode.KubeNameSec())
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
		deployment.Spec.Strategy.Type = v1beta1.RecreateDeploymentStrategyType
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
					Env: []v1.EnvVar{
						v1.EnvVar{
							Name: "SIA_API_PASSWORD",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "sia-secret",
									},
									Key: "siaapipassword",
								},
							},
						},
						v1.EnvVar{
							Name: "SIA_WALLET_PASSWORD",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: siaNode.KubeNameSec(),
									},
									Key: "siawalletpassword",
								},
							},
						},
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

	_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_DEPLOYED)
	return err
}

func pollKubeDeployed(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeDeployed: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)
	pod, _ := getPod(clientset, siaNode)
	if pod == nil || pod.Name == "" {
		return nil
	}
	log.Println("Found pod " + siaNode.KubeNameDep())
	_, err := updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_INSTANCED)
	return err
}

func pollKubeInstanced(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeInstanced: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get Sia client: " + err.Error())
		return err
	}

	var resp api.ConsensusGET
	if err = client.Get("/consensus", &resp); err != nil {
		log.Println("Got error checking for instance: " + err.Error())
		return nil
	}

	log.Println("Finished snapshotting Sia node " + siaNode.Shortcode)

	_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_SNAPSHOTTED)

	return err
}

func pollKubeSnapshotted(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeSnapshotted: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get Sia client: " + err.Error())
		return err
	}

	var resp api.ConsensusGET
	if err = client.Get("/consensus", &resp); err != nil {
		log.Println("Got error checking for consensus: " + err.Error())
		return nil
	}

	if resp.Synced {
		log.Println("Finished syncing blockchain on Sia node " + siaNode.Shortcode)
		_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_SYNCHRONIZED)
		return err
	}

	return nil
}

func pollKubeSynchronized(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeSynchronized: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	pods := clientset.Pods(kubeNamespace)
	secrets := clientset.Secrets(kubeNamespace)

	client, err := siaNode.SiaClient()
	if err != nil {
		return err
	}

	seed, err := getWalletSeed(siaNode.ID)
	if err != nil {
		log.Println("Got error checking wallet seed: " + err.Error())
		return err
	}

	// If for whatever reason we found one, we can init with seed
	if seed != nil && seed.Words != "" {
		var resp map[string]interface{}
		vals := url.Values{}
		//vals.Set("force", "true")
		vals.Set("seed", seed.Words)
		if err = client.Post("/wallet/init/seed", vals.Encode(), &resp); err != nil {
			log.Println("Got error initializing wallet: " + err.Error())
			return err
		}
	} else {
		var resp api.WalletInitPOST
		vals := url.Values{}
		//vals.Set("force", "true")
		if err = client.Post("/wallet/init", vals.Encode(), &resp); err != nil {
			log.Println("Got error initializing wallet: " + err.Error())
			return err
		}

		if resp.PrimarySeed == "" {
			log.Println("Could not initialize wallet")
			return fmt.Errorf("Could not initialize wallet")
		}

		seed, err = createWalletSeed(siaNode.ID, resp.PrimarySeed)
		if err != nil {
			log.Println("Got error saving wallet seed: " + err.Error())
			return err
		}
	}

	log.Println("Getting secret " + siaNode.KubeNameSec())
	secret, err := secrets.Get(siaNode.KubeNameSec(), metav1.GetOptions{})
	if err != nil {
		log.Println("Error getting secret from kubernetes: " + err.Error())
		return err
	}

	if bytes.Equal(secret.Data["siawalletpassword"], []byte(seed.Words)) {
		log.Println("Found that secret was already set for " + siaNode.Shortcode)
	} else {
		secret.Data["siawalletpassword"] = []byte(seed.Words)

		log.Println("Updating kubernetes secret with new info")
		if secret, err = secrets.Update(secret); err != nil {
			log.Println("Error updating secret from kubernetes: " + err.Error())
			return err
		}

		log.Println("Closing running containers so they receive the new secret")
		err = pods.DeleteCollection(getDeleteOpts(), metav1.ListOptions{
			LabelSelector: "app=" + siaNode.KubeNameApp(),
		})
		if err != nil {
			log.Println("Could not close running containers")
			return err
		}
	}

	log.Println("Initialized wallet on Sia node " + siaNode.Shortcode)
	_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_INITIALIZED)
	return err
}

func pollKubeInitialized(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeInitialized: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	client, err := siaNode.SiaClient()
	if err != nil {
		return err
	}

	log.Println("Checking wallet unlock " + siaNode.Shortcode)
	var resp api.WalletGET
	if err = client.Get("/wallet", &resp); err != nil {
		log.Println("Got error checking wallet unlock: " + err.Error())
		return nil
	}

	if resp.Unlocked {
		_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_UNLOCKED)
		return err
	} else {
		log.Println("Wallet found but not yet unlocked: " + siaNode.Shortcode)
	}
	return nil
}

func pollKubeUnlocked(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeUnlocked: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	prime, err := prime.Server(clientset)
	if err != nil {
		log.Println("Could not get a connection to the prime Sia node: " + err.Error())
		return err
	}

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get a connection to the sia node" + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	// First check to see if the wallet has enough already. If it does, skip to confirmed.
	var curResp api.WalletGET
	if err = client.Get("/wallet", &curResp); err != nil {
		log.Println("Could not get balance for " + siaNode.Shortcode + ": " + err.Error())
		return err
	}
	if curResp.ConfirmedSiacoinBalance.Cmp(siaNode.DesiredCurrency()) >= 0 {
		_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_CONFIRMED)
		if err != nil {
			log.Println("Could not update the SiaNode status to confirmed: " + err.Error())
			return err
		}
	}

	var address api.WalletAddressGET
	if err = client.Get("/wallet/address", &address); err != nil {
		log.Println("Could not get an address to " + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	// In this one we do the node change first before funding, because a failure
	// in this case results in a support case, but in the other case results in
	// a timing attack where an attacker could possibly drain the prime account
	_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_FUNDED)
	if err != nil {
		log.Println("Could not update the SiaNode status to funded: " + err.Error())
		return err
	}

	var vals = url.Values{}
	vals.Set("amount", siaNode.RequestedCurrency().String())
	vals.Set("destination", address.Address.String())
	var resp api.WalletSiacoinsPOST
	if err = prime.Post("/wallet/siacoins", vals.Encode(), &resp); err != nil {
		log.Println("Could not send coins to " + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	return nil
}

func pollKubeFunded(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeFunded: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get a connection to the sia node" + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	var resp api.WalletGET
	if err = client.Get("/wallet", &resp); err != nil {
		log.Println("Could not get balance for " + siaNode.Shortcode + ": " + err.Error())
		return err
	}
	if resp.ConfirmedSiacoinBalance.Cmp(siaNode.DesiredCurrency()) >= 0 {
		_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_CONFIRMED)
		if err != nil {
			log.Println("Could not update the SiaNode status to confirmed: " + err.Error())
			return err
		}
	}

	return nil
}

func pollKubeConfirmed(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeConfirmed: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get a connection to the sia node" + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	vals := url.Values{}
	vals.Set("funds", siaNode.AllowanceCurrency().String())
	vals.Set("hosts", fmt.Sprintf("%d", siaHosts))
	vals.Set("period", fmt.Sprintf("%d", siaContractPeriod))
	vals.Set("renewwindow", fmt.Sprintf("%d", siaRenewWindow))
	var resp map[string]interface{}
	if err := client.Post("/renter", vals.Encode(), &resp); err != nil {
		log.Println("Could not set renter " + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_CONFIGURED)
	if err != nil {
		log.Println("Could not update the SiaNode status to configured: " + err.Error())
		return err
	}

	return nil
}

func pollKubeConfigured(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeConfigured: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get a connection to the sia node" + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	var resp struct {
		Contracts []map[string]interface{} `json:"contracts"`
	}
	if err := client.Get("/renter/contracts", &resp); err != nil {
		log.Println("Could not set renter " + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	numContracts := len(resp.Contracts)

	if numContracts >= siaNeededContracts {
		_, err := updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_READY)
		if err != nil {
			log.Println("Could not update the SiaNode status to configured: " + err.Error())
			return err
		}
	} else {
		log.Println(fmt.Sprintf("Not enough contracts yet (%d/%d) of total %d", numContracts, siaNeededContracts, siaHosts))
	}

	return nil
}

func pollKubeStopping(clientset *kubernetes.Clientset, siaNode *models.SiaNode) error {
	log.Println("PollKubeStopping: " + siaNode.Shortcode)
	StartFlight(siaFlightPrefix, siaNode)
	defer StopFlight(siaFlightPrefix, siaNode)

	volumeClaims := clientset.PersistentVolumeClaims(kubeNamespace)
	deployments := clientset.AppsV1beta1Client.Deployments(kubeNamespace)
	services := clientset.CoreV1Client.Services(kubeNamespace)
	secrets := clientset.Secrets(kubeNamespace)

	deleteOpts := getDeleteOpts()

	client, err := siaNode.SiaClient()
	if err != nil {
		log.Println("Could not get a connection to the sia node" + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	var curResp api.WalletGET
	log.Println("Getting current wallet balance before stopping: " + siaNode.Shortcode)
	if err = client.Get("/wallet", &curResp); err != nil {
		log.Println("Could not get balance for " + siaNode.Shortcode + ": " + err.Error())
		return err
	}

	total := curResp.ConfirmedSiacoinBalance.Add(curResp.UnconfirmedIncomingSiacoins).Sub(curResp.UnconfirmedOutgoingSiacoins)
	if total.Cmp(types.SiacoinPrecision.Mul64(siaMinerFees)) < 0 {
		total = types.NewCurrency64(0)
	}

	// If it still has a balance, we need to send it back to the prime node
	if !total.IsZero() {
		log.Printf("Sending back %s from %s", total.String(), siaNode.Shortcode)
		prime, err := prime.Server(clientset)
		if err != nil {
			log.Println("Could not get a connection to the prime Sia node: " + err.Error())
			return err
		}

		var address api.WalletAddressGET
		if err = prime.Get("/wallet/address", &address); err != nil {
			log.Println("Could not get an address to " + siaNode.Shortcode + ": " + err.Error())
			return err
		}

		var vals = url.Values{}
		vals.Set("amount", total.String())
		vals.Set("destination", address.Address.String())
		var resp api.WalletSiacoinsPOST
		if err = prime.Post("/wallet/siacoins", vals.Encode(), &resp); err != nil {
			log.Println("Could not send coins to prime node: " + err.Error())
			return err
		}
	}

	log.Println("Getting new wallet balance before stopping: " + siaNode.Shortcode)
	if err = client.Get("/wallet", &curResp); err != nil {
		log.Println("Could not get balance for " + siaNode.Shortcode + ": " + err.Error())
		return err
	}
	total = curResp.ConfirmedSiacoinBalance.Add(curResp.UnconfirmedIncomingSiacoins).Sub(curResp.UnconfirmedOutgoingSiacoins)
	log.Printf("Got wallet balance of %s for %s\n", total.String(), siaNode.Shortcode)
	if total.Cmp(types.SiacoinPrecision.Mul64(siaMinerFees)) < 0 {
		total = types.NewCurrency64(0)
	}
	if !total.IsZero() {
		log.Printf("Found outgoing siacoins for %s, waiting...\n", siaNode.Shortcode)
		return nil
	}

	for i := 0; i < siaNode.MinioInstancesRequested; i++ {
		if err = deleteMinioInstance(clientset, siaNode, i); err != nil {
			log.Println("Could not delete minio instance: " + err.Error())
			return err
		}
	}

	log.Println("Deleting SiaNode deployment: " + siaNode.KubeNameDep())
	if err = deployments.Delete(siaNode.KubeNameDep(), deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: deployment " + siaNode.KubeNameDep())
	}
	log.Println("Deleting SiaNode secret: " + siaNode.KubeNameSec())
	if err = secrets.Delete(siaNode.KubeNameSec(), deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: secret " + siaNode.KubeNameSec())
	}
	log.Println("Deleting SiaNode service: " + siaNode.KubeNameSer())
	if err = services.Delete(siaNode.KubeNameSer(), deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: service " + siaNode.KubeNameSer())
	}
	log.Println("Deleting SiaNode persistent volume claim: " + siaNode.KubeNameVol())
	if err = volumeClaims.Delete(siaNode.KubeNameVol(), deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: pvc " + siaNode.KubeNameVol())
	}

	_, err = updateSiaNodeStatus(siaNode.ID, models.SIANODE_STATUS_STOPPED)
	if err != nil {
		log.Println("Could not update the SiaNode status to stopped: " + err.Error())
		return err
	}

	return nil
}

// Utils

func getDeleteOpts() *metav1.DeleteOptions {
	var deleteGracePeriod int64 = 60
	deletePolicy := metav1.DeletePropagationForeground
	deleteOpts := &metav1.DeleteOptions{}
	deleteOpts.GracePeriodSeconds = &deleteGracePeriod
	deleteOpts.PropagationPolicy = &deletePolicy
	return deleteOpts
}

func getPod(clientset *kubernetes.Clientset, siaNode *models.SiaNode) (*v1.Pod, error) {
	opts := metav1.ListOptions{LabelSelector: "app=" + siaNode.KubeNameApp()}
	pods, err := clientset.Pods(kubeNamespace).List(opts)
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting pod from kubernetes: " + err.Error())
		return nil, err
	}

	if pods == nil || errors.IsNotFound(err) || pods.Items == nil || len(pods.Items) == 0 {
		log.Println("Not found yet: " + siaNode.KubeNameApp())
		return nil, err
	}

	return &pods.Items[0], nil
}

// Administrative API Calls

func getPendingSiaNodes() ([]*models.SiaNode, error) {
	return getSiaNodes("/sianodes/pending/all")
}

func getReadyOrphanedSiaNodes() ([]*models.SiaNode, error) {
	return getSiaNodes("/sianodes/orphaned/ready")
}

func getSiaNodes(urlPart string) ([]*models.SiaNode, error) {
	url := fmt.Sprintf("%s%s?secret=%s", URLRoot, urlPart, SiaCDNSecretKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Could not create request GET " + url)
		return nil, err
	}

	res, err := cliClient.Do(req)
	if err != nil {
		log.Println("Error making getSiaNodes request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read getSiaNodes response: " + err.Error())
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

func updateSiaNodeStatus(id uuid.UUID, status string) (*models.SiaNode, error) {
	url := fmt.Sprintf("%s/sianodes/status?secret=%s", URLRoot, SiaCDNSecretKey)

	reqBodyData := struct {
		Id     uuid.UUID `json:"id"`
		Status string    `json:"status"`
	}{Id: id, Status: status}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(reqBodyData)
	if err != nil {
		log.Println("Could not encode sia node update json: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Println("Could not create request POST " + url)
		return nil, err
	}

	res, err := cliClient.Do(req)
	if err != nil {
		log.Println("Error making updateSiaNodeStatus request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read updateSiaNodeStatus response: " + err.Error())
		return nil, err
	}

	var node struct {
		SiaNode *models.SiaNode `json:"sianode"`
	}
	if err = json.Unmarshal(body, &node); err != nil {
		log.Println("Could not decode response: " + err.Error())
		return nil, err
	}
	return node.SiaNode, nil
}

func updateSiaNodeMinioInstancesActivated(id uuid.UUID, instances int) (*models.SiaNode, error) {
	url := fmt.Sprintf("%s/sianodes/id/%s?secret=%s", URLRoot, id.String(), SiaCDNSecretKey)

	reqBodyData := struct {
		Instances int `json:"minio_instances_activated"`
	}{Instances: instances}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(reqBodyData)
	if err != nil {
		log.Println("Could not encode sia node update json: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Println("Could not create request POST " + url)
		return nil, err
	}

	res, err := cliClient.Do(req)
	if err != nil {
		log.Println("Error making updateSiaNodeMinioInstancesActivated request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read updateSiaNodeMinioInstancesActivated response: " + err.Error())
		return nil, err
	}

	var node struct {
		SiaNode *models.SiaNode `json:"sianode"`
	}
	if err = json.Unmarshal(body, &node); err != nil {
		log.Println("Could not decode response: " + err.Error())
		return nil, err
	}
	return node.SiaNode, nil
}

func createWalletSeed(siaNodeID uuid.UUID, words string) (*models.WalletSeed, error) {
	url := fmt.Sprintf("%s/wallets/%s/seed?secret=%s", URLRoot, siaNodeID, SiaCDNSecretKey)

	reqBodyData := struct {
		Words string `json:"words"`
	}{Words: words}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(reqBodyData)
	if err != nil {
		log.Println("Could not encode sia wallet seed create json: " + err.Error())
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Println("Could not create request POST " + url)
		return nil, err
	}

	res, err := cliClient.Do(req)
	if err != nil {
		log.Println("Error making createWalletSeed request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read createWalletSeed response: " + err.Error())
		return nil, err
	}

	var resp struct {
		WalletSeed *models.WalletSeed `json:"wallet_seed"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Println("Could not decode response: " + err.Error())
		return nil, err
	}
	return resp.WalletSeed, nil
}

func getWalletSeed(siaNodeID uuid.UUID) (*models.WalletSeed, error) {
	url := fmt.Sprintf("%s/wallets/%s/seed?secret=%s",
		URLRoot, siaNodeID.String(), SiaCDNSecretKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Could not create request GET " + url)
		return nil, err
	}

	res, err := cliClient.Do(req)
	if err != nil {
		log.Println("Error making getWalletSeed request: " + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read getWalletSeed response: " + err.Error())
		return nil, err
	}

	var resp struct {
		Seed *models.WalletSeed `json:"wallet_seed"`
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		log.Println("Could not decode response: " + err.Error())
		return nil, err
	}

	return resp.Seed, nil
}

//////////////////////////////////
//////// MINIO STUFF /////////////
//////////////////////////////////

func deployMinio(clientset *kubernetes.Clientset, siaNode *models.SiaNode, instance int) error {
	name := siaNode.KubeNameMinio(instance)
	nfsName := siaNode.KubeNameMinioNFS(instance)
	hostname := siaNode.MinioHostname(instance)
	pvName := nfsName + "-pv"
	pvcName := nfsName + "-pvc"
	certName := fmt.Sprintf("%s-cert", name)
	mountPath := fmt.Sprintf("/minio%d", instance+1)

	deployments := clientset.AppsV1beta1Client.Deployments(kubeNamespace)
	services := clientset.CoreV1Client.Services(kubeNamespace)
	secrets := clientset.Secrets(kubeNamespace)
	volumeClaims := clientset.PersistentVolumeClaims(kubeNamespace)
	volumes := clientset.PersistentVolumes()
	ingresses := clientset.Ingresses(kubeNamespace)

	if siaNode.Status != models.SIANODE_STATUS_READY {
		log.Printf("Waiting for sianode %s status to be: %s (currently: %s)", siaNode.Shortcode, models.SIANODE_STATUS_READY, siaNode.Status)
		return nil
	}

	// First check for nfs volume claim
	log.Println("Checking NFS volume claim: " + nfsName)
	nfsClaim, err := volumeClaims.Get(nfsName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting volume claim from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if nfsClaim == nil || errors.IsNotFound(err) {
		nfsClaim = &v1.PersistentVolumeClaim{}
		nfsClaim.Name = nfsName
		nfsClaim.Namespace = kubeNamespace
		nfsClaim.Spec = v1.PersistentVolumeClaimSpec{
			StorageClassName: &minioStorageClass,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{v1.ResourceName("storage"): minioDefaultStorage},
			},
		}

		log.Println("Creating nfs volume claim " + nfsName)
		nfsClaim, err = volumeClaims.Create(nfsClaim)
		if err != nil {
			log.Println("Error creating nfs volume claim: " + err.Error())
			return err
		}
	} else {
		log.Println("Found nfs volume claim " + nfsName)
	}

	// Check for the NFS service
	log.Println("Checking NFS service: " + nfsName)
	nfsService, err := services.Get(nfsName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting service from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if nfsService == nil || errors.IsNotFound(err) {
		nfsService = &v1.Service{}
		nfsService.Name = nfsName
		nfsService.Namespace = kubeNamespace
		nfsService.Spec = v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
				v1.ServicePort{Name: "nfs", Port: 2049},
				v1.ServicePort{Name: "mountd", Port: 20048},
				v1.ServicePort{Name: "rpcbind", Port: 111},
			},
			Selector: map[string]string{"app": nfsName},
		}
		log.Println("Creating nfs service " + nfsName)
		nfsService, err = services.Create(nfsService)
		if err != nil {
			log.Println("Error creating nfs service: " + err.Error())
			return err
		}
	} else {
		log.Println("Found nfs service " + nfsName)
	}

	// Now check for nfs deployment
	log.Println("Checking NFS deployment: " + nfsName)
	nfsDeployment, err := deployments.Get(nfsName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting deployment from kubernetes: " + err.Error())
		return err
	}
	// If nfs deployment doesn't exist, create it
	if nfsDeployment == nil || errors.IsNotFound(err) {
		nfsDeployment := &v1beta1.Deployment{}
		nfsDeployment.Name = nfsName
		nfsDeployment.Namespace = kubeNamespace
		nfsDeployment.Spec = v1beta1.DeploymentSpec{Template: v1.PodTemplateSpec{}}
		nfsDeployment.Spec.Strategy.Type = v1beta1.RecreateDeploymentStrategyType
		nfsDeployment.Spec.Template.Labels = map[string]string{"app": nfsName}
		nfsDeployment.Spec.Template.Spec = v1.PodSpec{
			Volumes: []v1.Volume{
				v1.Volume{
					Name: "nfs",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: nfsName,
						},
					},
				},
			},
			Containers: []v1.Container{
				v1.Container{
					Name:            name,
					Image:           "gcr.io/google_containers/volume-nfs:0.8",
					ImagePullPolicy: v1.PullAlways,
					Ports: []v1.ContainerPort{
						v1.ContainerPort{Name: "nfs", ContainerPort: 2049},
						v1.ContainerPort{Name: "mountd", ContainerPort: 20048},
						v1.ContainerPort{Name: "rpcbind", ContainerPort: 111},
					},
					SecurityContext: &v1.SecurityContext{
						Privileged: &securityContextPrivileged,
					},
					VolumeMounts: []v1.VolumeMount{
						v1.VolumeMount{Name: "nfs", MountPath: "/exports"},
					},
				},
			},
		}
		log.Println("Creating nfs deployment " + name)
		nfsDeployment, err = deployments.Create(nfsDeployment)
		if err != nil {
			log.Println("Error creating nfs deployment: " + err.Error())
			return err
		}
	} else {
		log.Println("Found nfs deployment " + nfsName)
	}

	// Check for the Sia persistent volume
	log.Println("Checking minio persistent volume: " + pvName)
	pv, err := volumes.Get(pvName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting persistent volume from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if pv == nil || errors.IsNotFound(err) {
		pv = &v1.PersistentVolume{}
		pv.Name = pvName
		pv.Namespace = kubeNamespace
		pv.Spec = v1.PersistentVolumeSpec{
			Capacity: v1.ResourceList{
				"storage": minioDefaultStorage,
			},
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteMany,
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Server: nfsService.Spec.ClusterIP,
					Path:   "/",
				},
			},
		}
		pv.Labels = map[string]string{"app": pvName}
		log.Println("Creating nfs persistent volume " + pvName)
		pv, err = volumes.Create(pv)
		if err != nil {
			log.Println("Error creating nfs persistent volume: " + err.Error())
			return err
		}
	} else {
		log.Println("Found nfs persistent volume " + pvName)
	}

	// Now check for minio volume claim
	log.Println("Checking minio persistent volume claim: " + pvcName)
	claim, err := volumeClaims.Get(pvcName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting volume claim from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if claim == nil || errors.IsNotFound(err) {
		claim = &v1.PersistentVolumeClaim{}
		claim.Name = pvcName
		claim.Namespace = kubeNamespace
		claim.Spec = v1.PersistentVolumeClaimSpec{
			StorageClassName: &nfsStorageClass,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{v1.ResourceName("storage"): minioDefaultStorage},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": pvName},
			},
		}

		log.Println("Creating volume claim " + pvcName)
		claim, err = volumeClaims.Create(claim)
		if err != nil {
			log.Println("Error creating volume claim: " + err.Error())
			return err
		}
	} else {
		log.Println("Found volume claim " + pvcName)
	}

	// Now check the sia deployment to make sure it mounts the minio volume
	siaDeployment, err := deployments.Get(siaNode.KubeNameDep(), metav1.GetOptions{})
	if err != nil {
		log.Println("Error getting deployment from kubernetes: " + err.Error())
		return err
	}

	siaDeploymentChanged := false

	var volume *v1.Volume
	for _, vol := range siaDeployment.Spec.Template.Spec.Volumes {
		if vol.Name == pvcName {
			volume = &vol
			break
		}
	}
	// If it doesn't include the volume in the spec, add it
	if volume == nil {
		volumes := siaDeployment.Spec.Template.Spec.Volumes
		volumes = append(volumes, v1.Volume{
			Name: pvcName,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		})
		siaDeployment.Spec.Template.Spec.Volumes = volumes
		siaDeploymentChanged = true
	}

	var volumeMount *v1.VolumeMount
	for _, container := range siaDeployment.Spec.Template.Spec.Containers {
		for _, vm := range container.VolumeMounts {
			if vm.Name == pvcName {
				volumeMount = &vm
				break
			}
		}
	}
	// If it doesn't mount the volume in the containers, mount it on each one,
	// but making sure to do so in read-only mode.
	if volumeMount == nil {
		newContainers := []v1.Container{}
		for _, container := range siaDeployment.Spec.Template.Spec.Containers {
			container.VolumeMounts = append(container.VolumeMounts, v1.VolumeMount{
				Name:      pvcName,
				MountPath: mountPath,
			})
			newContainers = append(newContainers, container)
		}
		siaDeployment.Spec.Template.Spec.Containers = newContainers
		siaDeploymentChanged = true
	}

	if siaDeploymentChanged {
		log.Println("Changing the Sia deployment in the process of deploying Minio " + name)
		siaDeployment.Spec.Strategy.Type = v1beta1.RecreateDeploymentStrategyType // TODO: Remove this line sometime, this is just catching old ones to switch it over
		siaDeployment, err = deployments.Update(siaDeployment)
		if err != nil {
			log.Println("Error updating sia deployment with new volume info in kubernetes: " + err.Error())
			return err
		}
	}

	// Now check for service
	log.Println("Checking minio service: " + name)
	service, err := services.Get(name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting minio service from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if service == nil || errors.IsNotFound(err) {
		service = &v1.Service{}
		service.Name = name
		service.Namespace = kubeNamespace
		service.Spec = v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
				v1.ServicePort{Port: 9000, TargetPort: intstr.FromInt(9000), Protocol: v1.ProtocolTCP},
			},
			Selector: map[string]string{"app": name},
		}
		log.Println("Creating service " + name)
		service, err = services.Create(service)
		if err != nil {
			log.Println("Error creating service: " + err.Error())
			return err
		}
	} else {
		log.Println("Found service " + name)
	}

	// Check for secret
	log.Println("Checking minio secret: " + name)
	secret, err := secrets.Get(name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting secret from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if secret == nil || errors.IsNotFound(err) {
		secret = &v1.Secret{}
		secret.Name = name
		secret.Namespace = kubeNamespace
		secret.Type = v1.SecretTypeOpaque
		secret.Data = map[string][]byte{
			"accesskey": []byte(siaNode.MinioAccessKey),
			"secretkey": []byte(siaNode.MinioSecretKey),
		}
		log.Println("Creating secret " + name)
		secret, err = secrets.Create(secret)
		if err != nil {
			log.Println("Error creating secret: " + err.Error())
			return err
		}
	} else {
		log.Println("Found secret " + name)
	}

	// Check for deployment
	log.Println("Checking minio deployment: " + name)
	deployment, err := deployments.Get(name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting deployment from kubernetes: " + err.Error())
		return err
	}
	// If deployment doesn't exist, create it
	if deployment == nil || errors.IsNotFound(err) {
		deployment := &v1beta1.Deployment{}
		deployment.Name = name
		deployment.Namespace = kubeNamespace
		deployment.Spec = v1beta1.DeploymentSpec{Template: v1.PodTemplateSpec{}}
		deployment.Spec.Strategy.Type = v1beta1.RecreateDeploymentStrategyType
		deployment.Spec.Template.Labels = map[string]string{"app": name}
		deployment.Spec.Template.Spec = v1.PodSpec{
			Volumes: []v1.Volume{
				v1.Volume{
					Name: pvcName,
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				},
			},
			Containers: []v1.Container{
				v1.Container{
					Name:            name,
					Image:           "gcr.io/gradientzoo-1233/siacdn-minio:latest",
					ImagePullPolicy: v1.PullAlways,
					Ports: []v1.ContainerPort{
						v1.ContainerPort{ContainerPort: 9000},
					},
					VolumeMounts: []v1.VolumeMount{
						v1.VolumeMount{Name: pvcName, MountPath: mountPath},
					},
					Env: []v1.EnvVar{
						v1.EnvVar{
							Name: "SIA_API_PASSWORD",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "sia-secret",
									},
									Key: "siaapipassword",
								},
							},
						},
						v1.EnvVar{
							Name: "MINIO_ACCESS_KEY",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: name,
									},
									Key: "accesskey",
								},
							},
						},
						v1.EnvVar{
							Name: "MINIO_SECRET_KEY",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: name,
									},
									Key: "secretkey",
								},
							},
						},
						v1.EnvVar{
							Name:  "SIA_DAEMON_ADDR",
							Value: siaNode.KubeNameSer() + ".sia.svc.cluster.local:9980",
						},
						v1.EnvVar{
							Name:  "SIA_CACHE_DIR",
							Value: filepath.Join(mountPath, "siacache"),
						},
						v1.EnvVar{
							Name:  "SIA_DB_FILE",
							Value: filepath.Join(mountPath, "sia.db"),
						},
						v1.EnvVar{
							Name:  "SIA_CACHE_MAX_SIZE_BYTES",
							Value: "90000000000",
						},
						v1.EnvVar{
							Name:  "SIA_CACHE_PURGE_AFTER_SEC",
							Value: "345600",
						},
						v1.EnvVar{
							Name:  "SIA_BACKGROUND_UPLOAD",
							Value: "0",
						},
					},
				},
			},
		}
		log.Println("Creating deployment " + name)
		deployment, err = deployments.Create(deployment)
		if err != nil {
			log.Println("Error creating deployment: " + err.Error())
			return err
		}
	} else {
		log.Println("Found deployment " + name)
	}

	// Check for ingress
	log.Println("Checking minio ingress: " + name)
	ing, err := ingresses.Get(name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting minio ingress from kubernetes: " + err.Error())
		return err
	}
	// If it doesn't exist, create it
	if ing == nil || errors.IsNotFound(err) {
		ing = &extensions.Ingress{}
		ing.Name = name
		ing.Namespace = kubeNamespace
		ing.Annotations = map[string]string{
			"kubernetes.io/tls-acme":                   "true",
			"kubernetes.io/ingress.class":              "nginx",
			"ingress.kubernetes.io/force-ssl-redirect": "true",
			"ingress.kubernetes.io/proxy-body-size":    "4g",
		}
		ing.Spec = extensions.IngressSpec{
			Backend: &extensions.IngressBackend{
				ServiceName: name,
				ServicePort: intstr.FromInt(9000),
			},
			TLS: []extensions.IngressTLS{extensions.IngressTLS{
				SecretName: certName,
				Hosts:      []string{hostname},
			}},
			Rules: []extensions.IngressRule{extensions.IngressRule{
				Host: hostname,
				IngressRuleValue: extensions.IngressRuleValue{
					HTTP: &extensions.HTTPIngressRuleValue{
						Paths: []extensions.HTTPIngressPath{extensions.HTTPIngressPath{
							Path: "/",
							Backend: extensions.IngressBackend{
								ServiceName: name,
								ServicePort: intstr.FromInt(9000),
							},
						}},
					},
				},
			}},
		}
		log.Println("Creating ingress " + name)
		ing, err = ingresses.Create(ing)
		if err != nil {
			log.Println("Error creating ingress: " + err.Error())
			return err
		}
	} else {
		log.Println("Found ingress " + name)
	}

	// Wait for cert
	log.Println("Checking minio cert: " + certName)
	cert, err := secrets.Get(certName, metav1.GetOptions{})
	if err != nil {
		log.Println("Error getting cert (" + name + ") from kubernetes: " + err.Error())
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if len(cert.Data) == 0 {
		log.Println("Waiting for certificate: " + name)
		return fmt.Errorf("Waiting for certificate: %s", name)
	}

	siaNode, err = updateSiaNodeMinioInstancesActivated(siaNode.ID, siaNode.MinioInstancesActivated+1)
	if err != nil {
		log.Println("Could not update the number of minio instances activated: " + err.Error())
		return err
	}

	log.Println("Successfully activated Minio instance " + name + "!")
	return nil
}

func deleteMinioInstance(clientset *kubernetes.Clientset, siaNode *models.SiaNode, instance int) error {
	name := siaNode.KubeNameMinio(instance)
	nfsName := siaNode.KubeNameMinioNFS(instance)
	pvName := nfsName + "-pv"
	pvcName := nfsName + "-pvc"

	deployments := clientset.AppsV1beta1Client.Deployments(kubeNamespace)
	services := clientset.CoreV1Client.Services(kubeNamespace)
	secrets := clientset.Secrets(kubeNamespace)
	volumeClaims := clientset.PersistentVolumeClaims(kubeNamespace)
	volumes := clientset.PersistentVolumes()
	ingresses := clientset.Ingresses(kubeNamespace)

	deleteOpts := getDeleteOpts()

	log.Println("Deleting minio ingress: " + name)
	if err := ingresses.Delete(name, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: ingress " + name)
	}
	log.Println("Deleting minio deployment: " + name)
	if err := deployments.Delete(name, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: deployment " + name)
	}
	log.Println("Deleting minio secret: " + name)
	if err := secrets.Delete(name, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: secret " + name)
	}
	log.Println("Deleting minio service: " + name)
	if err := services.Delete(name, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: service " + name)
	}

	// Now check the sia deployment to make sure it doesn't mount the minio volume
	siaDeployment, err := deployments.Get(siaNode.KubeNameDep(), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Println("Error getting deployment from kubernetes: " + err.Error())
		return err
	}
	// If there was no error, we have a deployment to modify
	if err == nil {
		siaDeploymentChanged := false

		newVolumes := make([]v1.Volume, 0, len(siaDeployment.Spec.Template.Spec.Volumes))
		for _, vol := range siaDeployment.Spec.Template.Spec.Volumes {
			if vol.Name != pvcName {
				newVolumes = append(newVolumes, vol)
			}
		}
		// If it doesn't include the volume in the spec, add it
		if len(newVolumes) != len(siaDeployment.Spec.Template.Spec.Volumes) {
			siaDeploymentChanged = true
		}

		newContainers := make([]v1.Container, 0, len(siaDeployment.Spec.Template.Spec.Containers))
		for _, container := range siaDeployment.Spec.Template.Spec.Containers {
			newVolumeMounts := make([]v1.VolumeMount, 0, len(container.VolumeMounts))
			for _, vm := range container.VolumeMounts {
				if vm.Name != pvcName {
					newVolumeMounts = append(newVolumeMounts, vm)
				}
			}
			if len(newVolumeMounts) != len(container.VolumeMounts) {
				container.VolumeMounts = newVolumeMounts
				siaDeploymentChanged = true
			}
			newContainers = append(newContainers, container)
		}

		if siaDeploymentChanged {
			log.Println("Changing the Sia deployment in the process of deleting Minio " + name)
			siaDeployment.Spec.Template.Spec.Containers = newContainers
			siaDeployment, err = deployments.Update(siaDeployment)
			if err != nil {
				log.Println("Error updating sia deployment with new volume info in kubernetes: " + err.Error())
				return err
			}
		}
	}

	log.Println("Deleting minio NFS volume claim: " + pvcName)
	if err = volumeClaims.Delete(pvcName, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: pvc " + pvcName)
	}
	log.Println("Deleting minio NFS volume: " + pvName)
	if err = volumes.Delete(pvName, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: pv " + pvName)
	}
	log.Println("Deleting minio NFS deployment: " + nfsName)
	if err = deployments.Delete(nfsName, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: deployment " + nfsName)
	}
	log.Println("Deleting minio NFS service: " + nfsName)
	if err = services.Delete(nfsName, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: service " + nfsName)
	}
	log.Println("Deleting minio NFS volume claim: " + nfsName)
	if err = volumeClaims.Delete(nfsName, deleteOpts); err != nil && !errors.IsNotFound(err) {
		return err
	} else if errors.IsNotFound(err) {
		log.Println("Not found: pvc " + nfsName)
	}

	return nil
}
