package database_configs

import (
	"fmt"
	products_model "go-backend/app/models/products"
	users_model "go-backend/app/models/users"
	"go-backend/internal/constants"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host        string
	Port        string
	User        string
	DBName      string
	Schema      string
	SSLMode     string
	Password    string
	AutoMigrate bool
	DB          *pgxpool.Pool
}

var db *gorm.DB

var tables = []interface{}{
	&users_model.User{},
	&products_model.Product{},
}

func ConnectDatabase() (*gorm.DB, error) {

	config := &Config{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		User:        os.Getenv("DB_USER"),
		DBName:      os.Getenv("DB_NAME"),
		Password:    os.Getenv("DB_PASSWORD"),
		AutoMigrate: true,
		SSLMode:     "disable",
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if config.AutoMigrate {
		err := AutoMigrateModels(db)
		if err != nil {
			log.Fatalf("Error migrating database: %v", err)
			return nil, err
		}
	}

	log.Println("Database connected and tables migrated successfully!!!")
	return db, nil
}

func AutoMigrateModels(db *gorm.DB) error {
	fmt.Println("Auto migrating models...")
	for _, model := range tables {
		err := db.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("failed to auto migrate model %T: %v", model, err)
		}
	}
	return nil
}

func AddDatabaseContext(c fiber.Ctx) error {
	c.Locals(constants.Db, db)
	return c.Next()
}
