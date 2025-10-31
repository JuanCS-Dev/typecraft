package database

import (
	"fmt"
	"log"

	"github.com/JuanCS-Dev/typecraft/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB é a instância global do banco de dados
var DB *gorm.DB

// Connect estabelece conexão com o banco de dados
func Connect(databaseURL string) error {
	var err error
	
	// Configurar logger do GORM
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco: %w", err)
	}
	
	log.Println("✅ Conexão com banco de dados estabelecida")
	
	return nil
}

// Migrate executa as migrations automáticas
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("banco de dados não inicializado")
	}
	
	log.Println("🔄 Executando migrations...")
	
	err := DB.AutoMigrate(
		&domain.Project{},
		&domain.Job{},
	)
	
	if err != nil {
		return fmt.Errorf("falha nas migrations: %w", err)
	}
	
	log.Println("✅ Migrations concluídas")
	
	return nil
}

// Close fecha a conexão com o banco de dados
func Close() error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Close()
}

// Health verifica se o banco está respondendo
func Health() error {
	if DB == nil {
		return fmt.Errorf("banco de dados não inicializado")
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Ping()
}
