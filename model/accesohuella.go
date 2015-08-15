package model

import (
	"time"
)

type Register struct {
	RegId  int       `gorm:"column:isAccesoRegistro;primary_key"`
	UserId int       `gorm:"column:idPersona"`
	ClubId int       `gorm:"column:idUn"`
	Date   time.Time `gorm:"column:fecha"`
}

func (r Register) TableName() string {
	return "registroacceso"
}
