package repository

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	err := Db.QueryRow("SELECT * from orders WHERE id =$1", ID).Scan(&order.Id, &order.CustomerVkID, &order.CustomersComment, &order.ExecutorVkID, &order.DisciplineID, &order.DateOrder, &order.DateFinish, &order.Price, &order.PayoutAdmin, &order.PayoutExecutors, &order.OrderTask, pq.Array(&order.DocsUrl), pq.Array(&order.ImagesUrl))
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

func WriteUrl(Db *sqlx.DB, VkID int, attachments []object.MessagesMessageAttachment) {
	var docsURL, imagesURL []string
	var num int
	for _, val := range attachments {
		switch val.Type {
		case "doc":
			docsURL = append(docsURL, val.Doc.URL)
		case "photo":
			for i, a := range val.Photo.Sizes {
				if a.Type == "z" {
					num = i
					break
				}
				if a.Type == "x" {
					num = i
				}
			}
			imagesURL = append(imagesURL, val.Photo.Sizes[num].URL)
		default:
			break
		}
	}
	ID, err := GetIDOrder(Db, VkID)
	log.Println(docsURL)
	log.Println(imagesURL)
	res, err := Db.Exec("UPDATE orders SET (docs_url, images_url) = ($1, $2) WHERE id=$3", pq.Array(docsURL), pq.Array(imagesURL), ID)
	log.Println(res)
	if err != nil {
		log.WithError(err).Error("can`t record docs or images url")
	}

}
