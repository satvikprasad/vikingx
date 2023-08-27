package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/satvikprasad/vikingx/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

func NewDB() (*DbInstance, error) {
	godotenv.Load(".env")

	fmt.Println(os.Getenv("DB_PASSWORD"))

	dsn := fmt.Sprintf(`host=db user=%s password=%s dbname=%s port=5432 
        sslmode=disable`,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Trade{})

	dbInstance := &DbInstance{
		Db: db,
	}

	return dbInstance, nil
}
