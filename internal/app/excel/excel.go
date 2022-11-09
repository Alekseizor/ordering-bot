package excel

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"time"
)

const (
	layout = "02.01.2006"
)

func CreateRespTable(db *sqlx.DB, firstDateStr string, secondDateStr string, close string, ID ...int) (*excelize.File, error) {
	var payoutAdmin bool
	var rows *sql.Rows
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
	f.SetCellValue("Sheet1", "K1", "Статус заказа")
	f.SetCellValue("Sheet1", "L1", "Статус оплаты")
	firstDate, err := time.Parse(layout, firstDateStr)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		return nil, err
	}
	secondDate, err := time.Parse(layout, secondDateStr)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		return nil, err
	}
	defer rows.Close()
	if ID == nil {
		if close == "Все заказы" {
			rows, err = db.Query("SELECT * FROM orders WHERE date_order > ?1 OR date_order < ?2", firstDate, secondDate)
			if err != nil {
				return nil, err
			}
		} else {
			if close == "Оплата" {
				payoutAdmin = true

			} else if close == "Возврат" {
				payoutAdmin = false
			}
			rows, err = db.Query("SELECT * FROM orders WHERE payout_admin=?1 AND (date_order > ?2 OR date_order < ?3)", payoutAdmin, firstDate, secondDate)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
		}
	} else {
		IDFirst := ID[0]
		if close == "Все заказы" {
			rows, err = db.Query("SELECT * FROM orders WHERE executor_vk_id=?1 AND (date_order > ?2 OR date_order < ?3)", IDFirst, firstDate, secondDate)
			if err != nil {
				return nil, err
			}
		} else {
			if close == "Оплата" {
				payoutAdmin = true

			} else if close == "Возврат" {
				payoutAdmin = false
			}
			rows, err = db.Query("SELECT * FROM orders WHERE executor_vk_id=?1 AND payout_admin=?2 AND (date_order > ?3 OR date_order < ?4)", IDFirst, payoutAdmin, firstDate, secondDate)
			if err != nil {
				return nil, err
			}
		}
	}
	//var orders []ds.Order
	for rows.Next() {
		var order ds.Order
		if err := rows.Scan(&order); err != nil {
			return nil, err
		}
		f.SetCellValue("Sheet1", "A1", "Номер заказа")
		f.SetCellValue("Sheet1", "B1", "Заказчик")
		f.SetCellValue("Sheet1", "C1", "Заказчик")
		f.SetCellValue("Sheet1", "D1", "Дата оформления заказа")
		f.SetCellValue("Sheet1", "E1", "Дисциплина")
		f.SetCellValue("Sheet1", "F1", "Дата завершения заказа")
		f.SetCellValue("Sheet1", "G1", "Полная стоимость")
		f.SetCellValue("Sheet1", "H1", "Процент исполнителя")
		f.SetCellValue("Sheet1", "I1", "Прибыль исполнителя")
		f.SetCellValue("Sheet1", "J1", "Прибыль сервиса")
		f.SetCellValue("Sheet1", "K1", "Статус заказа")
		f.SetCellValue("Sheet1", "L1", "Статус оплаты")
		//orders = append(orders, order)
		f.SetCellValue("Sheet1", "A2", order.Id)
		f.SetCellValue("Sheet1", "B1", *order.ExecutorVkID)
		f.SetCellValue("Sheet1", "C1", order.CustomerVkID)
		f.SetCellValue("Sheet1", "D1", order.DateOrder)
		f.SetCellValue("Sheet1", "E1", order.DisciplineID)
		f.SetCellValue("Sheet1", "F1", order.DateFinish)
		f.SetCellValue("Sheet1", "G1", *order.Price)
		f.SetCellValue("Sheet1", "H1", order.PercentExecutor)
		f.SetCellValue("Sheet1", "I1", (*order.Price)*uint64(order.PercentExecutor)/100)
		f.SetCellValue("Sheet1", "J1", (*order.Price)*(100-uint64(order.PercentExecutor))/100)
		f.SetCellValue("Sheet1", "K1", order.)
		f.SetCellValue("Sheet1", "L1", "Статус оплаты")
		//if err := rows.Scan(&order.ID, &alb.Title, &alb.Artist,
		//	&alb.Price, &alb.Quantity); err != nil {
		//	return albums, err
		//}
		//albums = append(albums, album)
	}
	return f, nil
}
