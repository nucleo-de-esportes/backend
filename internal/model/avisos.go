package model


import "time"


type Avisos struct{

	Id      uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo  string           `gorm:"not null" json:"titulo"`
	Mensagem string			 `gorm:"not null" json:"mensagem"`
	Status string			 `gorm:"not null" json:"status"`
	Data_envio time.Time      `gorm:"autoCreateTime" json:"data_envio"`
	Destinatario Destinatario `gorm:"not null; embedded" json:"destinatario"`

}

type Destinatario struct{
	Tipo_destinatario string `json:"tipo_destinatario"`
	id_destinatario uint	 `json:"id_destinatario"`
}


func(Avisos) TableName() string{
	return "avisos"
}