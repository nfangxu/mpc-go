package conf

import (
	"encoding/json"
	"os"
)

type Conf struct {
	Self    string      `json:"self"`
	Nodes   []Node      `json:"nodes"`
	Parties []Party     `json:"parties"`
	DSL     []Component `json:"dsl"`
	Runtime Runtime     `json:"runtime"`
}

func (c *Conf) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(c)
}

type Node struct {
	NodeId     string `json:"node_id"`
	PartyId    string `json:"party_id"`
	Address    string `json:"address"`
	ListenAddr string `json:"listen_addr"`
}

type Party struct {
	NodeId     string `json:"node_id"`
	PartyId    string `json:"party_id"`
	Address    string `json:"address"`
	ListenAddr string `json:"listen_addr"`
}

type Component struct {
	Name      string             `json:"name"`
	Input     map[string][]Input `json:"input"`
	Output    map[string]Output  `json:"output"`
	Parameter map[string]any     `json:"parameter"`
}

type Input struct {
	JobId     string   `json:"job_id"`
	TaskId    string   `json:"task_id"`
	PartyId   string   `json:"party_id"`
	Filename  Filename `json:"filename"`
	ContainsY bool     `json:"contains_y"`
	Y         string   `json:"y"`
}

type Output struct {
	Dataset string `json:"dataset"`
}

type Filename struct {
	Dataset string `json:"dataset"`
}

// Runtime 预留的一些其他配置
type Runtime struct {
}
