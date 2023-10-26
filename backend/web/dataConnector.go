package web

import (
	"encoding/json"
	"goPipeline/graph"
	"goPipeline/services"
	"goPipeline/utils"
	"goPipeline/variables"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

func sendToStream() {

	id := util.GenerateUUID()
	extra := ""
	r := stream.Request{ID: id, Extra: extra}
	r.Send(map[string]string{
		"out1": "",
	})
}

func variable(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	if len(name) == 0 {
		respond := RespondMsg{Success: false, Data: nil}
		jsonResp, err := json.Marshal(respond)
		if err != nil {
			log.Errorf("Error happened in JSON marshal. Err: %s", err.Error())
		}
		w.Write(jsonResp)
		return
	}
	switch req.Method {
	case "GET":
		if val, ok := variables.GlobalVariables[name]; ok {
			respond := RespondMsg{Success: true, Data: val}
			jsonResp, err := json.Marshal(respond)
			if err != nil {
				log.Errorf("Error happened in JSON marshal. Err: %s", err.Error())
			}
			w.Write(jsonResp)
			return
		} else {
			respond := RespondMsg{Success: true, Data: nil}
			jsonResp, err := json.Marshal(respond)
			if err != nil {
				log.Errorf("Error happened in JSON marshal. Err: %s", err.Error())
			}
			w.Write(jsonResp)
			return
		}
	case "POST":
		var data interface{}
		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			log.Errorf("Error happened in JSON decoding. Err: %s", err.Error())
			respond := RespondMsg{Success: false, Data: nil}
			jsonResp, _ := json.Marshal(respond)
			w.Write(jsonResp)
			return
		}
		variables.GlobalVariables[name] = data
		respond := RespondMsg{Success: true, Data: nil}
		jsonResp, _ := json.Marshal(respond)
		w.Write(jsonResp)
		return
	case "DELETE":
		delete(variables.GlobalVariables, name)
		respond := RespondMsg{Success: true, Data: nil}
		jsonResp, _ := json.Marshal(respond)
		w.Write(jsonResp)
		return
	default:
	}

}

func RunWeb(appType string) {

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Infof("connected: %s", s.ID())
		return nil
	})

	server.OnEvent("/", "components.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, graph.GraphInst.Components}
	})

	server.OnEvent("/", "graph.update", func(s socketio.Conn, msg utils.GraphConfig) RespondMsg {
		graph.GraphInst.Update(msg)
		services.ServicesManager.Update(msg)
		return RespondMsg{true, graph.GraphInst.Config}
	})

	server.OnEvent("/", "graph.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, graph.GraphInst.Config}
	})

	server.OnEvent("/", "process.run", func(s socketio.Conn, msg interface{}) RespondMsg {
		//传入单步运行的模式，进行运行
		id := util.GenerateUUID()
		go graph.GraphInst.Run(map[string]string{}, id, "", server, true)
		return RespondMsg{true, nil}
	})

	server.OnEvent("/", "process.stop", func(s socketio.Conn, msg interface{}) RespondMsg {
		graph.GraphInst.Stop()
		return RespondMsg{true, nil}
	})

	server.OnEvent("/", "process.status.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		nodeStatus := make(map[string]int)
		for _, node := range graph.GraphInst.Nodes {
			nodeStatus[node.Id] = node.Status
		}
		return RespondMsg{true, map[string]interface{}{"status": graph.GraphInst.PipelineStatus, "nodes": nodeStatus}}
	})

	server.OnEvent("/", "graph.status.set", func(s socketio.Conn, msg map[string]interface{}) RespondMsg {
		graph.GraphInst.Status = uint(msg["status"].(float64))
		if graph.GraphInst.Status == 0 {
			services.ServicesManager.Release()
			graph.GraphInst.Release()
		} else {
			services.ServicesManager.Deploy(&graph.GraphInst)
			graph.GraphInst.Initialize()
		}
		return RespondMsg{true, graph.GraphInst.Status}
	})

	server.OnEvent("/", "graph.status.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, graph.GraphInst.Status}
	})

	server.OnEvent("/", "result.visualize", func(s socketio.Conn, msg string) RespondMsg {
		for _, node := range graph.GraphInst.Nodes {
			if node.Id == msg {
				return RespondMsg{true, node.OutputData}
			}
		}
		return RespondMsg{true, ""}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Infof("meet error: %s", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Infof("closed %s", msg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Infof("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("statics")))
	http.HandleFunc("/variable", variable)

	http.ListenAndServe("0.0.0.0:"+WebServerPort, nil)
}
