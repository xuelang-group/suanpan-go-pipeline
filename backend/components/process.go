package components

type JsonExtractor struct {
	Node
}

func (c *JsonExtractor) Main(inputData RequestData) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
