package web

import (
	"goPipeline/graph"
	"goPipeline/utils"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
	"github.com/xuelang-group/suanpan-go-sdk/util"
)

var GraphConfig interface{}

type ColorConfigType struct {
	Fields       []string          `json:"fields"`
	ColorMapping map[string]string `json:"colorMapping"`
}

func sendToStream() {

	id := util.GenerateUUID()
	extra := ""
	r := stream.Request{ID: id, Extra: extra}
	r.Send(map[string]string{
		"out1": "",
	})
}

func RunWeb() {

	graph.GraphInst.LoadComponents()

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Infof("connected: %s", s.ID())
		return nil
	})

	server.OnEvent("/", "components.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, graph.GraphInst.Components}
	})

	server.OnEvent("/", "graph.update", func(s socketio.Conn, msg utils.GraphConfig) RespondMsg {
		graph.GraphInst.Config = msg
		return RespondMsg{true, msg}
	})

	server.OnEvent("/", "graph.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, graph.GraphInst.Config}
	})

	server.OnEvent("/", "process.run", func(s socketio.Conn, msg ColorConfigType) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "process.status.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "graph.status.set", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "graph.status.get", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "result.visualize.json", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "result.download", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "result.visualize", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
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

	http.ListenAndServe("0.0.0.0:"+WebServerPort, nil)
}
