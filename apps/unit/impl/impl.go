package impl

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/upmio/unitctl/apps/unit"
)

type impl struct {
	kubeclientset kubernetes.Interface
}

func NewUnitImpl() (unit.UnitClient, error) {
	// create incluster config object
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("create incluster config fail, error: %v", err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create clientset fail, error: %v", err)
	}

	return &impl{
		kubeclientset: clientset,
	}, nil
}

func (i *impl) GetSecret(ctx context.Context, namespace, secretName string) (unit.SecretInfo, error) {
	secret, err := i.kubeclientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret.Data, nil
}

func (i *impl) GetConfigmap(ctx context.Context, namespace, configMapName string) (unit.ConfigMapInfo, error) {
	configMap, err := i.kubeclientset.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return configMap.Data, nil
}
