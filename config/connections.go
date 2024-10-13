package config

import (
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectToDB() *gorm.DB {
	log.Infoln("Connecting to database...")
	var err error
	dbDsn := os.Getenv("DB_DSN")

	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Error connecting to database. Error: ", err)
	}

	return db
}

func ConnectToNatsBroker() *nats.Conn {
	log.Infoln("Connecting to nats...")
	natsDsn := os.Getenv("NATS_DSN")
	nc, _ := nats.Connect(natsDsn)
	return nc
}
