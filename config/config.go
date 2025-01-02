package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

func ParseToml(filepath string, v any) {
	if _, err := toml.DecodeFile(filepath, &v); err != nil {
		logrus.Fatal(err)
	}
}
