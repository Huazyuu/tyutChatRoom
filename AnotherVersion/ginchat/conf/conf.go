package conf

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var GlobalConf AppConf

type AppConf struct {
	App        App        `yaml:"app"`
	Mysql      Mysql      `yaml:"mysql"`
	Log        Log        `yaml:"log"`
	PictureBed PictureBed `yaml:"picture_bed"`
}

type App struct {
	Port           string `yaml:"port"`
	UploadFilePath string `yaml:"upload_file_path"`
	CookieKey      string `yaml:"cookie_key"`
	ServeType      string `yaml:"serve_type"`
	Mode           string `yaml:"mode"`
}
type Mysql struct {
	Dsn string `yaml:"dsn"`
}
type Log struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}

type PictureBed struct {
	Type   string `yaml:"type"`
	ApiKey string `yaml:"api_key"`
}

const ConfigFile = "./conf/settings.yaml"

func InitCore(path string) {
	if path == "" {
		path = ConfigFile
	}
	conf := &AppConf{}
	yamlConf, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error : %s", err))
	}
	err = yaml.Unmarshal(yamlConf, &conf)
	if err != nil {
		log.Fatalf("conf init Unmarshal: %v", err)
	}
	GlobalConf = *conf
}
