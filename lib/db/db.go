package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"imabad.dev/do/lib/utils"
)

//ConnectDB opens and returns a connection the database using the central configuration
func ConnectDB() (*gorm.DB, error) {
	config := utils.GetConfig()
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database))
	return db, err
}
