package components

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"github.com/go-gota/gota/dataframe"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type scriptData struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

func pyScriptMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	inputdata := getScriptInputData(currentNode)
	log.Infof("节点%s(%s)接收到参数 %s", currentNode.Key, currentNode.Id, inputdata)
	// inputsStringArr := make([]string, 0)
	// for _, inputString := range inputStrings {
	// 	inputsStringArr = append(inputsStringArr, inputString)
	// }
	// var inputdata = strings.Join(inputsStringArr, ",")
	var script = currentNode.Config["script"].(string)
	var nodeid = currentNode.Id
	params := url.Values{}

	Url, err := url.Parse("http://0.0.0.0:8080/data/?nodeid=10112&inputdata=fdsf&script=scr")
	if err != nil {
		log.Errorf("can not run script with error: %s", err.Error())
		return map[string]interface{}{}, nil
	}
	params.Set("nodeid", nodeid)
	params.Set("inputdata", inputdata)
	params.Set("script", script)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	defer resp.Body.Close()
	stdout, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(stdout))
	// cmdStrings := make([]string, 0)
	// cmdStrings = append(cmdStrings, "scripts/pyRuntime.py")
	// for _, inputString := range inputStrings {
	//         cmdStrings = append(cmdStrings, inputString)
	// }
	// cmdStrings = append(cmdStrings, "--script")
	// cmdStrings = append(cmdStrings, currentNode.Config["script"].(string))
	// cmd := exec.Command("python3", cmdStrings...)
	// stdout, err := cmd.Output()
	// if err != nil {
	//         log.Infof("can not run script with error: %s", err.Error())
	//         return map[string]interface{}{}, nil
	// }
	outs := []scriptData{}
	err1 := json.Unmarshal(stdout, &outs)
	if err1 != nil {
		log.Errorf("can not solve output data with error: %s", err.Error())
		return map[string]interface{}{}, nil
	}
	return getScriptOutputData(outs, currentNode), nil
}

// func getScriptInputData(currentNode Node) string {
// 	inputDatas := make([]scriptData, 0)
// 	log.Infof("当前节点%s,  %s", currentNode.Id, currentNode.InputData)
// 	for _, v := range currentNode.InputData {
// 		inputData := scriptData{}
// 		switch i := v.(type) {
// 		case dataframe.DataFrame:
// 			os.Mkdir(currentNode.Id, os.ModePerm)
// 			tmpPath := currentNode.Id + "/data.csv"
// 			os.Remove(tmpPath)
// 			file, err := os.Create(tmpPath)
// 			if err != nil {
// 				log.Error("无法创建临时文件")
// 			}
// 			i.WriteCSV(file)
// 			inputData.Data = tmpPath
// 			inputData.Type = "csv"
// 		default:
// 			inputData.Data = i
// 			inputData.Type = "json"
// 		}
// 		// inputString, _ := json.Marshal(inputData)
// 		inputDatas = append(inputDatas, inputData)
// 	}
// 	inputString, _ := json.Marshal(inputDatas)
// 	return string(inputString)
// }

func getScriptInputData(currentNode Node) string {
	keys := make([]string, 0, len(currentNode.InputData))
	for k, _ := range currentNode.InputData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	inputDatas := make([]scriptData, len(currentNode.InputData))
	for _, k := range keys {
		inputData := scriptData{}
		switch i := currentNode.InputData[k].(type) {
		case dataframe.DataFrame:
			os.Mkdir(currentNode.Id, os.ModePerm)
			tmpPath := currentNode.Id + "/data.csv"
			os.Remove(tmpPath)
			file, err := os.Create(tmpPath)
			if err != nil {
				log.Error("无法创建临时文件")
			}
			i.WriteCSV(file)
			inputData.Data = tmpPath
			inputData.Type = "csv"
		default:
			inputData.Data = i
			inputData.Type = "json"
		}
		// inputString, _ := json.Marshal(inputData)
		numStr := k[2:] // 从第3个字符开始截取
		if num, err := strconv.Atoi(numStr); err == nil {
			inputDatas[num-1] = inputData
		}
	}
	inputString, _ := json.Marshal(inputDatas)
	return string(inputString)
}

func readCsv(filePath string) dataframe.DataFrame {
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Errorf("Can not open csv file: %s, with error: %s", filePath, err.Error())
	}
	defer func() {
		csvFile.Close()
		err = os.Remove(filePath)
		if err != nil {
			log.Errorf("Can not remove csv file: %s, with error: %s", filePath, err.Error())
		}
	}()
	df := dataframe.ReadCSV(csvFile)
	return df
}

func getScriptOutputData(outputs []scriptData, currentNode Node) map[string]interface{} {
	outputDatas := map[string]interface{}{}
	idx := 0
	for port := range currentNode.OutputData {
		if len(outputs) >= idx+1 {
			switch outputs[idx].Type {
			case "csv":
				outputDatas[port] = readCsv(outputs[idx].Data.(string))
			case "json":
				outputDatas[port] = outputs[idx].Data
			}
		} else {
			break
		}
		idx += 1

	}
	return outputDatas
}
