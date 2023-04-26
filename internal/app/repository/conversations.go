package repository

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func CreateConversations(Db *sqlx.DB, orderID int) (err error) {
	_, err = Db.Exec("INSERT INTO conversations (order_id) VALUES ($1)", orderID)
	if err != nil {
		log.WithError(err).Error("can`t record docs")
	}
	return err
}

func AddChatID(Db *sqlx.DB, ChatID int, isExec bool) (err error) {
	if isExec {
		_, err = Db.Exec("UPDATE conversations SET executors_conversation_id = $1", ChatID)
		if err != nil {
			log.WithError(err).Error("can`t record docs")
		}
	} else {
		_, err = Db.Exec("UPDATE conversations SET customers_conversation_id = $1", ChatID)
		if err != nil {
			log.WithError(err).Error("can`t record docs")
		}
	}

	return err
}

func GetConversationID(Db *sqlx.DB, orderID int, isExec bool) (int, error) {
	conversation := ds.Conversation{}
	err := Db.QueryRow("SELECT customers_conversation_id, executors_conversation_id FROM conversations WHERE order_id = $1", orderID).Scan(&conversation.CustomersConversationID, &conversation.ExecutorsConversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No executor with vk_id founded")
		} else {
			log.Println(err)
			log.Println("Query error")
		}
	}
	if isExec {
		return conversation.CustomersConversationID, err
	} else {
		return conversation.ExecutorsConversationID, err
	}

}
