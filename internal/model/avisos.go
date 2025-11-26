package model

import "time"

type Aviso struct {
	Id            uint           `gorm:"primaryKey;autoIncrement" json:"id_aviso"`
	Titulo        string         `gorm:"not null" json:"titulo"`
	Mensagem      string         `gorm:"not null" json:"mensagem"`
	Status        string         `gorm:"not null" json:"status"`
	DataEnvio     time.Time      `gorm:"autoCreateTime" json:"data_envio"`
	Destinatarios []Destinatario `json:"destinatarios" gorm:"foreignKey:AvisoID;constraint:OnDelete:CASCADE"`
}

func (Aviso) TableName() string {
	return "avisos"
}
