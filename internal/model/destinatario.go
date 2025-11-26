package model

type Destinatario struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	AvisoID   uint   `json:"-"`
	Tipo      string `json:"tipo" binding:"required"`
	DestinoID uint   `json:"id"`
}

func (Destinatario) TableName() string {
	return "destinatarios"
}
