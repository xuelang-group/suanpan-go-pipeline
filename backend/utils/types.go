package utils

type Data struct {
	Value string
	Type  string
}

type Graph struct {
}

type Component struct {
	Type          string            `json:"type"`
	TypeLabel     string            `json:"typeLabel"`
	Category      string            `json:"category"`
	CategoryLabel string            `json:"categoryLabel"`
	Name          string            `json:"name"`
	Key           string            `json:"key"`
	HelpUrl       string            `json:"helpUrl"`
	Parameters    []Parameter       `json:"parameters"`
	Ports         map[string][]Port `json:"ports"`
}

type Parameter struct {
	Key      string
	Name     string
	Type     string
	Required bool
}

type Port struct {
	Id   string
	Name string
}
