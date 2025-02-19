package conf

import "fmt"

type Mysql struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Config          string `yaml:"config"` // 高级配置,如charset
	DB              string `yaml:"db"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	LogLevel        string `yaml:"log_level"` // 日志等级
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxConnLifeTime int    `yaml:"max_conn_lifetime"`
}

func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.User, m.Password, m.Host, m.Port, m.DB, m.Config)
}
