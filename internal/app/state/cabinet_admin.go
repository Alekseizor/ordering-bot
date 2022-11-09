package state

import (
	"database/sql"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

//////////////////////////////////////////////////////////
type CabinetAdmin struct {
}

func (state CabinetAdmin) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка таблицы заказов" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else if messageText == "Назначить исполнителя" {
		AddExecutor{}.PreviewProcess(ctc)
		return &AddExecutor{}
	} else if messageText == "Управлять исполнителями" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else if messageText == "Изменить реквизиты" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else if messageText == "Рассылка" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else if messageText == "Назад в главное меню" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		state.PreviewProcess(ctc)
		return &CabinetAdmin{}
	}
}

func (state CabinetAdmin) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выбери нужный пункт")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Выгрузка таблицы заказов", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назначить исполнителя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Управлять исполнителями", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить реквизиты", "", "secondary")
	k.AddRow()
	k.AddTextButton("Рассылка", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад в главное меню", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state CabinetAdmin) Name() string {
	return "CabinetAdmin"
}

//////////////////////////////////////////////////////////
type UnloadTable struct {
}

func (state UnloadTable) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Незакрытые заказы" {
		UncloseOrder{}.PreviewProcess(ctc)
		return &UncloseOrder{}
	} else if messageText == "Закрытые заказы" {
		СloseOrder{}.PreviewProcess(ctc)
		return &СloseOrder{}
	} else if messageText == "Общая таблица" {
		UncloseOrder{}.PreviewProcess(ctc)
		return &StartState{}
	} else if messageText == "Назад" {
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else {
		state.PreviewProcess(ctc)
		return &UnloadTable{}
	}
}

func (state UnloadTable) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите таблицу")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Незакрытые заказы", "", "secondary")
	k.AddRow()
	k.AddTextButton("Закрытые заказы", "", "secondary")
	k.AddRow()
	k.AddTextButton("Общая таблица", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state UnloadTable) Name() string {
	return "UnloadTable"
}

//////////////////////////////////////////////////////////
type UncloseOrder struct {
}

func (state UncloseOrder) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка" {
		UnloadingUnclose{}.PreviewProcess(ctc)
		return &UnloadingUnclose{}
	} else if messageText == "Очистка" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "Назад" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else {
		state.PreviewProcess(ctc)
		return &UncloseOrder{}
	}
}

func (state UncloseOrder) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите следующее действие")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Выгрузка", "", "secondary")
	k.AddRow()
	k.AddTextButton("Очистка", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state UncloseOrder) Name() string {
	return "UncloseOrder"
}

//////////////////////////////////////////////////////////
type СloseOrder struct {
}

func (state СloseOrder) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка" {
		ExecHistoryOrders{}.PreviewProcess(ctc)
		return &ExecHistoryOrders{}
	} else if messageText == "Очистка" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "Назад" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else {
		state.PreviewProcess(ctc)
		return &СloseOrder{}
	}
}

func (state СloseOrder) PreviewProcess(ctc ChatContext) {
	UncloseOrder{}.PreviewProcess(ctc)
}

func (state СloseOrder) Name() string {
	return "СloseOrder"
}

//////////////////////////////////////////////////////////
type GeneralTable struct {
}

func (state GeneralTable) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка" {
		ExecHistoryOrders{}.PreviewProcess(ctc)
		return &ExecHistoryOrders{}
	} else if messageText == "Очистка" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "Назад" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else {
		state.PreviewProcess(ctc)
		return &GeneralTable{}
	}
}

func (state GeneralTable) PreviewProcess(ctc ChatContext) {
	UncloseOrder{}.PreviewProcess(ctc)
}

func (state GeneralTable) Name() string {
	return "GeneralTable"
}

//////////////////////////////////////////////////////////
type UnloadingUnclose struct {
}

func (state UnloadingUnclose) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "За весь период" {
		AllUnloadingUnclose{}.PreviewProcess(ctc)
		return &AllUnloadingUnclose{}
	} else if messageText == "По дате(с ДД.ММ.ГГ до ДД.ММ.ГГ)" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "За месяц" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "За 7 дней" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "За день" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
	} else if messageText == "Назад" {
		UncloseOrder{}.PreviewProcess(ctc)
		return &UncloseOrder{}
	} else {
		state.PreviewProcess(ctc)
		return &UnloadingUnclose{}
	}
}

func (state UnloadingUnclose) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите период выгрузки")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("За весь период", "", "secondary")
	k.AddRow()
	k.AddTextButton("По дате(с ДД.ММ.ГГ до ДД.ММ.ГГ)", "", "secondary")
	k.AddRow()
	k.AddTextButton("За месяц", "", "secondary")
	k.AddRow()
	k.AddTextButton("За 7 дней", "", "secondary")
	k.AddRow()
	k.AddTextButton("За день", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state UnloadingUnclose) Name() string {
	return "UnloadingUnclose"
}

//////////////////////////////////////////////////////////
type AllUnloadingUnclose struct {
}

