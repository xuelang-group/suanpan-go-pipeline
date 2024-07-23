package components

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type scriptData struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

func pyScriptMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	inputdata := getScriptInputData(currentNode)
	var script = currentNode.Config["script"].(string)
	var nodeid = currentNode.Id
	// params := url.Values{}

	Url := "http://0.0.0.0:8080/data/"
	// if err != nil {
	// 	log.Errorf("can not run script with error: %s", err.Error())
	// 	return map[string]interface{}{}, err
	// }
	// params.Set("nodeid", nodeid)
	// params.Set("inputdata", inputdata)
	// params.Set("script", script)
	// params.Set("messageid", inputData.ID)
	// params.Set("extra", inputData.Extra)
	payloads := map[string]string{"nodeid": nodeid, "inputdata": inputdata, "script": script, "messageid": inputData.ID, "extra": inputData.Extra}
	jsonStr, _ := json.Marshal(payloads)
	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Errorf("Python脚本编辑器(%s)调用python服务API报错: %s", currentNode.Id, err.Error())
		return map[string]interface{}{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	//如果参数中有中文参数,这个方法会进行URLEncode
	// Url.RawQuery = params.Encode()
	// urlPath := Url.String()
	log.Debugf("Python脚本编辑器(%s)调用python服务API: %s", currentNode.Id, Url)
	log.Debugf("Python脚本编辑器(%s)调用python服务API参数: %s", currentNode.Id, string(jsonStr))
	// resp, err := http.Post(urlPath)
	if err != nil {
		log.Errorf("Python脚本编辑器(%s)调用python服务API报错: %s", currentNode.Id, err.Error())
		return map[string]interface{}{}, err
	}
	if resp.Status == "400 Bad Request" {
		defer resp.Body.Close()
		stdout, _ := io.ReadAll(resp.Body)
		log.Errorf("Python脚本编辑器(%s)调用python服务API报错: %s", currentNode.Id, string(stdout))
		return map[string]interface{}{}, errors.New(string(stdout))
	}
	defer resp.Body.Close()
	stdout, err := io.ReadAll(resp.Body)
	log.Infof("Python脚本编辑器(%s)调用python服务API返回: %s", currentNode.Id, string(stdout))
	outs := map[string]scriptData{}
	err1 := json.Unmarshal(stdout, &outs)
	if err1 != nil {
		log.Errorf("Python脚本编辑器(%s)调用python服务API返回结果无法解析: %s", currentNode.Id, err.Error())
		return map[string]interface{}{}, err1
	}
	return getScriptOutputData(outs, currentNode), nil
}

func getScriptInputData(currentNode Node) string {
	inputDatas := map[string]scriptData{}
	// inputPorts := make([]string, 0)
	// for port := range currentNode.InputData {
	// 	inputPorts = append(inputPorts, port)
	// }
	// sort.Slice(inputPorts, func(i, j int) bool {
	// 	return inputPorts[i] < inputPorts[j]
	// })
	for port, v := range currentNode.InputData {
		// v := currentNode.InputData[port]
		inputData := scriptData{}
		switch i := v.(type) {
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
		// inputDatas = append(inputDatas, inputData)
		inputDatas[port] = inputData
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

func getScriptOutputData(outputs map[string]scriptData, currentNode Node) map[string]interface{} {
	outputDatas := map[string]interface{}{}
	for port := range currentNode.OutputData {
		switch outputs[port].Type {
		case "csv":
			outputDatas[port] = readCsv(outputs[port].Data.(string))
		case "json":
			outputDatas[port] = outputs[port].Data
		}

	}
	return outputDatas
}
