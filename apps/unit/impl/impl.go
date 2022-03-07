package impl

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/upmio/unitctl/apps/unit"
)

type impl struct {
	kubeclientset kubernetes.Interface
	logger        *zap.SugaredLogger
}

func NewUnit(logger *zap.SugaredLogger) (unit.UnitClient, error) {
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
		logger:        logger,
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

func (i *impl) GetMysqlSet(ctx context.Context, namespace, svcGroupName string) (unit.MysqlSet, error) {
	podList, err := i.kubeclientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set{
			unit.SvcGroupNameLabel: svcGroupName,
			unit.SvcTypeLabel:      "mysql",
		}.String(),
	})

	if err != nil {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, fmt.Errorf("get mysql pod list zero")
	}

	var mysqlSet = make(unit.MysqlSet, 0)

	for _, pod := range podList.Items {
		if pod.Labels[unit.ReadOnlyLabel] == "false" {
			podIp := pod.Status.PodIP
			for _, container := range pod.Spec.Containers {
				if container.Name == "mysql" {
					podPort := int(container.Ports[0].ContainerPort)
					var mysql = unit.NewMysql(podIp, podPort)
					mysqlSet = append(mysqlSet, mysql)
					i.logger.Infof("get master server %s:%d", mysql.IpAddr, mysql.Port)
				}
			}
		}
	}

	if len(mysqlSet) > 1 {
		return nil, fmt.Errorf("get readonly label false pod result count more than 1")
	}

	return mysqlSet, nil
}
