package stream

import (
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type Component interface {
	InitHandler()
	CallHandler(stream.Request)
	SioHandler()
}

type DefaultComponents struct {
}

func (c DefaultComponents) InitHandler() {
	log.Info("Init function not implement.")
}

func (c DefaultComponents) CallHandler(r stream.Request) {
	log.Info("Stream call function not implement.")
}

func (c DefaultComponents) SioHandler() {
	log.Info("Socket io server not implement.")
}
