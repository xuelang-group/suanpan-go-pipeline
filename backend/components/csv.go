package components

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/storage"
)

func csvDownloaderMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	needBasename := currentNode.Config["needBasename"].(bool)
	args := config.GetArgs()
	tmpPath := path.Join(args[fmt.Sprintf("--storage-%s-temp-store", args["--storage-type"])], currentNode.InputData["in1"].(string), currentNode.Id, "data.csv")
	tmpKey := currentNode.InputData["in1"].(string)
	if needBasename {
		tmpKey = path.Join(currentNode.InputData["in1"].(string), "data.csv")
	}
	os.MkdirAll(filepath.Dir(tmpPath), os.ModePerm)
	storageErr := storage.FGetObject(tmpKey, tmpPath)
	if storageErr != nil {
		log.Errorf("Can not download file: %s, with error: %s", tmpKey, storageErr.Error())
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{"out1": tmpPath}, nil
}
