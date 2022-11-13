package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func ChangeRequisites(Db *sqlx.DB, message string) {
	_, err := Db.Exec("UPDATE requisites SET requisites = $1 WHERE requisites_id = 1", message)
	if err != nil {
		log.WithError(err).Error("can`t set requisites in requisites table")
	}
}

func GetRequisites(Db *sqlx.DB) string {
	var message string
	err := Db.QueryRow("SELECT requisites from requisites WHERE requisites_id =1").Scan(&message)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with requisites_id unknown")
		} else {
			log.Println("Couldn't find the line with the requisites_id")
		}
		log.Error(err)
	}
	return message
}
