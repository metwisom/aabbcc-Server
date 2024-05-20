package model

type Data struct {
	Id     int    `json:"id" db:"id"`
	Aspect string `json:"aspect" db:"aspect" validate:"required"`
	Value  string `json:"value" db:"value" validate:"required"`
	Time   string `json:"time" db:"time_from" validate:"required"`
}

type Auth struct {
	Id       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}

type Register struct {
	Login    string `json:"login" db:"login" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
}
