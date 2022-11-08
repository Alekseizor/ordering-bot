package repository

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strconv"
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
	err := Db.QueryRow("SELECT * from orders WHERE id =$1", ID).Scan(&order.Id, &order.CustomerVkID, &order.CustomersComment, &order.ExecutorVkID, &order.TypeOrder, &order.DisciplineID, &order.DateOrder, &order.DateFinish, &order.Price, &order.PayoutAdmin, &order.PayoutExecutors, &order.OrderTask)
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

func GetCompleteOrder(Db *sqlx.DB, VkID int) (string, error) {
	var output string
	ID, err := GetIDOrder(Db, VkID)
	if err != nil {
		log.WithError(err).Error("can`t get id with user vk id")
		return output, err
	}
	order, err := GetOrder(Db, ID)
	if err != nil {
		return output, err
	}
	disciplineName, err := GetDisciplineName(Db, order.DisciplineID)
	if err != nil {
		return output, err
	}
	dateFinish := strconv.Itoa(order.DateFinish.Day()) + "." + order.DateFinish.Format("01") + "." + strconv.Itoa(order.DateFinish.Year())
	var orderTask string
	if order.OrderTask != nil {
		orderTask = *order.OrderTask
	}
	state := GetState(Db, VkID)
	switch state {
	case "ChoiceTime", "CommentOrder": //todo: На стейте TaskOrder при нажатии назад пишет второй вариант вместо первого (хз надо ли фиксить)
		output = "Ваш заказ:\nВид работы - " + order.TypeOrder + "\nДисциплина - " + disciplineName + "\nДата выполнения - " + dateFinish + "\nВремя выполнения - " + order.DateFinish.Format("15:04")
		break
	case "TaskOrder", "EditType", "EditDiscipline", "EditDate", "EditTime", "EditTaskOrder", "EditCommentOrder", "OrderChange", "OrderCancel", "OrderCompleted":
		if order.CustomersComment != nil {
			customerComment := *order.CustomersComment
			output = "Проверьте заказ:\nВид работы - " + order.TypeOrder + "\nДисциплина - " + disciplineName + "\nДата выполнения - " + dateFinish + "\nВремя выполнения - " + order.DateFinish.Format("15:04") + "\nИнформация по заказу - " + orderTask + "\nКомментарий к заказу - " + customerComment //вывод заказа пользователя
			break
		} else {
			output = "Проверьте заказ:\nВид работы - " + order.TypeOrder + "\nДисциплина - " + disciplineName + "\nДата выполнения - " + dateFinish + "\nВремя выполнения - " + order.DateFinish.Format("15:04") + "\nИнформация по заказу - " + orderTask //вывод заказа пользователя
			break
		}
	default:
		break
	}
	return output, nil
}
