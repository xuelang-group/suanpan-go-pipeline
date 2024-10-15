package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

type RequestMsg struct {
	AppId    string `json:"appId"`
	NodeId   string `json:"nodeId"`
	UserId   string `json:"userId"`
	NodePort int    `json:"nodePort"`
	Port     int    `json:"port"`
}

func FindPort(host string, ports []string) string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ports), func(i, j int) { ports[i], ports[j] = ports[j], ports[i] })
	availablePort := ports[0]
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			log.Infof("Port %s is available.", port)
			availablePort = port
			break
		}
		if conn != nil {
			defer conn.Close()
			log.Errorf("Port %s is not available, try next port.", port)
		}
	}
	return availablePort
}

func MakeRange(min int, max int) []string {
	a := make([]string, max-min+1)
	for i := range a {
		a[i] = strconv.FormatInt(int64(min+i), 10)
	}
	return a
}

func RegisterPort(virtualport string, realport string) {
	e := config.GetEnv()
	protocol := "https"
	if e.SpHostTls == "false" {
		protocol = "http"
	}

	url := fmt.Sprintf("%s://localhost:%s/app/service/register", protocol, e.SpPort)

	virtualportInt, _ := strconv.ParseInt(virtualport, 10, 32)
	realportInt, _ := strconv.ParseInt(realport, 10, 32)
	param := RequestMsg{AppId: e.SpAppId, NodeId: e.SpNodeId, UserId: e.SpUserId, NodePort: int(virtualportInt), Port: int(realportInt)}
	jsonParam, _ := json.Marshal(param)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonParam))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
