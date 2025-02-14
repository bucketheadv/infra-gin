package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/bucketheadv/infra-gin/components/apollo"
	"github.com/bucketheadv/infra-gin/components/rocket"
	"github.com/bucketheadv/infra-gin/components/xxljob"
	"github.com/bucketheadv/infra-gin/db"
	"github.com/go-redis/redis/v8"
)

type ServerConf struct {
	Port int
}

type Conf struct {
	Server   ServerConf
	Apollo   apollo.Conf
	XxlJob   xxljob.Conf
	MySQL    map[string]db.MySQLConf
	Redis    map[string]redis.Options
	RocketMQ map[string]rocket.Conf
}

func Parse(configFile string, config *Conf) error {
	if _, err := toml.DecodeFile(configFile, config); err != nil {
		return err
	}
	return nil
}
