package v1

import (
	"aabbcc-Server/internal/model"
	"aabbcc-Server/internal/pkg/db"
	"aabbcc-Server/internal/pkg/server"
)

func DataPost(body *server.FiberRequest) map[string]interface{} {

	data := model.Data{}
	err := body.GetBody(&data)
	if err != nil {
		return map[string]interface{}{"code": 400, "error": err.Error()}
	}

	db.Database.InsertData(&data)
	response := map[string]interface{}{"message": "Данные успешно добавлены в базу данных"}
	return response
}
