package components

import (
	"goPipeline/utils"
	"goPipeline/variables"
	"strconv"
	"strings"
	"time"
)

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
			if recieveAll && utils.SlicesContain(currentNode.TriggeredPorts, currentNode.Config["triggerPort"].(string)) {
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
			if recieveAll && utils.SlicesContain(currentNode.TriggeredPorts, currentNode.Config["triggerPort"].(string)) {
				for port, data := range currentNode.InputData {
					result[strings.Replace(port, "in", "out", -1)] = data
				}
			}
		}
	}
	return result, nil
}

func globalVariableSetterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	varname := currentNode.Config["name"].(string)
	variables.GlobalVariables[varname] = currentNode.InputData["in1"]
	return map[string]interface{}{"out1": "success"}, nil
}

func globalVariableGetterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	varname := currentNode.Config["name"].(string)
	if val, ok := variables.GlobalVariables[varname]; ok {
		return map[string]interface{}{"out1": val}, nil
	} else {
		return map[string]interface{}{}, nil
	}
}

func globalVariablDeleterMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	varname := currentNode.Config["name"].(string)
	if _, ok := variables.GlobalVariables[varname]; ok {
		delete(variables.GlobalVariables, varname)
		return map[string]interface{}{"out1": "success"}, nil
	} else {
		return map[string]interface{}{}, nil
	}
}

func dalayMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	duration := currentNode.Config["duration"].(string)
	intDuration, _ := strconv.Atoi(duration)
	time.Sleep(time.Duration(intDuration) * time.Second)
	result := make(map[string]interface{})
	for port, data := range currentNode.InputData {
		result[strings.Replace(port, "in", "out", -1)] = data
	}
	for port := range currentNode.InputData {
		currentNode.InputData[port] = nil
	}
	return result, nil
}
