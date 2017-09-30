package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NebulousLabs/Sia/api"
	"github.com/NebulousLabs/Sia/types"
	"github.com/google/uuid"
	"github.com/thegreatdb/siacdn/backend/kube"
	"github.com/thegreatdb/siacdn/backend/randstring"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const siaPerTB = 210.39 // https://siastats.info/storage_pricing.html
//const siaPerTB = 2.1 // Just for testing
const slopMultiple = 1.05
const numTerms = 1.2
const kubeNamespace = "sia"

const SIANODE_STATUS_CREATED = "created"           // This object has been created
const SIANODE_STATUS_DEPLOYED = "deployed"         // Deployment yaml has been sent to kube
const SIANODE_STATUS_INSTANCED = "instanced"       // Pod instance has been seen
const SIANODE_STATUS_SNAPSHOTTED = "snapshotted"   // Bootstrap snapshot script has finished
const SIANODE_STATUS_SYNCHRONIZED = "synchronized" // Blockchain has synchronized
const SIANODE_STATUS_INITIALIZED = "initialized"   // Wallet has been initialized
const SIANODE_STATUS_UNLOCKED = "unlocked"         // Wallet has been unlocked
const SIANODE_STATUS_FUNDED = "funded"             // Account has received initial funding
const SIANODE_STATUS_CONFIRMED = "confirmed"       // Funding has been confirmed by the server
const SIANODE_STATUS_CONFIGURED = "configured"     // Allowance has been set
const SIANODE_STATUS_READY = "ready"               // Everything is ready to go and contracts are all set
const SIANODE_STATUS_STOPPED = "stopped"           // A user or administrator has stopped the node
const SIANODE_STATUS_DEPLETED = "depleted"         // All the SiaCoins in the accoint have been used up
const SIANODE_STATUS_ERROR = "error"               // An error has occurred in the process

type SiaNode struct {
	ID                      uuid.UUID `json:"id"`
	Shortcode               string    `json:"shortcode"`
	AccountID               uuid.UUID `json:"account_id"`
	Capacity                float32   `json:"capacity"`
	Status                  string    `json:"status"`
	MinioInstancesRequested int       `json:"minio_instances_requested"`
	MinioInstancesActivated int       `json:"minio_instances_activated"`
	MinioAccessKey          string    `json:"minio_access_key"`
	MinioSecretKey          string    `json:"minio_secret_key"`
	CreatedTime             time.Time `json:"created_time"`
}

func NewSiaNode(accountID uuid.UUID, capacity float32) (*SiaNode, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &SiaNode{
		ID:          id,
		Shortcode:   randstring.New(8),
		AccountID:   accountID,
		Capacity:    capacity,
		Status:      SIANODE_STATUS_CREATED,
		CreatedTime: time.Now().UTC(),
	}, nil
}

func (sn *SiaNode) Copy() *SiaNode {
	var cpy SiaNode
	cpy = *sn
	return &cpy
}

func (sn *SiaNode) SiaClient() (*api.Client, error) {
	clientset, clusterType, err := kube.KubeClient()
	if err != nil {
		return nil, err
	}
	services := clientset.Services(kubeNamespace)
	service, err := services.Get(sn.KubeNameSer(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, fmt.Errorf("Service not found: " + sn.Shortcode)
	}
	var ip string
	if clusterType == kube.ClusterTypeExternal {
		log.Println("Found external cluster, assuming local tunnel. " +
			"Run the following command:\n\tkubectl --namespace=sia port-forward " +
			"$(kubectl --namespace=sia get pods -l app=" + sn.KubeNameApp() +
			" -o jsonpath=\"{.items[0].metadata.name}\") 9980")
		ip = "localhost"
	} else {
		ip = service.Spec.ClusterIP
	}
	client := api.NewClient(ip+":9980", os.Getenv("SIA_API_PASSWORD"))
	return client, nil
}

func (sn *SiaNode) AllowanceCurrency() types.Currency {
	return types.SiacoinPrecision.MulFloat(float64(siaPerTB * sn.Capacity))
}

func (sn *SiaNode) RequestedCurrency() types.Currency {
	return types.SiacoinPrecision.MulFloat(float64(siaPerTB * numTerms * slopMultiple * sn.Capacity))
}

func (sn *SiaNode) DesiredCurrency() types.Currency {
	return types.SiacoinPrecision.MulFloat(float64(siaPerTB * numTerms * sn.Capacity))
}

func (sn *SiaNode) ValidateStatus() error {
	switch sn.Status {
	case SIANODE_STATUS_CREATED,
		SIANODE_STATUS_DEPLOYED,
		SIANODE_STATUS_INSTANCED,
		SIANODE_STATUS_SNAPSHOTTED,
		SIANODE_STATUS_SYNCHRONIZED,
		SIANODE_STATUS_INITIALIZED,
		SIANODE_STATUS_UNLOCKED,
		SIANODE_STATUS_FUNDED,
		SIANODE_STATUS_CONFIRMED,
		SIANODE_STATUS_CONFIGURED,
		SIANODE_STATUS_READY,
		SIANODE_STATUS_STOPPED,
		SIANODE_STATUS_DEPLETED,
		SIANODE_STATUS_ERROR:
		return nil
	default:
		return fmt.Errorf("Invalid SiaNode status: '%s'", sn.Status)
	}
}

func (sn *SiaNode) Pending() bool {
	switch sn.Status {
	case SIANODE_STATUS_READY,
		SIANODE_STATUS_STOPPED,
		SIANODE_STATUS_DEPLETED,
		SIANODE_STATUS_ERROR:
		return false
	}
	return true
}

func (sn *SiaNode) KubeNameBase() string {
	return fmt.Sprintf("siacdn-%s", sn.Shortcode)
}

func (sn *SiaNode) KubeNameApp() string {
	return sn.KubeNameBase()
}

func (sn *SiaNode) KubeNameDep() string {
	return sn.KubeNameBase()
}

func (sn *SiaNode) KubeNameVol() string {
	return sn.KubeNameBase()
}

func (sn *SiaNode) KubeNameSer() string {
	return sn.KubeNameBase()
}

func (sn *SiaNode) KubeNameSec() string {
	return sn.KubeNameBase()
}

func (sn *SiaNode) KubeNameMinio(instance int) string {
	return fmt.Sprintf("siacdn-%s-minio%d", sn.Shortcode, instance)
}

func (sn *SiaNode) KubeNameMinioNFS(instance int) string {
	return fmt.Sprintf("siacdn-%s-minio%d-nfs", sn.Shortcode, instance)
}
