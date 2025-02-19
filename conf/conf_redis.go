package conf

import (
	"fmt"
)

type Redis struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	PoolSize int    `yaml:"pool_size"`
	DB       int    `yaml:"db"`
}

func (r Redis) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.IP, r.Port)
}
