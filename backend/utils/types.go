package utils

type Data struct {
	Value string
	Type  string
}

type Graph struct {
}

type Component struct {
	Type          string      `json:"type"`
	TypeLabel     string      `json:"typeLabel"`
	Category      string      `json:"category"`
	CategoryLabel string      `json:"categoryLabel"`
	Name          string      `json:"name"`
	Key           string      `json:"key"`
	HelpUrl       string      `json:"helpUrl"`
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
