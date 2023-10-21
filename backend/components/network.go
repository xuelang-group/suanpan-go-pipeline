package components

import (
	"encoding/json"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/web"
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio"
)

func getSio(uri string, path string, namespace string) (*socketio.ClientConn, error) {
	u, err := url.Parse(uri)
	if err != nil {
		logrus.Errorf("Parse url error: %v", err)
		return nil, err
	}
	schemeOpt := socketio.WithScheme("ws")
	if u.Scheme == "https" {
		schemeOpt = socketio.WithScheme("wss")
	}
	pathOpt := socketio.WithPath(path)

	// u = socketio.GetURL(u.Host, schemeOpt, pathOpt)

	// headerOpt := socketio.WithHeader(web.GetHeaders())
	// namespaceOpt := socketio.WithNamespace(namespace)

	u = socketio.GetURL(u.Host, schemeOpt, pathOpt)

	conn, err := socketio.NewClientConn(u.String(), &socketio.ClientOptions{
		Namespace:      namespace,
		Header:         web.GetHeaders(), // not working now
		Reconnect:      true,
		EventBufferMax: 1000,
	})

	return conn, err
}

func emitEvent(uri string, path string, namespace string, event string, data interface{}) {
	sio, err := getSio(uri, path, namespace)
	if err != nil {
		logrus.Errorf("Get sio error: %v", err)
		return
	}
	defer sio.Close()

	if data == nil {
		sio.Emit(event)
	} else {
		sio.Emit(event, data)
	}
}

func socketIOClientMain(currentNode Node, inputData RequestData) (map[string]interface{}, error) {
	url := currentNode.Config["url"].(string)
	path := currentNode.Config["path"].(string)
	namespace := currentNode.Config["namespace"].(string)
	event := currentNode.Config["event"].(string)
	data := currentNode.Config["data"].(string)
	if len(data) == 0 {
		emitEvent(url, path, namespace, event, nil)
	} else {
		var jsonData interface{}
		err := json.Unmarshal([]byte(data), &jsonData)
		if err != nil {
			log.Info("Can not convert data to json, send text...")
			emitEvent(url, path, namespace, event, data)
		} else {
			emitEvent(url, path, namespace, event, jsonData)
		}
	}
	return map[string]interface{}{"out1": "success"}, nil
}
