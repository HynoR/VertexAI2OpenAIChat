package main

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"os"
)

var ConfigInstance Config

type Config struct {
	Location   string `json:"location" yaml:"location"`
	KeyFile    string `json:"keyfile" yaml:"keyfile"`
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

type KeyCredit struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func GetProjectId(filepath string) string {
	content, err := os.ReadFile(filepath)
	var k KeyCredit
	err = json.Unmarshal(content, &k)
	if err != nil {
		panic(err)
	}
	return k.ProjectId
}

func main() {
	ConfigInstance = ReadConfig()

	InitVertexInstance(GetProjectId(ConfigInstance.KeyFile), ConfigInstance.Location, ConfigInstance.KeyFile)
	r := NewRouter()
	r.Run(ConfigInstance.ListenAddr)

}
