package utils

type Data struct {
	Value string
	Type  string
}

type GraphConfig struct {
	htime      uint
	scale      float32
	x          float64
	y          float64
	Nodes      []NodeConfig
	Connectors []Connector
}

type NodeConfig struct {
	Uuid       string
	Puuit      string
	Type       string
	Key        string
	Name       string
	Status     int
	x          float64
	y          float64
	Parameters []map[string]string
}

type Connector struct {
	Src map[string]string
	Tgt map[string]string
}

type NodeInfo struct {
	Info    map[string]string       `json:"info"`
	Inputs  map[string]NodePortType `json:"inputs"`
	Outputs map[string]NodePortType `json:"outputs"`
	Params  map[string]NodePortType `json:"params"`
}

type NodePortType struct {
	Uuid        string            `json:"uuid"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Subtype     string            `json:"subtype,omitempty"`
	Description map[string]string `json:"description,omitempty"`
}

type Component struct {
	Type          string      `json:"type"`
	TypeLabel     string      `json:"typeLabel" yaml:"typeLabel"`
	Category      string      `json:"category"`
	CategoryLabel string      `json:"categoryLabel" yaml:"categoryLabel"`
	Name          string      `json:"name"`
	Key           string      `json:"key"`
	HelpUrl       string      `json:"helpUrl" yaml:"helpUrl"`
	Parameters    []Parameter `json:"parameters,omitempty"`
	Ports         Ports       `json:"ports"`
}

type Parameter struct {
	Key      string      `json:"key"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Required bool        `json:"required,omitempty"`
	Default  interface{} `json:"default,omitempty"`
}

type Ports struct {
	In  []Port `json:"in,omitempty"`
	Out []Port `json:"out,omitempty"`
}

type Port struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
