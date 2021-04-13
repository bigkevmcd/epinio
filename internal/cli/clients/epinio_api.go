package clients

import (
	"fmt"

	"github.com/epinio/epinio/deployments"
	"github.com/epinio/epinio/internal/cli/config"
	"github.com/epinio/epinio/kubernetes"
	"github.com/pkg/errors"
)

// EpinioApiClient provides functionality for talking to an Epinio API
// server on Kubernetes
type EpinioApiClient struct {
	URL string
}

var epinioClientMemo *EpinioApiClient

func GetEpinioApiClient() (*EpinioApiClient, error) {
	if epinioClientMemo != nil {
		return epinioClientMemo, nil
	}

	configConfig, err := config.Load()
	if err != nil {
		return nil, err
	}

	cluster, err := kubernetes.GetCluster()
	if err != nil {
		return nil, err
	}

	epinioURL, err := getEpinioURL(configConfig, cluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve epinio api host")
	}

	epinioClient := &EpinioApiClient{
		URL: epinioURL,
	}

	epinioClientMemo = epinioClient

	return epinioClient, nil
}

// getEpinioURL finds the URL for epinio
func getEpinioURL(config *config.Config, cluster *kubernetes.Cluster) (string, error) {
	// Get the ingress
	ingresses, err := cluster.ListIngress(deployments.EpinioDeploymentID, "app.kubernetes.io/name=epinio")
	if err != nil {
		return "", errors.Wrap(err, "failed to list ingresses for epinio api server")
	}

	if len(ingresses.Items) < 1 {
		return "", errors.New("epinio api ingress not found")
	}

	if len(ingresses.Items) > 1 {
		return "", errors.New("more than one epinio api ingress found")
	}

	if len(ingresses.Items[0].Spec.Rules) < 1 {
		return "", errors.New("epinio api ingress has no rules")
	}

	if len(ingresses.Items[0].Spec.Rules) > 1 {
		return "", errors.New("epinio api ingress has more than on rule")
	}

	host := ingresses.Items[0].Spec.Rules[0].Host

	return fmt.Sprintf("%s://%s", config.EpinioProtocol, host), nil
}
