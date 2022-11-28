package components

func jsonExtractorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	if m, ok := currentNode.InputData["in1"].(map[string]interface{}); !ok {
		return map[string]interface{}{}, nil
	} else if _, ok := m[currentNode.Config["field"].(string)]; ok {
		return map[string]interface{}{"out1": m[currentNode.Config["field"].(string)]}, nil
	}
	return map[string]interface{}{}, nil
}
