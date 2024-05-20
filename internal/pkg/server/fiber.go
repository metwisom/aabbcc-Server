package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"reflect"
	"time"
)

type FiberWebServer struct {
	fiber *fiber.App
}

type FiberRequest struct {
	*fiber.Ctx
}

func ErrorNotFound(message string) map[string]interface{} {
	return map[string]interface{}{"code": 404, "error": message}
}
func ErrorConflict(message string) map[string]interface{} {
	return map[string]interface{}{"code": 409, "error": message}
}

func ErrorBadRequest(message string) map[string]interface{} {
	return map[string]interface{}{"code": 400, "error": message}
}

func ErrorInternalServerError(message string) map[string]interface{} {
	return map[string]interface{}{"code": 500, "error": message}
}
func ResponseOkStruct(field string, data any) map[string]interface{} {
	var myMap []map[string]interface{}
	jsonPresent, _ := json.Marshal(data)

	err := json.Unmarshal(jsonPresent, &myMap)
	if err != nil {
		fmt.Println("unmarshal error", err.Error())
	}
	return map[string]interface{}{"code": 200, field: myMap}
}
func ResponseOkString(field string, value string) map[string]interface{} {
	return map[string]interface{}{"code": 200, field: value}
}
func ResponseCreatedString(field string, value string) map[string]interface{} {
	return map[string]interface{}{"code": 201, field: value}
}

func (f *FiberWebServer) Listen(addr string) {

	err := f.fiber.Listen(addr)
	if err != nil {
		panic(err)
	}
}

func (f *FiberWebServer) Close() {

	err := f.fiber.Shutdown()
	if err != nil {
		panic(err)
	}
}

func (c *FiberRequest) GetBody(data any) map[string]interface{} {
	json.Unmarshal(c.BodyRaw(), data)

	v := reflect.ValueOf(data).Elem() // Предполагаем, что data - это указатель на структуру
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typee := t.Field(i)
		var required = typee.Tag.Get("validate") == "required"
		if required && field.String() == "" {
			return map[string]interface{}{"code": 400, "error": errors.New(fmt.Sprintf("field %s is required", typee.Tag.Get("json"))).Error()}
		}
	}

	return nil

}

func (c *FiberRequest) GetQuery(data any) map[string]interface{} {
	queryData := c.Queries()

	v := reflect.ValueOf(data).Elem() // Предполагаем, что data - это указатель на структуру
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typee := t.Field(i)
		var required = typee.Tag.Get("validate") == "required"
		if required && queryData[field.String()] == "" {
			return map[string]interface{}{"code": 400, "error": errors.New(fmt.Sprintf("field %s is required", typee.Tag.Get("json"))).Error()}
		}
		if field.CanSet() {
			// Устанавливаем новое значение
			field.SetString(queryData[typee.Tag.Get("query")])
		}
	}

	return nil
}

func (f *FiberWebServer) Post(path string, handler Handler) {
	f.fiber.Post(path, func(c *fiber.Ctx) error {
		var request Request = &FiberRequest{c}
		response := handler(request)
		var code = 200
		if response["code"] != nil {
			code = response["code"].(int)
			delete(response, "code")
		}
		return c.Status(code).JSON(response)
	})
}
func (f *FiberWebServer) Get(path string, handler Handler) {
	f.fiber.Get(path, func(c *fiber.Ctx) error {

		var request Request = &FiberRequest{c}
		response := handler(request)
		var code = 200
		if response["code"] != nil {
			code = response["code"].(int)
			delete(response, "code")
		}
		return c.Status(code).JSON(response)
	})
}

func NewFiberWebServer() *FiberWebServer {
	cfg := fiber.Config{}
	cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "*")
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(code).JSON(map[string]interface{}{"detail": err.Error()})
	}
	app := fiber.New(cfg)
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "*")
		c.Set("Access-Control-Allow-Headers", "*")
		return c.Next()
	})
	app.Options("*", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "*")
		c.Set("Access-Control-Allow-Headers", "*")
		return c.Status(200).JSON(map[string]interface{}{"detail": errors.New("rr")})
	})
	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendFile("./toofast.html")
		},
	}))
	return &FiberWebServer{
		fiber: app,
	}
}
