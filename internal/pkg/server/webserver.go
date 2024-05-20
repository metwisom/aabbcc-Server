package server

type WebServer interface {
	Listen(addr string)
	Close()
	Post(path string, handler Handler)
	Get(path string, handler Handler)
}

type Request interface {
	GetBody(data any) map[string]interface{}
	GetQuery(data any) map[string]interface{}
}

type Handler func(body Request) map[string]interface{}
