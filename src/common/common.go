package common

import (
	model2 "chat/src/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// 读取配置文件
func ReadConfig() model2.ConfigYaml {

	basePath, _ := os.Getwd()
	basePath = strings.Split(basePath, "chat")[0]
	configFilePath := path.Join(basePath, "chat", "src", "config", "config.yaml")
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Read config file fail: %s", err)
	}

	var ConfigYaml model2.ConfigYaml
	err = yaml.Unmarshal(yamlFile, &ConfigYaml)
	if err != nil {
		log.Fatalf("Unmarshal config file fail: %s", err)
	}

	return ConfigYaml
}
