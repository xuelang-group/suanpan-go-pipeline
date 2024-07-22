package components

import (
	"bytes"
	"goPipeline/utils"
	"goPipeline/variables"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	socketio "github.com/googollee/go-socket.io"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type NodeAction interface {
	Run(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	UpdateInput(inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	Main(inputData RequestData) (map[string]interface{}, error)
}

type Node struct {
	TriggeredPorts []string
	PreviousNodes  []*Node
	NextNodes      []*Node
	InputData      map[string]interface{}
	OutputData     map[string]interface{}
	PortConnects   map[string][]string
	Config         map[string]interface{}
	Id             string
	Key            string
	Run            func(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool, server *socketio.Server, runtimeErr *error)
	// dumpOutput    func(currentNode Node, outputData map[string]interface{})
	UpdateInput func(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool)
	loadInput   func(currentNode Node, inputData RequestData) error
	main        func(currentNode Node, inputData RequestData) (map[string]interface{}, error)
	initNode    func(currentNode Node) error
	releaseNode func(currentNode Node) error
	Status      int // 0: stoped 1： running -1：error
	// ServiceHandler services.Service
}

type RequestData struct {
	Data  string
	ID    string
	Extra string
}

func (c *Node) Init(nodeType string) {
	c.Run = Run
	c.UpdateInput = UpdateInput
	c.initNode = defaultInit
	c.releaseNode = defaultRelease
	// c.dumpOutput = dumpOutput
	switch nodeType {
	case "StreamIn":
		c.main = streamInMain
		c.loadInput = streamInLoadInput
	case "StreamOut":
		c.main = streamOutMain
	case "JsonExtractor":
		c.main = jsonExtractorMain
	case "String2Json":
		c.main = string2JsonMain
	case "DataSync":
		c.main = dataSyncMain
	case "GlobalVariableSetter":
		c.main = globalVariableSetterMain
	case "GlobalVariableGetter":
		c.main = globalVariableGetterMain
	case "GlobalVariableDeleter":
		c.main = globalVariablDeleterMain
	case "CsvDownloader":
		c.main = csvDownloaderMain
	case "CsvToDataFrame":
		c.main = CsvToDataFrameMain
	case "DataFrameToCsv":
		c.main = DataFrameToCsvMain
	case "ExecutePythonScript":
		c.main = pyScriptMain
	case "PostgresReader":
		c.main = postgresReaderMain
		c.initNode = postgresInit
	case "PostgresSqlExecuter":
		c.main = postgresExecutorMain
		c.initNode = postgresInit
	case "PostgresWriter":
		c.main = postgresWriterMain
		c.initNode = postgresInit
	case "OracleReader":
		c.main = oracleReaderMain
	case "OracleSqlExecutor":
		c.main = oracleExecutorMain
	case "OracleWriter":
		c.main = oracleWriterMain
	case "SQLServerReader":
		c.main = sqlServerReaderMain
	case "SQLServerSqlExecutor":
		c.main = sqlServerExecutorMain
	case "SQLServerWriter":
		c.main = sqlServerWriterMain
	case "HiveReader":
		c.main = hiveReaderMain
	case "HiveSqlExecutor":
		c.main = hiveExecutorMain
	case "HiveWriter":
		c.main = hiveWriterMain
	case "CsvUploader":
		c.main = csvUploaderMain
	case "SocketIOClient":
		c.main = socketIOClientMain
	case "Delay":
		c.main = dalayMain
	case "KafkaConsumer":
		c.main = kafkaConsumerMain
	case "KafkaProducer":
		c.main = kafkaProducerMain
	case "MysqlReader":
		c.main = mysqlReaderMain
		c.initNode = mysqlInit
		c.releaseNode = mysqlRlease
	case "MysqlExecutor":
		c.main = mysqlExecutorMain
		c.initNode = mysqlInit
		c.releaseNode = mysqlRlease
	case "MysqlWriter":
		c.main = mysqlWriterMain
		c.initNode = mysqlInit
		c.releaseNode = mysqlRlease
	case "MysqlJsonReader":
		c.main = mysqlJsonReaderMain
		c.initNode = mysqlInit
		c.releaseNode = mysqlRlease
	default:
	}
}

func (c *Node) Initialize() {
	c.initNode(*c)
}

func (c *Node) Release() {
	c.releaseNode(*c)
}

func Run(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool, server *socketio.Server, runtimeErr *error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("节点%s(%s)运行异常，错误日志：%s", currentNode.Key, currentNode.Id, err)
			*runtimeErr = err.(error)
		}
		wg.Done()
	}()
	select {
	case <-stopChan:
		log.Infof("节点%s(%s)运行被中断", currentNode.Key, currentNode.Id)
	default:
		receiveInputs := true
		for _, v := range currentNode.InputData {
			if v == nil {
				receiveInputs = false
				break
			}
		}
		if len(inputData.Data) > 0 || receiveInputs {
			currentNode.Status = 1
			outputData, err := currentNode.main(currentNode, inputData)
			if err != nil {
				*runtimeErr = err
				log.Debugf("节点%s(%s)运行失败: %s", currentNode.Key, currentNode.Id, err.Error())
				currentNode.Status = -1
				if server != nil {
					server.BroadcastToNamespace("/", "notify.process.status", map[string]int{currentNode.Id: -1})
					server.BroadcastToNamespace("/", "notify.process.error", map[string]string{currentNode.Id: err.Error()})
				}
			} else {
				log.Debugf("节点%s(%s)运行成功", currentNode.Key, currentNode.Id)
				readyToRun := make([]string, 0)
				triggeredPorts := make(map[string][]string)
				for port, data := range outputData { //map[out1:true]
					for _, tgt := range currentNode.PortConnects[port] {
						tgtInfo := strings.Split(tgt, "-")
						for i := range currentNode.NextNodes {
							if currentNode.NextNodes[i].Id == tgtInfo[0] {
								log.Debugf("数据下发到节点%s(%s)", currentNode.NextNodes[i].Key, currentNode.NextNodes[i].Id)
								tmpData := data
								if dataString, ok := tmpData.(string); ok {
									if strings.HasSuffix(dataString, ".csv") {
										basename := path.Base(dataString)
										dst := strings.Replace(dataString, basename, "data_"+currentNode.NextNodes[i].Id+".csv", -1)
										utils.CopyFile(dataString, dst)
										tmpData = dst
									}
								}
								currentNode.NextNodes[i].InputData[tgtInfo[1]] = tmpData
								triggeredPorts[currentNode.NextNodes[i].Id] = append(triggeredPorts[currentNode.NextNodes[i].Id], tgtInfo[1])
								if !utils.SlicesContain(readyToRun, currentNode.NextNodes[i].Id) {
									readyToRun = append(readyToRun, currentNode.NextNodes[i].Id)
								}
							}
						}
					}
					if dataString, ok := data.(string); ok {
						if strings.HasSuffix(dataString, ".csv") {
							_, err := os.Stat(dataString)
							if err == nil {
								err = os.Remove(dataString)
								if err != nil {
									log.Errorf("Can not remove csv file: %s, with error: %s", dataString, err.Error())
								}
							}
						}
					}
				}
				for i := range currentNode.NextNodes {
					if utils.SlicesContain(readyToRun, currentNode.NextNodes[i].Id) {
						currentNode.NextNodes[i].TriggeredPorts = triggeredPorts[currentNode.NextNodes[i].Id]
						wg.Add(1)
						go currentNode.NextNodes[i].Run(*currentNode.NextNodes[i], RequestData{ID: inputData.ID, Extra: inputData.Extra}, wg, stopChan, server, runtimeErr)
					}
				}
				currentNode.Status = 0
				if server != nil {
					server.BroadcastToNamespace("/", "notify.process.status", map[string]int{currentNode.Id: 0})
				}
			}
		}
	}
}

func UpdateInput(currentNode Node, inputData RequestData, wg *sync.WaitGroup, stopChan chan bool) {
	defer wg.Done()
	select {
	case <-stopChan:
		log.Info("Receive stop event")
	default:
		err := currentNode.loadInput(currentNode, inputData)
		if err != nil {
			log.Errorf("Error occur when running node: %s, error info: %s", currentNode.Key, err.Error())
		}
	}
}

// func dumpOutput(currentNode Node, outputData map[string]interface{}) {
// 	for port, data := range outputData { //map[out1:true]
// 		for _, tgt := range currentNode.PortConnects[port] {
// 			tgtInfo := strings.Split(tgt, "-")
// 			for _, node := range currentNode.NextNodes {
// 				if node.Id == tgtInfo[0] {
// 					node.InputData[tgtInfo[1]] = data
// 				}
// 			}
// 		}
// 	}
// }

func loadParameter(parameter string, vars map[string]interface{}) string {
	paramT := template.New("parameterLoader")
	paramT, err := paramT.Parse(parameter)
	if err != nil {
		log.Errorf("无法正常载入参数：%s", parameter)
		return parameter
	}
	var result bytes.Buffer
	for k, v := range variables.GlobalVariables {
		vars[k] = v
	}
	paramT.Execute(&result, vars)
	return result.String()
}

func defaultInit(currentNode Node) error {
	return nil
}

func defaultRelease(currentNode Node) error {
	return nil
}
