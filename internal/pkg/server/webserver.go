package server

type WebServer interface {
	Listen(addr string)
	Post(path string, handler Handler)
	Get(path string, handler Handler)
}

type Request interface {
	GetBody(data any) (any, error)
}

type Handler func(body *FiberRequest) map[string]interface{}
