package components

import (
	"encoding/json"
	"goPipeline/utils"
	"goPipeline/variables"
	"strconv"
	"strings"
	"time"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

func string2JsonMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	inputString := currentNode.InputData["in1"].(string)
	inputJson := make(map[string]interface{})
	err := json.Unmarshal([]byte(inputString), &inputJson)
	if err != nil {
		log.Infof("Json解析失败: %s", err)
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": inputJson}, nil
}

func jsonExtractorMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	inputString := currentNode.InputData["in1"].(string)
	inputJson := make(map[string]interface{})
	err := json.Unmarshal([]byte(inputString), &inputJson)
	if err != nil {
		log.Infof("Json解析失败: %s", err)
		return map[string]interface{}{}, nil
	}
	result := loadParameter(currentNode.Config["param1"].(string), inputJson)
	return map[string]interface{}{"out1": result}, nil
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
