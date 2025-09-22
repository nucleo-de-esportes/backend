package model

import (
	"time"

	"github.com/google/uuid"
)

type Presenca struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AulaID   uint      `gorm:"not null;index" json:"aula_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Presente bool      `gorm:"default:false" json:"presente"`
	CriadoEm time.Time `gorm:"autoCreateTime" json:"criado_em"`

	Aula Aula `gorm:"foreignKey:AulaID;constraint:OnDelete:CASCADE" json:"aula"`
	User User `gorm:"foreignKey:UserID;references:User_id;constraint:OnDelete:CASCADE" json:"user"`
}
