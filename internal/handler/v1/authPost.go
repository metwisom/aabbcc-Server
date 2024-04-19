package v1

import (
	"aabbcc-Server/internal/pkg/db"
	"aabbcc-Server/internal/pkg/server"
	"encoding/json"
	"fmt"
)

func AuthPost(_ *server.FiberRequest) map[string]interface{} {

	//data, _ := body.GetBody(&model.Auth{})

	result := db.Database.SelectData()

	var myMap []map[string]interface{}
	data, _ := json.Marshal(result)

	err := json.Unmarshal(data, &myMap)
	if err != nil {
		fmt.Println("unmarshal error", err.Error())
	}

	return map[string]interface{}{"items": myMap}
}
