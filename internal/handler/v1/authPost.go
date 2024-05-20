package v1

import (
	"aabbcc-Server/internal/model"
	"aabbcc-Server/internal/pkg/db"
	"aabbcc-Server/internal/pkg/jwt"
	"aabbcc-Server/internal/pkg/server"
)

func AuthPost(body server.Request) map[string]interface{} {

	postData := model.Auth{}
	if err := body.GetBody(&postData); err != nil {
		return err
	}

	data := db.Database.GetUserByLogin(postData.Login)
	if data == nil || data.Password != postData.Password {
		return server.ErrorNotFound("user_not_found")
	}

	token, err := jwt.CreateToken(data.Id)
	if err != nil {
		return server.ErrorInternalServerError(err.Error())
	}

	return server.ResponseOkString("token", token)
}
