package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"msg"`
}

type Handle func(ctx *gin.Context) *Response

type Server struct {
	sv   *gin.Engine
	port string
	key  string
	cert string
}

func NewSimpleServer(port, key, cert string) *Server {
	return &Server{
		sv:   gin.New(),
		cert: cert,
		key:  key,
		port: port,
	}
}

func (s *Server) Start() error {
	if len(s.cert) > 0 && len(s.key) > 0 {
		return s.sv.RunTLS(s.port, s.cert, s.key)
	}
	return s.sv.Run(s.port)
}

func (s *Server) AddWrap(handle ...Handle) {
	for _, v := range handle {
		filter := func(context *gin.Context) {
			response := v(context)
			if response != nil {
				returnData(context, response.Code, response.Msg, nil)
			}
		}
		s.sv.Use(filter)
	}

}

func NewResponse(code, msg string, data *interface{}) *Response {
	return &Response{code, msg, data}
}

func returnData(context *gin.Context, code, msg string, data *interface{}) {
	DataResponse := Response{code, msg, data}
	context.JSON(http.StatusOK, DataResponse)
}

func (s *Server) AddRoute(method, path string, handle Handle) {
	filter := func(context *gin.Context) {
		response := handle(context)
		returnData(context, response.Code, response.Msg, &response.Data)
	}

	switch method {
	case http.MethodGet:
		s.sv.GET(path, filter)
	case http.MethodPost:
		s.sv.POST(path, filter)
	case http.MethodPut:
		s.sv.PUT(path, filter)
	case http.MethodDelete:
		s.sv.DELETE(path, filter)
	default:
		panic("不支持")
	}
}
