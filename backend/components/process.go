package components

func jsonExtractorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func dataSyncMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if currentNode.Config["triggerPort"].(string) == "" {
		if currentNode.Config["empty"].(bool) {
			recieveAll := true
			for _, data := range currentNode.InputData {
				if data == nil {
					recieveAll = false
				}
			}
			if recieveAll {
				for port, data := range currentNode.InputData {
					result[port] = data
				}
				for port := range currentNode.InputData {
					currentNode.InputData[port] = nil
				}
			}
		} else {
			recieveAll := true
			for _, data := range currentNode.InputData {
				if data == nil {
					recieveAll = false
				}
			}
			if recieveAll {
				for port, data := range currentNode.InputData {
					result[port] = data
				}
			}
		}
	} else {
		if currentNode.Config["empty"].(bool) {
			recieveAll := true
			for _, data := range currentNode.InputData {
				if data == nil {
					recieveAll = false
				}
			}
			if recieveAll && currentNode.InputData[currentNode.Config["triggerPort"].(string)] != nil {
				for port, data := range currentNode.InputData {
					result[port] = data
				}
				for port := range currentNode.InputData {
					currentNode.InputData[port] = nil
				}
			}
		} else {
			recieveAll := true
			for _, data := range currentNode.InputData {
				if data == nil {
					recieveAll = false
				}
			}
			if recieveAll && currentNode.InputData[currentNode.Config["triggerPort"].(string)] != nil {
				for port, data := range currentNode.InputData {
					result[port] = data
				}
			}
		}
	}

	return result, nil
}
