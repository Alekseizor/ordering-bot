package repository

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func IsExecutor(Db *sqlx.DB, vkID int) (bool, error) {
	var exec ds.Executor
	err := Db.QueryRow("SELECT 1 from executors WHERE vk_id = $1", vkID).Scan(&exec.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id founded")
			return false, err
		} else {
			log.Println(err)
			log.Println("Query error")
			return false, err
		}
	}
	return true, nil
}

func DeleteExecutor(Db *sqlx.DB, vkID int) error {
	log.Println(vkID)
	_, err := Db.Exec("DELETE from executors WHERE vk_id = $1", vkID)
	if err != nil {
		log.Println("Executor with vk_id unknown. Delete failed")
		return err
	}
	return nil
}
