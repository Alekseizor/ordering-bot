package excel

import (
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

func CreateRespTable(db *sqlx.DB, firstDate string, secondDate string) *excelize.File {
	f := excelize.NewFile() //создали новый лист
	f.SetCellValue("Sheet1", "A1", "Номер заказа")
	f.SetCellValue("Sheet1", "B1", "Заказчик")
	f.SetCellValue("Sheet1", "C1", "Дата оформления заказа")
	f.SetCellValue("Sheet1", "D1", "Дисциплина")
	f.SetCellValue("Sheet1", "E1", "Дата завершения заказа")
	f.SetCellValue("Sheet1", "F1", "Полная стоимость")
	f.SetCellValue("Sheet1", "G1", "Процент исполнителя")
	f.SetCellValue("Sheet1", "H1", "Прибыль исполнителя")
	f.SetCellValue("Sheet1", "I1", "Прибыль сервиса")
	f.SetCellValue("Sheet1", "J1", "Статус заказа")
	f.SetCellValue("Sheet1", "K1", "Статус оплаты")
	return f
}

func CreateRespTableOneExec(db *sqlx.DB, firstDate string, secondDate string, id int) *excelize.File {
	f := excelize.NewFile() //создали новый лист
	f.SetCellValue("Sheet1", "A1", "Номер заказа")
	f.SetCellValue("Sheet1", "B1", "Заказчик")
	f.SetCellValue("Sheet1", "C1", "Дата оформления заказа")
	f.SetCellValue("Sheet1", "D1", "Дисциплина")
	f.SetCellValue("Sheet1", "E1", "Дата завершения заказа")
	f.SetCellValue("Sheet1", "F1", "Полная стоимость")
	f.SetCellValue("Sheet1", "G1", "Процент исполнителя")
	f.SetCellValue("Sheet1", "H1", "Прибыль исполнителя")
	f.SetCellValue("Sheet1", "I1", "Прибыль сервиса")
	f.SetCellValue("Sheet1", "J1", "Статус заказа")
	f.SetCellValue("Sheet1", "K1", "Статус оплаты")
	return f
}
