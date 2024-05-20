package v1

import (
	"aabbcc-Server/internal/pkg/db"
	"aabbcc-Server/internal/pkg/server"
	"fmt"
)

type DataGetQuery struct {
	From string `query:"from"`
}

func DataGet(data server.Request) map[string]interface{} {
	query := DataGetQuery{}

	data.GetQuery(&query)

	fmt.Println(query)

	result := db.Database.SelectData(query.From)

	return server.ResponseOkStruct("items", result)
}
