package config

import (
	"encoding/json"
	"io/ioutil"
)

type GraphConfig struct {
	Topology map[string]*NodeConfig `json:"topology"`
}

type NodeConfig struct {
	Inputs  map[string]int `json:"inputs,omitempty"`
	Outputs map[string]int `json:"outputs,omitempty"`
}

// ParseGraphConfig parses the JSON file.
func ParseGraphConfig(filePath string) (*GraphConfig, error) {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		return nil, e
	}
	var jsonObj GraphConfig
	e = json.Unmarshal(file, &jsonObj)
	if e != nil {
		return nil, e
	}
	return &jsonObj, nil
}
