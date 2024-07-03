package main

import (
	"context"
	"gopkg.in/yaml.v2"
	"os"
)

var ConfigInstance Config

type Config struct {
	Location   string `json:"location" yaml:"location"`
	KeyFile    string `json:"keyfile" yaml:"keyfile"`
	ProjectID  string `json:"projectid" yaml:"projectid"`
	AuthKey    string `json:"authkey" yaml:"authkey"`
	ListenAddr string `json:"listenaddr" yaml:"listenaddr"`
}

func ReadConfig() Config {
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var yamlConfig Config
	err = yaml.Unmarshal(content, &yamlConfig)
	if err != nil {
		panic(err)
	}
	return yamlConfig
}

func main() {
	ConfigInstance = ReadConfig()

	ctx := context.Background()
	v := NewVertexInstance(ConfigInstance.ProjectID, ConfigInstance.Location, ConfigInstance.KeyFile).VerTexAuth(ctx)
	VertexIns = *v
	r := NewRouter()
	r.Run(ConfigInstance.ListenAddr)

}
