package stream

import (
	"encoding/json"
	"fmt"
	"goPipeline/web"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type DataConnectorComponent struct {
	DefaultComponents
}

func (c *DataConnectorComponent) CallHandler(r stream.Request) {
	inputData := r.InputData(1)
	log.Info("receive request data: " + inputData)
	e := config.GetEnv()
	args := config.GetArgs()
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], inputData, "data.csv")
	tmpKey := path.Join(inputData, "data.csv")
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	mioerr := storage.FGetObject(tmpKey, tmpPath)
	if mioerr != nil {
		log.Errorf("Can not download file: %s, with error: %s", tmpKey, mioerr.Error())
		r.Send(map[string]string{
			"out1": inputData,
		})
		return
	}
	csvFile, err := os.Open(tmpPath)
	if err != nil {
		log.Errorf("Can not open csv file: %s, with error: %s", tmpPath, err.Error())
		r.Send(map[string]string{
			"out1": inputData,
		})
		return
	}
	defer func() {
		csvFile.Close()
		err = os.Remove(tmpPath)
		if err != nil {
			log.Errorf("Can not remove csv file: %s, with error: %s", tmpPath, err.Error())
		}
	}()
	web.OriginalDF = dataframe.ReadCSV(csvFile)
	if (len(web.ColorConfig.ColorMapping) == 0) || (len(web.ColorConfig.Fields) == 0) {
		log.Info("Can not find any mapping rules")
		r.Send(map[string]string{
			"out1": inputData,
		})
		return
	}
	colorSlice := make([]string, 0, web.OriginalDF.Nrow())
	if len(web.ColorConfig.Fields) == 1 {
		for _, data := range web.OriginalDF.Col(web.ColorConfig.Fields[0]).Records() {
			colorSlice = append(colorSlice, web.ColorConfig.ColorMapping[data])
		}
	} else {
		for i := 0; i < web.OriginalDF.Nrow(); i++ {
			elems := []string{}
			for _, col := range web.ColorConfig.Fields {
				elems = append(elems, web.OriginalDF.Col(col).Records()[i])
			}
			colorSlice = append(colorSlice, web.ColorConfig.ColorMapping[strings.Join(elems, "_")])
		}
	}
	colorSeries := dataframe.New(series.New(colorSlice, series.String, "Xuelang@Des@Color"))
	processedDF := web.OriginalDF.CBind(colorSeries)
	writeFile, writeErr := os.Create("data.csv")
	if writeErr != nil {
		log.Errorf("Can not open csv file: data.csv, with error: %s", writeErr.Error())
		r.Send(map[string]string{
			"out1": inputData,
		})
		return
	}
	defer func() {
		writeFile.Close()
		err = os.Remove("data.csv")
		if err != nil {
			log.Errorf("Can not remove csv file: data.csv, with error: %s", err.Error())
		}
	}()
	wErr := processedDF.WriteCSV(writeFile)
	if wErr != nil {
		log.Errorf("Can not write file: data.csv, with error: %s", wErr.Error())
		r.Send(map[string]string{
			"out1": inputData,
		})
		return
	}
	uploadDir := path.Join("studio", e.SpUserId, "tmp", e.SpAppId, strings.Join(strings.Split(r.ID, "-"), ""), e.SpNodeId, "out1")
	storage.FPutObject(path.Join(uploadDir, "data.csv"), "data.csv")
	r.Send(map[string]string{
		"out1": uploadDir,
	})
}

func (c *DataConnectorComponent) InitHandler() {
	// 获取环境变量
	e := config.GetEnv()
	// 获取命令行参数
	args := config.GetArgs()
	// web配置文件，云端路径
	web.DataKey = strings.Join([]string{"studio", e.SpUserId, "configs", e.SpAppId, e.SpNodeId, "params.json"}, "/")
	// web配置文件，本地路径
	web.DataPath = path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], "studio", e.SpUserId, "configs", e.SpAppId, e.SpNodeId, "params.json")
	os.MkdirAll(filepath.Dir(web.DataPath), os.ModePerm)
	log.Info(fmt.Sprintf("download file from %s to %s", web.DataKey, web.DataPath))
	err := storage.FGetObject(web.DataKey, web.DataPath)
	if err != nil {
		log.Info("Fail to Load Config File, init with default value...")
		web.ColorConfig = web.ColorConfigType{Fields: []string{}, ColorMapping: map[string]string{}}
	} else {
		jsonFile, err := os.Open(web.DataPath)
		if err != nil {
			log.Info(err.Error())
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &web.ColorConfig)
		log.Info(fmt.Sprintf("Successfully Loaded Config File %s.", web.DataPath))
	}
}

func (c *DataConnectorComponent) SioHandler() {
	go web.RunWeb()
}
