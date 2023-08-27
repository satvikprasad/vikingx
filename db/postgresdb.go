package db

import (
	"fmt"
	"os"

	"github.com/satvikprasad/vikingx/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	Db *gorm.DB
}

func NewDB() (*PostgresDB, error) {
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

	dbInstance := &PostgresDB{
		Db: db,
	}

	return dbInstance, nil
}

func (d *PostgresDB) CreateTrade(t *models.Trade) {
	d.Db.Create(&t)
}

func (d *PostgresDB) Trades() []*models.Trade {
	trades := []*models.Trade{}
	d.Db.Find(&trades)
	return trades
}
