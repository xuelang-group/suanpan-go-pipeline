package components

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"

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
	params := url.Values{}

	Url, err := url.Parse("http://0.0.0.0:8080/data/?nodeid=10112&inputdata=fdsf&script=scr")
	if err != nil {
		log.Infof("can not run script with error: %s", err.Error())
		return map[string]interface{}{}, nil
	}
	params.Set("nodeid", nodeid)
	params.Set("inputdata", inputdata)
	params.Set("script", script)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	log.Debugf("Python脚本编辑器(%s)调用python服务API: %s", currentNode.Id, urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		log.Infof("Python脚本编辑器(%s)调用python服务API报错: %s", currentNode.Id, err.Error())
		return map[string]interface{}{}, nil
	}
	defer resp.Body.Close()
	stdout, err := io.ReadAll(resp.Body)
	log.Debugf("Python脚本编辑器(%s)调用python服务API返回: %s", currentNode.Id, string(stdout))
	outs := []scriptData{}
	err1 := json.Unmarshal(stdout, &outs)
	if err1 != nil {
		log.Infof("can not solve output data with error: %s", err.Error())
		return map[string]interface{}{}, nil
	}
	return getScriptOutputData(outs, currentNode), nil
}

func getScriptInputData(currentNode Node) string {
	inputDatas := make([]scriptData, 0)
	inputPorts := make([]string, 0)
	for port := range currentNode.InputData {
		inputPorts = append(inputPorts, port)
	}
	sort.Slice(inputPorts, func(i, j int) bool {
		return inputPorts[i] < inputPorts[j]
	})
	for _, port := range inputPorts {
		v := currentNode.InputData[port]
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
		inputDatas = append(inputDatas, inputData)
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
