package model

type Local struct {
	Local_id int64  `json:"local_id" gorm:"primaryKey;autoIncrement"`
	Nome     string `json:"nome" gorm:"not null"`
	Campus   string `json:"campus" gorm:"not null"`
}

func (Local) TableName() string {
	return "local"
}
