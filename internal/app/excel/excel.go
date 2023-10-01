package excel

import (
	"database/sql"
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func CreateRespTable(db *sqlx.DB, firstDateStr string, secondDateStr string, close string, ID ...int) (*excelize.File, error) {
	var rows *sql.Rows
	var err error
	var requestStr, orderStatus string
	var exec ds.Executor
	var style int

	f := excelize.NewFile() //создали новый лист
	f.SetCellValue("Sheet1", "A1", "Номер заказа")
	f.SetCellValue("Sheet1", "B1", "Заказчик")
	f.SetCellValue("Sheet1", "C1", "Исполнитель")
	f.SetCellValue("Sheet1", "D1", "Дата оформления заказа")
	f.SetCellValue("Sheet1", "E1", "Дисциплина")
	f.SetCellValue("Sheet1", "F1", "Дата завершения заказа")
	f.SetCellValue("Sheet1", "G1", "Полная стоимость")
	f.SetCellValue("Sheet1", "H1", "Процент исполнителя")
	f.SetCellValue("Sheet1", "I1", "Прибыль исполнителя")
	f.SetCellValue("Sheet1", "J1", "Прибыль сервиса")
	f.SetCellValue("Sheet1", "K1", "Реквизиты исполнителя")
	f.SetCellValue("Sheet1", "L1", "Статус заказа")
	f.SetCellValue("Sheet1", "M1", "Статус оплаты")
	log.Println(close)
	if ID == nil {
		if close == "Общая таблица" {
			requestStr = fmt.Sprintf("SELECT * FROM orders WHERE date_order > '%s' AND date_order < '%s'", firstDateStr, secondDateStr)
		} else if close == "Закрытые заказы" {
			requestStr = fmt.Sprintf("SELECT * FROM orders WHERE verification_executor=true AND verification_customer=true AND date_order > '%s' AND date_order < '%s'", firstDateStr, secondDateStr)
		} else if close == "Незакрытые заказы" {
			requestStr = fmt.Sprintf("SELECT * FROM orders WHERE (verification_executor IS NULL OR verification_customer IS NULL) AND date_order > '%s' AND date_order < '%s'", firstDateStr, secondDateStr)
		}
	} else {
		IDFirst := ID[0]
		if close == "Общая таблица" {
			requestStr = fmt.Sprintf("SELECT * FROM orders WHERE executor_vk_id=%d AND date_order > '%s' AND date_order < '%s'", IDFirst, firstDateStr, secondDateStr)
		} else if close == "Закрытые заказы" {
			requestStr = fmt.Sprintf("SELECT * FROM orders WHERE executor_vk_id=%d AND verification_executor=true AND verification_customer=true AND date_order > '%s' AND date_order < '%s'", IDFirst, firstDateStr, secondDateStr)
		} else if close == "Незакрытые заказы" {
			requestStr = fmt.Sprintf("SELECT * FROM orders WHERE executor_vk_id=%d AND (verification_executor IS NULL OR verification_customer IS NULL) AND date_order > '%s' AND date_order < '%s'", IDFirst, firstDateStr, secondDateStr)
		}
	}
	rows, err = db.Query(requestStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var order ds.Order
	var price uint64
	var requisite string
	var executor, percentExecutor uint
	numberRow := 2
	for rows.Next() {
		if err := rows.Scan(&order.Id, &order.CustomerVkID, &order.CustomersComment, &order.ExecutorVkID, &order.TypeOrder, &order.DisciplineID, &order.DateOrder, &order.DateFinish, &order.Price, &order.PercentExecutor, &order.VerificationExecutor, &order.VerificationCustomer, &order.OrderTask); err != nil {
			log.Println(err)
			return nil, err
		}
		//orders = append(orders, order)
		if order.Price == nil {
			price = 0
		} else {
			price = *(order.Price)
		}
		if order.ExecutorVkID == nil {
			executor = 0
		} else {
			executor = *(order.ExecutorVkID)
			err := db.QueryRow("SELECT requisite from executors WHERE vk_id =$1", executor).Scan(&exec.Requisite)
			if err != nil {
				log.Error(err)
			}

			if exec.Requisite == nil {
				requisite = "Исполнитель не оставил свои реквизиты"
			} else {
				requisite = *exec.Requisite
			}
		}
		if order.PercentExecutor == nil {
			percentExecutor = 0
		} else {
			percentExecutor = *(order.PercentExecutor)
		}
		if order.VerificationCustomer == nil || order.VerificationExecutor == nil {
			orderStatus = "Не закрыт"
			style, err = f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{"#ae0000"}, Pattern: 1},
			})
			if err != nil {
				fmt.Println(err)
			}
		} else {
			orderStatus = "Закрыт"
			style, err = f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{"#9fff40"}, Pattern: 1},
			})
			if err != nil {
				fmt.Println(err)
			}
		}
		numberRowStr := strconv.Itoa(numberRow)
		f.SetCellValue("Sheet1", "A"+numberRowStr, order.Id)
		f.SetCellValue("Sheet1", "B"+numberRowStr, order.CustomerVkID)
		f.SetCellValue("Sheet1", "C"+numberRowStr, executor)
		f.SetCellValue("Sheet1", "D"+numberRowStr, order.DateOrder)
		f.SetCellValue("Sheet1", "E"+numberRowStr, order.DisciplineID)
		f.SetCellValue("Sheet1", "F"+numberRowStr, order.DateFinish)
		f.SetCellValue("Sheet1", "G"+numberRowStr, price)
		f.SetCellValue("Sheet1", "H"+numberRowStr, percentExecutor)
		f.SetCellValue("Sheet1", "I"+numberRowStr, price*uint64(percentExecutor)/100)
		f.SetCellValue("Sheet1", "J"+numberRowStr, price*(100-uint64(percentExecutor))/100)
		f.SetCellValue("Sheet1", "K"+numberRowStr, requisite)
		f.SetCellValue("Sheet1", "L"+numberRowStr, orderStatus)
		f.SetCellStyle("Sheet1", "L"+numberRowStr, "L"+numberRowStr, style)
		//f.SetCellValue("Sheet1", "K1", order.)
		//f.SetCellValue("Sheet1", "L1", "Статус оплаты")
		if err := f.SaveAs("Book1.xlsx"); err != nil {
			log.Println(err)
			return nil, err
		}
		numberRow++
	}
	return f, err
}

