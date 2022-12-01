package components

func jsonExtractorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func dataSyncMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	log.Info("start data sync !!!")
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
					result[strings.Replace(port, "in", "out", -1)] = data
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
					result[strings.Replace(port, "in", "out", -1)] = data
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
					result[strings.Replace(port, "in", "out", -1)] = data
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
					result[strings.Replace(port, "in", "out", -1)] = data
				}
			}
		}
	}
	// log.Infof("ly  datasync result %s ", result)
	return result, nil
}
