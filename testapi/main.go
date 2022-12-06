package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	params := url.Values{}

	Url, err := url.Parse("http://0.0.0.0:8080/data/?inputdata=fdsf&script=scr")
	if err != nil {
		panic(err.Error())
	}
	params.Set("inputdata", "{\"data\":\"1\",\"type\":\"json\"}"+","+"{\"data\":\"2\",\"type\":\"json\"}")
	params.Set("script", "outputs.append(int(inputs[0])+1)\noutputs.append(int(inputs[0])+2)")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(s))
	// cmdStrings := make([]string, 0)
	// cmdStrings = append(cmdStrings, "{\"data\":\"1\",\"type\":\"json\"}")
	// cmdStrings = append(cmdStrings, "outputs.append(int(inputs[0])+1)")
	// // res, err := http.Get("http://10.88.34.184:30080/proxr/100026/55182/9d40cc50743d11edad7e29ea3448dc30/8080/data/?inputdata=" + cmdStrings[0] + "&script=" + cmdStrings[0])
	// // res, err := http.Get("http://10.88.34.184:30080/proxr/100026/55182/9d40cc50743d11edad7e29ea3448dc30/8080/data/?inputdata=hello&script=world")
	// res, err := http.Get("http://0.0.0.0:8080/data/?inputdata=" + "{\"data\":\"1\",\"type\":\"json\"}" + "&script=" + "outputs.append(int(inputs[0])+1)")
	// if err != nil {
	// 	//...
	// }
	// defer res.Body.Close() //在回复后必须关闭回复的主体
	// body, err := ioutil.ReadAll(res.Body)
	// if err == nil {
	// 	fmt.Println(string(body))
	// }
}
