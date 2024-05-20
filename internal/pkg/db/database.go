package db

import (
	"aabbcc-Server/internal/model"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"log"
)

type database struct {
	connection *sqlx.DB
}

var Database = &database{}

func (d *database) Connect() {
	//connection, err := sqlx.Connect("sqlite3", "test.db")
	connection, err := sqlx.Connect("clickhouse", "tcp://127.0.0.1:9000?username=&debug=true")
	if err != nil {
		panic("Ошибка подключения к базе данных: " + err.Error())
	}
	d.connection = connection
}

func (d *database) Close() {
	err := d.connection.Close()
	if err != nil {
		panic("Ошибка закрытия базы данных: " + err.Error())
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
		"password" Text NOT NULL );
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
	query = `
	CREATE TABLE IF NOT EXISTS data (
		id UInt32,
		time_from DateTime DEFAULT now(),
		time_to DateTime DEFAULT now(),
		aspect String,
		value String
	) ENGINE = MergeTree()
	ORDER BY id`

	if _, err := d.connection.Exec(query); err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	query = `CREATE TABLE IF NOT EXISTS user (
		id UInt32,
		login String,
		password String
	) ENGINE = MergeTree()
	ORDER BY id`

	if _, err := d.connection.Exec(query); err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}

	query = `CREATE TABLE IF NOT EXISTS token (
		id UInt32,
		user_fk UInt32,
		token String,
		expire DateTime
	) ENGINE = MergeTree()
	ORDER BY id`
	if _, err := d.connection.Exec(query); err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
}

func (d *database) LastData(data *model.Data) int {
	query := "SELECT id FROM data WHERE aspect = ? ORDER BY time_from DESC LIMIT 1"
	rows, err := d.connection.Query(query, data.Aspect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		var oldData model.Data
		if err := rows.Scan(&oldData.Id); err != nil {
			log.Fatal(err)
			return oldData.Id
		}
	}

	//rows, err := d.connection.Query(query, data.Aspect)
	//if err != nil {
	//	log.Fatal("Ошибка получения последнего id:", err)
	//}
	//
	//if rows.Next() {
	//
	//}
	//rows.Close()
	return 0
}

func (d *database) InsertData(data *model.Data) {
	id := d.LastData(data)

	tx, err := d.connection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if id == 0 {
		//query := "INSERT INTO data (aspect,value) VALUES (?,?)"
		stmt, err := tx.Prepare("INSERT INTO data (aspect,value) VALUES (?,?)")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := stmt.Exec(data.Aspect, data.Value); err != nil {
			log.Fatal(err)
		}

		//if _, err := d.connection.Exec(query, data.Aspect, data.Value); err != nil {
		//	log.Println("Ошибка вставки данных:", err)
		//}
	} else {
		//query := "UPDATE data SET time_to = CURRENT_TIMESTAMP WHERE id = ?;"
		//if _, err := d.connection.Exec(query, id); err != nil {
		//	log.Println("Ошибка обновления данных:", err)
		//}
		//d.updateCache = make(map[int]*model.Data)
		stmt, err := tx.Prepare("UPDATE data SET time_to = CURRENT_TIMESTAMP WHERE id = ?;")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := stmt.Exec(id); err != nil {
			log.Fatal(err)
		}

	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}
func (d *database) SelectData(max string) []model.Data {
	query := ""
	if max == "1" {
		query = "SELECT * FROM (SELECT id, aspect,value,time_from FROM data ORDER BY time_from desc LIMIT 80) t ORDER BY time_from "
	} else {
		query = "SELECT * FROM (SELECT id, aspect,value,time_from FROM data ORDER BY time_from desc LIMIT 80) t ORDER BY time_from "
	}

	rows, err := d.connection.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	datas := make([]model.Data, 0)

	for rows.Next() {
		var data model.Data
		if err := rows.Scan(&data.Id, &data.Aspect, &data.Value, &data.Time); err != nil {
			log.Fatal(err)
		}
		datas = append(datas, data)
	}
	return datas
	//
	//rows, _ := d.connection.Query(query, max)
	//
	//datas := make([]model.Data, 0)
	//
	//for rows.Next() {
	//	var data model.Data
	//	_ = rows.Scan(&data.Id, &data.Aspect, &data.Value, &data.Time)
	//	datas = append(datas, data)
	//}
	//rows.Close()
	//return datas
}

func (d *database) GetUserByLogin(login string) *model.Auth {
	query := "SELECT id,login,password FROM user WHERE login = ?"
	rows, _ := d.connection.Query(query, login)
	for rows.Next() {
		var data model.Auth
		_ = rows.Scan(&data.Id, &data.Login, &data.Password)
		rows.Close()
		return &data
	}
	rows.Close()
	return nil
}

func (d *database) CreateUser(data *model.Register) {
	fmt.Println(data.Login, data.Password)

	query := "INSERT INTO user (login,password) VALUES (?,?)"
	if _, err := d.connection.Exec(query, data.Login, data.Password); err != nil {
		log.Println("Ошибка вставки данных:", err)
	}

}
