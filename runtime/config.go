package runtime

import (
	"io/ioutil"
	"path/filepath"

	"github.com/free5gc/openapi/models"
	"gopkg.in/yaml.v2"
)

var PopulateConfig Config

type Config struct {
	Mongo  Mongodb  `yaml:"mongo"`
	MCC    string   `yaml:"mcc"`
	MNC    string   `yaml:"mnc"`
	Key    string   `yaml:"key"`
	OP     string   `yaml:"op"`
	SQN    string   `yaml:"sqn"`
	AMF    string   `yaml:"amf"`
	Slices []Slice  `yaml:"slices"`
	IMSI   []string `yaml:"imsi"`
}

type Mongodb struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Slice struct {
	Sst   int32  `yaml:"sst"`
	Sd    string `yaml:"sd"`
	VarQI uint8  `yaml:"varqi"`
	Dnn   string `yaml:"dnn"`
	// IPV4, IPV6, IPV4V6, UNSTRUCTURED, ETHERNET. Default is IPV4
	PduSessionType *models.PduSessionType `yaml:"pdu-session-type,omitempty"`
}

// ParseConf read the yaml file and populate the Config instancce
func ParseConf(file string) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &PopulateConfig)
	if err != nil {
		return err
	}
	return nil
}
