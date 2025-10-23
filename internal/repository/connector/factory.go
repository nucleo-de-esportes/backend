package connector

import (
	"fmt"
	"github.com/nucleo-de-esportes/backend/internal/config"
	"gorm.io/gorm"
)

// DatabaseConnector define o contrato para os conectores de banco.
type DatabaseConnector interface {
	Connect(cfg config.DatabaseConfig) (*gorm.DB, error)
}

// New retorna a implementação do conector baseada no driver.
func New(driver string) (DatabaseConnector, error) {
	switch driver {
	case "postgres":
		return &PostgresConnector{}, nil
	case "oracle":
		return &OracleConnector{}, nil
	default:
		return nil, fmt.Errorf("driver de banco não suportado: %s", driver)
	}
}
