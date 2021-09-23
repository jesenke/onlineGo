package memory

import (
	"net/rpc"
)

type Result struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data *interface{} `json:"data"`
}

type Param struct {
	Param []byte `json:"code"`
}

//service.Method
func Call(domain string, serviceMethod string, param Param) Result {
	client, _ := rpc.DialHTTP("tcp", domain)
	var result Result
	err := client.Call(serviceMethod, param, &result)
	if err != nil {
		result.Code = 401
		result.Msg = err.Error()
		return result
	} else {
		return result
	}
}
