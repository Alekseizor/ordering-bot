package repository

import (
	"database/sql"
	"errors"
	"fmt"
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
	err := Db.QueryRow("SELECT * from orders WHERE id =$1", ID).Scan(&order.Id, &order.CustomerVkID, &order.CustomersComment, &order.ExecutorVkID, &order.TypeOrder, &order.DisciplineID, &order.DateOrder, &order.DateFinish, &order.Price, &order.PercentExecutor, &order.VerificationExecutor, &order.VerificationCustomer, &order.OrderTask)
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
	case "ChoiceTime", "CommentOrder", "ConfirmationOrder": //todo: На стейте TaskOrder при нажатии назад пишет второй вариант вместо первого (хз надо ли фиксить)
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

func ClearTable(db *sqlx.DB, firstDateStr string, secondDateStr string, close string) error {
	var requestStr string
	if close == "Общая таблица" {
		requestStr = fmt.Sprintf("DELETE FROM orders WHERE date_order > '%s' AND date_order < '%s'", firstDateStr, secondDateStr)
	} else if close == "Закрытые заказы" {
		requestStr = fmt.Sprintf("DELETE FROM orders WHERE verification_executor=true AND verification_customer=true AND date_order > '%s' AND date_order < '%s'", firstDateStr, secondDateStr)
	} else if close == "Незакрытые заказы" {
		requestStr = fmt.Sprintf("DELETE FROM orders WHERE (verification_executor IS NULL OR verification_customer IS NULL) AND date_order > '%s' AND date_order < '%s'", firstDateStr, secondDateStr)
	}
	_, err := db.Exec(requestStr)
	if err != nil {
		log.WithError(err).Error("cant delete orders")
		return err
	}
	return nil
}

func AddingExecutor(db *sqlx.DB, executorOrder ds.ExecutorOrder) error {
	order, err := GetOrder(db, executorOrder.OrderID)
	if err != nil {
		log.WithError(err).Error("couldn't request an order")
		return err
	}
	if order.ExecutorVkID != nil {
		return errors.New("the executor has already been selected")
	}
	executor, err := GetExecutorByID(db, executorOrder.ExecutorID)
	if err != nil {
		log.WithError(err).Error("failed to request a executor by ID")
		return err
	}
	_, err = db.Exec("UPDATE orders SET executor_vk_id=$1,price=$2,percent_executor=$3 WHERE id=$4", executor.VkID, executorOrder.Price, executor.PercentExecutor, executorOrder.OrderID)
	if err != nil {
		log.WithError(err).Error("cant delete orders")
		return err
	}
	return nil
}

func FinishOrder(db *sqlx.DB, orderID int, isExec bool) error {
	if isExec {
		_, err := db.Exec("UPDATE orders SET verification_executor=$1 WHERE id=$2", true, orderID)
		if err != nil {
			log.WithError(err).Error("cant delete orders")
			return err
		}
	} else {
		_, err := db.Exec("UPDATE orders SET verification_customer=$1 WHERE id=$2", true, orderID)
		if err != nil {
			log.WithError(err).Error("cant delete orders")
			return err
		}
	}

	return nil
}
