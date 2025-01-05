package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/bucketheadv/infragin"
	"github.com/bucketheadv/infragin/components/apollo"
	"github.com/bucketheadv/infragin/components/rocket"
	"github.com/bucketheadv/infragin/components/xxljob"
	"github.com/bucketheadv/infragin/db"
	"github.com/go-redis/redis/v8"
)

type Conf struct {
	Server   infragin.ServerConf      `json:"Server"`
	Apollo   apollo.Conf              `json:"Apollo"`
	XxlJob   xxljob.Conf              `json:"XxlJob"`
	MySql    map[string]db.MySqlConf  `json:"MySQL"`
	Redis    map[string]redis.Options `json:"Redis"`
	RocketMQ map[string]rocket.Conf   `json:"RocketMQ"`
}

func Parse(configFile string, config *Conf) error {
	if _, err := toml.DecodeFile(configFile, config); err != nil {
		return err
	}
	return nil
}
