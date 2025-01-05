package apollo

import (
	"cmp"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/bucketheadv/infracore"
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

type Conf struct {
	Enabled        bool
	AppID          string
	Cluster        string
	NamespaceName  string
	IP             string
	IsBackupConfig bool
}

var apolloClient agollo.Client

func InitClient(c Conf, onSuccess func()) {
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

func AssignConfigValueTo[T cmp.Ordered | bool](namespace, key string, value *T) {
	var s = Namespace(namespace).GetValue(key)
	if s == "" {
		return
	}
	infracore.ConvertStringTo(s, value)
}

func ApplicationValue(key string) string {
	return apolloClient.GetValue(key)
}

func Namespace(namespace string) *storage.Config {
	return apolloClient.GetConfig(namespace)
}

func NamespaceValue[T cmp.Ordered | bool](namespace, key string) T {
	v := Namespace(namespace).GetValue(key)
	var t T
	infracore.ConvertStringTo(v, &t)
	return t
}