func (state AllUnloadingUnclose) Process(ctc ChatContext, msg object.MessagesMessage) State {
	//messageText := msg.Text
	//if messageText == "Да" {
	//	table := excel.CreateRespTableOneExec(ctc.Db, "12.11.2006", "12.11.2099", ctc.User.VkID)
	//
	//} else if messageText == "Нет" {
	//	table := excel.CreateRespTable(ctc.Db, "12.11.2006", "12.11.2099")
	//	UnloadingUnclose{}.PreviewProcess(ctc)
	//	return &UnloadingUnclose{}
	//} else if messageText == "Назад" {
	//	UnloadingUnclose{}.PreviewProcess(ctc)
	//	return &UnloadingUnclose{}
	//} else {
	//	state.PreviewProcess(ctc)
	//	return &AllUnloadingUnclose{}
	//}
	state.PreviewProcess(ctc)
	return &AllUnloadingUnclose{}
}

func (state AllUnloadingUnclose) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("По конкретному исполнителю?")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state AllUnloadingUnclose) Name() string {
	return "AllUnloadingUnclose"
}

//////////////////////////////////////////////////////////
type AddExecutor struct {
}

func (state AddExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else if messageText == "Добавить исполнителя" {
		AddExecID{}.PreviewProcess(ctc)
		return &AddExecID{}
	} else {
		state.PreviewProcess(ctc)
		return &AddExecutor{}
	}
}

func (state AddExecutor) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выбери нужный пункт")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Добавить исполнителя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state AddExecutor) Name() string {
	return "AddExecutor"
}

//////////////////////////////////////////////////////////
type AddExecID struct {
}

func (state AddExecID) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		AddExecutor{}.PreviewProcess(ctc)
		return &AddExecutor{}
	} else {
		messageInt, err := strconv.Atoi(messageText)
		if err != nil {
			state.PreviewProcess(ctc)
			return &AddExecID{}
		} else {
			_, err := ctc.Db.ExecContext(*ctc.Ctx, "INSERT INTO executors(vk_id) VALUES ($1) ON CONFLICT (vk_id) DO UPDATE SET vk_id = $1", messageInt)
			if err != nil {
				log.WithError(err).Error("cant add Executor on state ChoiceDiscipline")
				state.PreviewProcess(ctc)
				return &AddExecID{}
			}
			AddExecDisciplines{}.PreviewProcess(ctc)
			return &AddExecDisciplines{}
		}
	}
}

func (state AddExecID) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID ВК нового исполнителя")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state AddExecID) Name() string {
	return "AddExecID"
}

//////////////////////////////////////////////////////////
type AddExecDisciplines struct {
}

func (state AddExecDisciplines) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		AddExecID{}.PreviewProcess(ctc)
		return &AddExecID{}
	} else {
		disciplines := strings.Split(messageText, " ")
		for _, val := range disciplines {
			num, err := strconv.Atoi(val)
			if err != nil || (num < 1) || (num > 52) {
				state.PreviewProcess(ctc)
				return &AddExecDisciplines{}
			}
		}
		var VkID int
		err := ctc.Db.QueryRow("SELECT vk_id from executors ORDER BY id DESC LIMIT 1").Scan(&VkID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Row with customer_vk_id unknown")
			} else {
				log.Println("Couldn't find the line with the order")
			}
			state.PreviewProcess(ctc)
			return &AddExecDisciplines{}
		}
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE executors SET disciplines_id = $1 WHERE vk_id=$2", pq.Array(disciplines), VkID)
		if err != nil {
			log.WithError(err).Error("cant add Executor on state ChoiceDiscipline")
			state.PreviewProcess(ctc)
			return &AddExecID{}
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(ctc.User.VkID)
		vkID := strconv.Itoa(VkID)
		b.Message("Исполнитель создан\nID исполнителя: " + vkID + "\nНомера выбранных предметов: " + messageText)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	}
}

func (state AddExecDisciplines) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	b.Message("1. MATLAB\n2. MS Office(Word, excel, Access)\n3. Mathcad\n4. Аналитическая геометрия\n5. Английский язык\n6. Детали машин\n7. Дискретная математика\n8. Инженерная и компьютерная графика\n9. Интегралы и дифференциальные уравнения\n10. Информатика\n11. История\n12. Кратные интегралы и ряды\n13. Культурология\n14. Линейная алгебра\n15. Математика\n16. Математический анализ\n17. Материаловедение\n18. Менеджмент\n19. Метрология\n20. Механика жидкости и газа\n21. Начертательная геометрия\n22. Организация производства\n23. Основы конструирования приборов\n24. Основы теории цепей\n25. Основы технологии приборостроения\n26. Политология\n27. Правоведение\n28. Практика\n29. Прикладная статистика\n30. Психология\n31. Системный анализ и принятие решений\n32. Сопротивление материалов\n33. Социология\n34. Теоретическая механика\n35. Теоретические основы электротехники\n36. Теория вероятностей\n37. Теория механизмов и машин\n38. Теория поля\n39. Теория функции комплексных переменных и операционное исчисление\n40. Теория функции нескольких переменных\n41. Термодинамика\n42. Технология конструкционных материалов\n43. Уравнения математической физики\n44. Физика\n45. Физкультура\n46. Философия\n47. Финансирование инновационной деятельности\n48. Химия\n49. Цифровые устройства и микропроцессоры\n50. Экономика\n51. Электроника\n52. Электротехника")
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	b.Message("Введите номера предметов через пробел")
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state AddExecDisciplines) Name() string {
	return "AddExecDisciplines"
}
