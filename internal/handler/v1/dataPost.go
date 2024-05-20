package v1

import (
	"aabbcc-Server/internal/model"
	"aabbcc-Server/internal/pkg/db"
	"aabbcc-Server/internal/pkg/server"
)

type IncomeData struct {
	Items []model.Data `json:"items"`
}

func DataPost(body server.Request) map[string]interface{} {

	data := IncomeData{}
	err := body.GetBody(&data)
	if err != nil {
		return err
	}

	//fmt.Println(data)

	for _, v := range data.Items {
		v.Id = 1

		db.Database.InsertData(&v)
	}

	return server.ResponseOkString("message", "Данные успешно добавлены в базу данных")
}
