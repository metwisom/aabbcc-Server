package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type FiberWebServer struct {
	fiber *fiber.App
}

type FiberRequest struct {
	*fiber.Ctx
}

func (f *FiberWebServer) Listen(addr string) {

	err := f.fiber.Listen(addr)
	if err != nil {
		return
	}
}

func (c *FiberRequest) GetBody(data any) error {
	json.Unmarshal(c.BodyRaw(), data)

	v := reflect.ValueOf(data).Elem() // Предполагаем, что data - это указатель на структуру
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typee := t.Field(i)
		var required = typee.Tag.Get("validate") == "required"
		if required && field.String() == "" {
			return c.Status(400).JSON(map[string]interface{}{"code": 400, "error": errors.New(fmt.Sprintf("field %s is required", typee.Tag.Get("json"))).Error()})
		}
	}

	return nil

}

func (f *FiberWebServer) Post(path string, handler Handler) {
	f.fiber.Post(path, func(c *fiber.Ctx) error {
		request := &FiberRequest{c}
		response := handler(request)
		var code = 200
		if response["code"] != nil {
			code = response["code"].(int)
		}
		return c.Status(code).JSON(response)
	})
}
func (f *FiberWebServer) Get(path string, handler Handler) {
	f.fiber.Get(path, func(c *fiber.Ctx) error {

		response := handler(nil)
		return c.JSON(response)
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
	return &FiberWebServer{
		fiber: app,
	}
}
