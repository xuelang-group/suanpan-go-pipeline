package web

import (
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

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Infof("connected: %s", s.ID())
		return nil
	})

	server.OnEvent("/", "get.fields", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "set.fields", func(s socketio.Conn, msg []string) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "get.column", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "set.config", func(s socketio.Conn, msg ColorConfigType) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "get.config", func(s socketio.Conn, msg interface{}) RespondMsg {
		return RespondMsg{true, []string{}}
	})

	server.OnEvent("/", "clear.config", func(s socketio.Conn, msg interface{}) RespondMsg {
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
