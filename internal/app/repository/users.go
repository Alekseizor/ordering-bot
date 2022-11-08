package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func GetState(Db *sqlx.DB, VkID int) string {
	var state string
	err := Db.QueryRow("SELECT state from users WHERE vk_id =$1", VkID).Scan(&state)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with id unknown")
		} else {
			log.Println("Couldn't find the users state")
		}
		log.Error(err)
	}
	return state
}
