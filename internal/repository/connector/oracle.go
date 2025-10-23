package connector

import (
	"fmt"

	oracle "github.com/godoes/gorm-oracle"
	"github.com/nucleo-de-esportes/backend/internal/config"
	"gorm.io/gorm"
)

type OracleConnector struct{}

func (o *OracleConnector) Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s/%s@%s", cfg.User, cfg.Password, cfg.Name)
	return gorm.Open(oracle.Open(dsn), &gorm.Config{})
}
