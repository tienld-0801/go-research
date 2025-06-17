package database_configs

import (
	"fmt"
	order_model "go-backend/app/models/orders"
	products_model "go-backend/app/models/products"
	receipt_model "go-backend/app/models/receipts"
	users_model "go-backend/app/models/users"
	"go-backend/internal/constants"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	&order_model.Order{},
	&receipt_model.Receipt{},
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

	const defaultName = "Admin"
	const defaultEmail = "admin@gmail.com"
	const defaultPassword = "admin123"

	var count int64
	if err := db.Model(&users_model.User{}).Where("email = ?", defaultEmail).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin := users_model.User{
			Name:     defaultName,
			Email:    defaultEmail,
			Password: string(hashedPassword),
			Role:     constants.RoleAdmin,
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		fmt.Println("✅ Admin account created: username=admin, password=admin123")
	} else {
		fmt.Println("ℹ️ Admin account already exists, skip seeding.")
	}

	return nil
}

func AddDatabaseContext(c fiber.Ctx) error {
	c.Locals(constants.Db, db)
	return c.Next()
}
