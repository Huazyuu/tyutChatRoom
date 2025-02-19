package core

import (
	"fmt"
	"gin-gorilla/conf"
	"gin-gorilla/global"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const ConfigFile = "./conf/settings.yaml"

func InitCore() {
	c := &conf.Config{}
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error : %s", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("conf init Unmarshal: %v", err)
	}
	global.Config = c
	logrus.Info("conf init success")
}
