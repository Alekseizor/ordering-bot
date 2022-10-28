package repository

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func GetIDOrder(Db *sqlx.DB, VkID int) (int, error) {
	var ID int
	err := Db.QueryRow("SELECT id from orders WHERE customer_vk_id =$1 ORDER BY id DESC LIMIT 1", VkID).Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with customer_vk_id unknown")
		} else {
			log.Println("Couldn't find the line with the order")
		}
		log.Error(err)
		return -1, err
	}
	return ID, err
}

func GetOrder(Db *sqlx.DB, ID int) (ds.Order, error) {
	var order ds.Order
	err := Db.QueryRow("SELECT * from orders WHERE id =$1", ID).Scan(&order.Id, &order.CustomerVkID, &order.CustomersComment, &order.ExecutorVkID, &order.DisciplineID, &order.DateOrder, &order.DateFinish, &order.Price, &order.PayoutAdmin, &order.PayoutExecutors)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with id unknown")
		} else {
			log.Println("Couldn't find the line with the order")
		}
		log.Error(err)
	}

	return order, err
}
