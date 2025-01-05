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
	Server   infragin.ServerConf
	Apollo   apollo.Conf
	XxlJob   xxljob.Conf
	MySql    map[string]*db.MySqlConf
	Redis    map[string]*redis.Options
	RocketMQ map[string]*rocket.Conf
}

func Parse(configFile string, config *Conf) error {
	if _, err := toml.DecodeFile(configFile, config); err != nil {
		return err
	}
	return nil
}