func CreateExecTable(db *sqlx.DB) (*excelize.File, error) {
	var rows *sql.Rows
	var err error
	var requestStr string
	var executor ds.Executor

	f := excelize.NewFile() //создали новый лист
	f.SetCellValue("Sheet1", "A1", "ID исполнителя")
	f.SetCellValue("Sheet1", "B1", "VK_ID исполнителя")
	f.SetCellValue("Sheet1", "C1", "Дисциплины")
	f.SetCellValue("Sheet1", "D1", "Процент исполнителя")
	f.SetCellValue("Sheet1", "E1", "Рейтинг исполнителя")
	f.SetCellValue("Sheet1", "F1", "Количество оценок")
	f.SetCellValue("Sheet1", "G1", "Прибыль")
	f.SetCellValue("Sheet1", "H1", "Количество заказов")
	f.SetCellValue("Sheet1", "I1", "Реквизиты")

	requestStr = fmt.Sprintf("SELECT * FROM executors ")
	rows, err = db.Query(requestStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var requisite string

	numberRow := 2
	for rows.Next() {
		var DisciplinesID []sql.NullInt64
		if err := rows.Scan(&executor.Id, &executor.VkID, pq.Array(&DisciplinesID), &executor.PercentExecutor, &executor.Rating, &executor.AmountRating, &executor.Profit, &executor.AmountOrders, &executor.Requisite); err != nil {
			log.Println(err)
			return nil, err
		}

		var disciplines []int64
		for _, val := range DisciplinesID {
			disciplines = append(disciplines, val.Int64)
		}
		if executor.Requisite == nil {
			requisite = "Исполнитель не оставил свои реквизиты"
		} else {
			requisite = *executor.Requisite
		}

		numberRowStr := strconv.Itoa(numberRow)
		f.SetCellValue("Sheet1", "A"+numberRowStr, executor.Id)
		f.SetCellValue("Sheet1", "B"+numberRowStr, executor.VkID)
		f.SetCellValue("Sheet1", "C"+numberRowStr, disciplines)
		f.SetCellValue("Sheet1", "D"+numberRowStr, executor.PercentExecutor)
		f.SetCellValue("Sheet1", "E"+numberRowStr, executor.Rating)
		f.SetCellValue("Sheet1", "F"+numberRowStr, executor.AmountRating)
		f.SetCellValue("Sheet1", "G"+numberRowStr, executor.Profit)
		f.SetCellValue("Sheet1", "H"+numberRowStr, executor.AmountOrders)
		f.SetCellValue("Sheet1", "I"+numberRowStr, requisite)

		if err := f.SaveAs("Book2.xlsx"); err != nil {
			log.Println(err)
			return nil, err
		}
		numberRow++
	}
	return f, err
}
