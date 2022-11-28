package utils

type Data struct {
	Value string
	Type  string
}

type GraphConfig struct {
	Htime      uint         `json:"htime,omitempty"`
	Scale      float32      `json:"scale"`
	X          float64      `json:"x"`
	Y          float64      `json:"y"`
	Nodes      []NodeConfig `json:"nodes"`
	Connectors []Connector  `json:"connectors"`
}

type NodeConfig struct {
	Uuid       string          `json:"uuid"`
	Puuit      string          `json:"puuid"`
	Type       string          `json:"type"`
	Key        string          `json:"key"`
	Name       string          `json:"name"`
	Status     int             `json:"status"`
	X          float64         `json:"x"`
	Y          float64         `json:"y"`
	Parameters []NodeParameter `json:"parameters"`
}

type NodeParameter struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Connector struct {
	Src map[string]string `json:"src"`
	Tgt map[string]string `json:"tgt"`
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
	Category      string      `json:"category,omitempty"`
	CategoryLabel string      `json:"categoryLabel,omitempty" yaml:"categoryLabel"`
	Name          string      `json:"name"`
	Key           string      `json:"key"`
	HelpUrl       string      `json:"helpUrl" yaml:"helpUrl"`
	Parameters    []Parameter `json:"parameters"`
	Ports         Ports       `json:"ports"`
}

type Parameter struct {
	Key      string      `json:"key"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Required bool        `json:"required,omitempty"`
	Default  interface{} `json:"default,omitempty"`
	Options  []Option    `json:"options,omitempty"`
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
type Ports struct {
	In  []Port `json:"in,omitempty"`
	Out []Port `json:"out,omitempty"`
}

type Port struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
