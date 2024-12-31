package infrastructure

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"test1/model"
)

func openConnection() (*gorm.DB, error) {
	connectSQL := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" dbname=" + dbName +
		" password=" + dbPassword +
		" sslmode=disable"
	log.Println(connectSQL)

	db, err := gorm.Open(postgres.Open(connectSQL), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		ErrLog.Printf("Not connect to database: %+v\n", err)
		return nil, err
	}
	return db, nil
}
func CloseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		ErrLog.Printf("Not close connection: %+v\n", err)
	}
	dbSQL.Close()
}

// InitDatabase open connection and migrate database
func InitDatabase(allowMigrate bool) error {
	var err error
	db, err = openConnection()
	if err != nil {
		return err
	}

	if allowMigrate {
		log.Println("Migrating database...")
		db.DisableForeignKeyConstraintWhenMigrating = false
		db.DisableForeignKeyConstraintWhenMigrating = false
		if err := db.AutoMigrate(
			&model.Profile{},
			&model.User{},
			&model.UserRole{},
		); err != nil {
			log.Println("Failed to migrate database: ", err)
			return err
		}
		log.Println("Done migrating database")
	}
	return nil
}
