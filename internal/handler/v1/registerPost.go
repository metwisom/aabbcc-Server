package v1

import (
	"aabbcc-Server/internal/model"
	"aabbcc-Server/internal/pkg/db"
	"aabbcc-Server/internal/pkg/server"
)

func RegisterPost(body server.Request) map[string]interface{} {

	postData := model.Register{}
	if err := body.GetBody(&postData); err != nil {
		return err
	}

	data := db.Database.GetUserByLogin(postData.Login)
	if data != nil {
		return server.ErrorConflict("user_already_exist")
	}

	db.Database.CreateUser(&postData)

	return server.ResponseCreatedString("result", "success")
}
