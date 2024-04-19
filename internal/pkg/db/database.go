package db

import (
	"aabbcc-Server/internal/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type database struct {
	connection *sqlx.DB
}

var Database = &database{}

func (d *database) Connect() {
	var err error
	d.connection, err = sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
}

func (d *database) Close() {
	err := d.connection.Close()
	if err != nil {
		return
	}
}

func (d *database) CreateTable() {
	query := `
    BEGIN;
	CREATE TABLE IF NOT EXISTS data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		time_to TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		aspect TEXT NOT NULL,
		value TEXT NOT NULL	
	);
	CREATE TABLE IF NOT EXISTS "user"(
		"id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"login" Text NOT NULL,
		"pawssord" Text NOT NULL );
	CREATE INDEX IF NOT EXISTS "index_login" ON "user"( "login" );
	CREATE TABLE IF NOT EXISTS "token"(
		"id" Integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"user_fk" Integer,
		"token" Text,
		"expire" DateTime );
	CREATE INDEX IF NOT EXISTS "index_token" ON "token"( "token" );
	CREATE INDEX IF NOT EXISTS "index_user_fk" ON "token"( "user_fk" );
	COMMIT;
	`
	if _, err := d.connection.Exec(query); err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
}

func (d *database) LastData(data *model.Data) int {
	query := "SELECT id,value FROM data WHERE aspect =? ORDER BY time_from DESC LIMIT 1"
	rows, _ := d.connection.Query(query, data.Aspect)

	if rows.Next() {
		var oldData model.Data
		_ = rows.Scan(&oldData.Id, &oldData.Value)
		if oldData.Value == data.Value {
			rows.Close()
			return oldData.Id
		}
	}
	rows.Close()
	return 0
}

func (d *database) InsertData(data *model.Data) {
	fmt.Println(data.Aspect, data.Value)
	id := d.LastData(data)
	if id == 0 {
		fmt.Println("NO")
		fmt.Println(data)
		query := "INSERT INTO data (aspect,value) VALUES (?,?)"
		if _, err := d.connection.Exec(query, data.Aspect, data.Value); err != nil {
			log.Println("Ошибка вставки данных:", err)
		}
	} else {
		query := "UPDATE data SET time_to = CURRENT_TIMESTAMP WHERE id =?"
		if _, err := d.connection.Exec(query, id); err != nil {
			log.Println("Ошибка обновления данных:", err)
		}
	}

}
func (d *database) SelectData() []model.Data {
	query := "SELECT * FROM (SELECT aspect,value,time FROM data ORDER BY time desc LIMIT 80) t ORDER BY time "
	rows, _ := d.connection.Query(query)
	//fmt.Println(rows)
	datas := make([]model.Data, 0)
	// iterate over each row
	for rows.Next() {
		var data model.Data
		_ = rows.Scan(&data.Aspect, &data.Value, &data.Time)
		datas = append(datas, data)
	}
	return datas
}
