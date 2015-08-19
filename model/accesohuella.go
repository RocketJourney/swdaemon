package model

import (
	"time"
)

type Register struct {
	RegId      int       `gorm:"column:isAccesoRegistro;primary_key"`
	UserId     int       `gorm:"column:idPersona"`
	ClubId     int       `gorm:"column:idUn"`
	Date       time.Time `gorm:"column:fecha"`
	Hour       string    `gorm:"column:hora"`
	Status     int       `gorm:"column:status"`
	Line       int       `gorm:"column:carril"`
	WayId      bool      `gorm:"column:idSentido"`
	Message    string    `gorm:"column:mensaje"`
	EmployeeId int       `gorm:"column:idEmpleado"`
	TypeId     int       `gorm:"column:idTipoRegistroAcceso"`
}

func (r Register) TableName() string {
	return "registroacceso"
}
