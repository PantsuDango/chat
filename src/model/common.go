package model

type ConfigYaml struct {
	Mysql  ConfigMysql `yaml:"Mysql"`
	Server Server      `yaml:"Server"`
}

type ConfigMysql struct {
	IP       string `yaml:"IP"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

type Server struct {
	Port string `yaml:"Port"`
}
