package components

import (
	"cmp"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	core "github.com/bucketheadv/infra-core"
	"github.com/sirupsen/logrus"
)

type apolloChangeListener struct{}

func (c *apolloChangeListener) OnChange(event *storage.ChangeEvent) {
	for k, v := range event.Changes {
		logrus.Infof("apollo %v config changed, key: %v, old value: %v, new value: %v",
			event.Namespace, k, v.OldValue, v.NewValue)
	}
}

func (c *apolloChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	logrus.Infof("Apollo config pull, namespace [%s] updated to latest version", event.Namespace)
}

type ApolloConf struct {
	Enabled        bool   `json:"enabled"`
	AppID          string `json:"appId"`
	Cluster        string `json:"cluster"`
	NamespaceName  string `json:"namespaceName"`
	IP             string `json:"ip"`
	IsBackupConfig bool   `default:"true" json:"isBackupConfig"`
}

var apolloClient agollo.Client

func InitApolloClient(c ApolloConf, onSuccess func()) {
	if !c.Enabled {
		return
	}
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		var appConfig = &config.AppConfig{
			AppID:          c.AppID,
			Cluster:        c.Cluster,
			NamespaceName:  c.NamespaceName,
			IP:             c.IP,
			IsBackupConfig: c.IsBackupConfig,
		}
		return appConfig, nil
	})

	if err != nil {
		logrus.Infof("初始化Apollo失败, %s", err.Error())
		return
	}

	client.AddChangeListener(&apolloChangeListener{})
	apolloClient = client
	if onSuccess != nil {
		onSuccess()
	}
}

func ApolloApplicationConfig(key string) string {
	return apolloClient.GetValue(key)
}

func ApolloNamespace(namespace string) *storage.Config {
	return apolloClient.GetConfig(namespace)
}

func ApolloNamespaceValue[T cmp.Ordered | bool](namespace, key string) T {
	v := ApolloNamespace(namespace).GetValue(key)
	var t T
	core.ConvertStringTo(v, &t)
	return t
}
