package prime

import (
	"fmt"
	"log"
	"os"

	"github.com/NebulousLabs/Sia/api"
	"github.com/thegreatdb/siacdn/siacdn-backend/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const kubeNamespace = "sia"

func Server(clientset *kubernetes.Clientset) (*api.Client, error) {
	_, clusterType, err := kube.KubeClient()
	if err != nil {
		return nil, err
	}
	var host string
	if clusterType == kube.ClusterTypeExternal {
		log.Println("Found external Sia Prime cluster, assuming local tunnel. " +
			"Run the following command:\n\tkubectl --namespace=sia port-forward " +
			"$(kubectl --namespace=sia get pods -l app=sia-prime" +
			" -o jsonpath=\"{.items[0].metadata.name}\") 9985:9980")
		host = "localhost:9985"
	} else {
		services := clientset.Services(kubeNamespace)
		service, err := services.Get("sia-prime-service", metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if service == nil {
			return nil, fmt.Errorf("Prime service not found")
		}
		host = service.Spec.ClusterIP + ":9980"
	}
	client := api.NewClient(host, os.Getenv("SIA_API_PASSWORD"))
	return client, nil
}
