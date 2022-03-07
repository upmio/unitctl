package unit

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"
)

const (
	ReadOnlyLabel     = "dbscale.proxysql.readonly"
	SvcGroupNameLabel = "dbscale.service.group"
	SvcTypeLabel      = "dbscale.service.image.name"
)

type UnitClient interface {
	GetSecret(ctx context.Context, namespace string, secretName string) (SecretInfo, error)
	GetConfigmap(ctx context.Context, namespace, configMapName string) (ConfigMapInfo, error)
	GetMysqlSet(ctx context.Context, namespace, svcGroupName string) (MysqlSet, error)
}

type SecretInfo map[string][]byte

func (s SecretInfo) Marshal() ([]byte, error) {
	var output = make(map[string]string, 0)
	for key, value := range s {
		output[key] = string(value)
	}
	return json.Marshal(output)
}

type ConfigMapInfo map[string]string

func (c ConfigMapInfo) CreateConfig(fileDir string) error {
	for key, value := range c {
		filePath := path.Join(fileDir, key)

		err := ioutil.WriteFile(filePath, []byte(value), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

type Mysql struct {
	IpAddr string
	Port   int
}

func NewMysql(ip string, port int) *Mysql {
	return &Mysql{
		IpAddr: ip,
		Port:   port,
	}
}

type MysqlSet []*Mysql
