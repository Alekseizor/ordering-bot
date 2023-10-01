package state

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/excel"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var firstDateStr, secondDateStr, close string

type CabinetAdmin struct {
}

func (state CabinetAdmin) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка таблицы заказов" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else if messageText == "Выгрузка таблицы исполнителей" {
		UnloadTableExec{}.PreviewProcess(ctc)
		return &UnloadTableExec{}
	} else if messageText == "Назначить исполнителя" {
		AddExecutor{}.PreviewProcess(ctc)
		return &AddExecutor{}
	} else if messageText == "Управлять исполнителями" {
		ManageExecutors{}.PreviewProcess(ctc)
		return &ManageExecutors{}
	} else if messageText == "Изменить реквизиты" {
		ChangeRequisites{}.PreviewProcess(ctc)
		return &ChangeRequisites{}
	} else if messageText == "Рассылка" {
		Newsletter{}.PreviewProcess(ctc)
		return &Newsletter{}
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
	k.AddTextButton("Выгрузка таблицы исполнителей", "", "secondary")
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

// ////////////////////////////////////////////////////////
type UnloadTable struct {
}

func (state UnloadTable) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Незакрытые заказы" || messageText == "Закрытые заказы" || messageText == "Общая таблица" {
		close = messageText
		SelectionUnload{}.PreviewProcess(ctc)
		return &SelectionUnload{}
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

// ////////////////////////////////////////////////////////
type UnloadTableExec struct {
}

func (state UnloadTableExec) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Да" {
		table, err := excel.CreateExecTable(ctc.Db)
		if err != nil {
			log.Println(err)
			state.PreviewProcess(ctc)
			return &CabinetAdmin{}
		}
		tableBuffer, err := table.WriteToBuffer()
		if err != nil {
			log.Println(err)
			state.PreviewProcess(ctc)
			return &CabinetAdmin{}
		}
		tableBytes := tableBuffer.Bytes()
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Ваша таблица:")
		b.PeerID(ctc.User.VkID)
		doc, err := ctc.Vk.UploadMessagesDoc(ctc.User.VkID, "doc", "Book2.xlsx", "", bytes.NewReader(tableBytes))
		if err != nil {
			log.Error(err)
			state.PreviewProcess(ctc)
			return &CabinetAdmin{}
		}
		b.Attachment(doc.Type + strconv.Itoa(doc.Doc.OwnerID) + "_" + strconv.Itoa(doc.Doc.ID))
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else if messageText == "Назад" {
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else {
		state.PreviewProcess(ctc)
		return &UnloadTableExec{}
	}
}

func (state UnloadTableExec) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выгрузить таблицу?")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state UnloadTableExec) Name() string {
	return "UnloadTableExec"
}

// ////////////////////////////////////////////////////////
type SelectionUnload struct {
}

func (state SelectionUnload) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка" {
		SelectionDateUnload{}.PreviewProcess(ctc)
		return &SelectionDateUnload{}
	} else if messageText == "Очистка" {
		SelectionDateClear{}.PreviewProcess(ctc)
		return &SelectionDateClear{}
	} else if messageText == "Назад" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else {
		state.PreviewProcess(ctc)
		return &SelectionUnload{}
	}
}

func (state SelectionUnload) PreviewProcess(ctc ChatContext) {
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

func (state SelectionUnload) Name() string {
	return "SelectionUnload"
}

// ////////////////////////////////////////////////////////
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

// ////////////////////////////////////////////////////////
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

// ////////////////////////////////////////////////////////
type AddExecDisciplines struct {
}

func (state AddExecDisciplines) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		AddExecID{}.PreviewProcess(ctc)
		return &AddExecID{}
	} else {
		disciplines := strings.Split(messageText, " ")
		var disciplinesBD []int
		uniqueDisciplines := make(map[string]bool)
		for _, val := range disciplines {
			num, err := strconv.Atoi(val)
			if err != nil || (num < 1) || (num > 52) {
				state.PreviewProcess(ctc)
				return &AddExecDisciplines{}
			}
			if uniqueDisciplines[val] {
				continue
			} else {
				disciplinesBD = append(disciplinesBD, num)
				uniqueDisciplines[val] = true
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
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE executors SET disciplines_id = $1 WHERE vk_id=$2", pq.Array(disciplinesBD), VkID)
		if err != nil {
			log.WithError(err).Error("cant add Executor on state ChoiceDiscipline")
			state.PreviewProcess(ctc)
			return &AddExecID{}
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(ctc.User.VkID)
		vkID := strconv.Itoa(VkID)
		b.Message("Исполнитель создан\nID исполнителя: " + vkID + "\nНомера выбранных предметов: " + strings.Trim(strings.Replace(fmt.Sprint(disciplinesBD), " ", " ", -1), "[]"))
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

// ////////////////////////////////////////////////////////
type ManageExecutors struct {
}

func (state ManageExecutors) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else if messageText == "Удалить исполнителя" {
		DeleteExecutorID{}.PreviewProcess(ctc)
		return &DeleteExecutorID{}
	} else if messageText == "Изменить предметы исполнителю" {
		ChangeExecutorsDisciplinesID{}.PreviewProcess(ctc)
		return &ChangeExecutorsDisciplinesID{}
	} else if messageText == "Изменить комиссию для исполнителя" {
		ChangeExecutorsCommissionID{}.PreviewProcess(ctc)
		return &ChangeExecutorsCommissionID{}
	} else {
		state.PreviewProcess(ctc)
		return &ManageExecutors{}
	}
}

func (state ManageExecutors) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	b.Message("Выберите нужный пункт")
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Удалить исполнителя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить предметы исполнителю", "", "secondary")
	k.AddRow()
	k.AddTextButton("Изменить комиссию для исполнителя", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ManageExecutors) Name() string {
	return "ManageExecutors"
}

var ExecutorDeleteID int
var ExecutorChangeID int
var ExecutorChangeCommissionID int

// ////////////////////////////////////////////////////////
type DeleteExecutorID struct {
}

func (state DeleteExecutorID) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		ManageExecutors{}.PreviewProcess(ctc)
		return &ManageExecutors{}
	} else {
		ExecutorDeleteID = 0
		execID, err := strconv.Atoi(messageText)
		ExecutorDeleteID = execID
		if err != nil {
			state.PreviewProcess(ctc)
			return &DeleteExecutorID{}
		} else {
			log.Println(ExecutorDeleteID)
			check, _ := repository.IsExecutor(ctc.Db, ExecutorDeleteID)
			log.Println(check)
			if check {
				DeleteExecutor{}.PreviewProcess(ctc)
				return &DeleteExecutor{}
			} else {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Такого исполнителя не существует. Проверьте ID ВК")
				b.PeerID(ctc.User.VkID)
				_, err := ctc.Vk.MessagesSend(b.Params)
				if err != nil {
					log.Println("Failed to send message on state DeleteExecutorID")
					log.Error(err)
				}
				state.PreviewProcess(ctc)
				return &DeleteExecutorID{}
			}
		}
	}
}

func (state DeleteExecutorID) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID ВК исполнителя")
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

func (state DeleteExecutorID) Name() string {
	return "DeleteExecutorID"
}

// ////////////////////////////////////////////////////////
type DeleteExecutor struct {
}

func (state DeleteExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Нет" {
		ManageExecutors{}.PreviewProcess(ctc)
		return &ManageExecutors{}
	} else if messageText == "Да" {
		_ = repository.DeleteExecutor(ctc.Db, ExecutorDeleteID)
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Исполнитель удалён")
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to send message on state DeleteExecutor")
			log.Error(err)
		}
		ExecutorDeleteID = 0
		ManageExecutors{}.PreviewProcess(ctc)
		return &ManageExecutors{}
	} else {
		state.PreviewProcess(ctc)
		return &DeleteExecutor{}
	}
}

func (state DeleteExecutor) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Удалить исполнителя " + string(ExecutorDeleteID) + "?")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state DeleteExecutor) Name() string {
	return "DeleteExecutor"
}

// ////////////////////////////////////////////////////////
type ChangeExecutorsDisciplinesID struct {
}

func (state ChangeExecutorsDisciplinesID) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		ManageExecutors{}.PreviewProcess(ctc)
		return &ManageExecutors{}
	} else {
		ExecutorChangeID = 0
		execID, err := strconv.Atoi(messageText)
		ExecutorChangeID = execID
		if err != nil {
			state.PreviewProcess(ctc)
			return &ChangeExecutorsDisciplinesID{}
		} else {
			check, _ := repository.IsExecutor(ctc.Db, ExecutorChangeID)
			log.Println(check)
			if check {
				ChangeExecutorsDisciplines{}.PreviewProcess(ctc)
				return &ChangeExecutorsDisciplines{}
			} else {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Такого исполнителя не существует. Проверьте ID ВК")
				b.PeerID(ctc.User.VkID)
				_, err := ctc.Vk.MessagesSend(b.Params)
				if err != nil {
					log.Println("Failed to send message on state DeleteExecutorID")
					log.Error(err)
				}
				state.PreviewProcess(ctc)
				return &ChangeExecutorsDisciplinesID{}
			}
		}
	}
}

func (state ChangeExecutorsDisciplinesID) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID ВК исполнителя")
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

func (state ChangeExecutorsDisciplinesID) Name() string {
	return "ChangeExecutorsDisciplinesID"
}

// ////////////////////////////////////////////////////////
type ChangeExecutorsDisciplines struct {
}

func (state ChangeExecutorsDisciplines) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		ChangeExecutorsDisciplinesID{}.PreviewProcess(ctc)
		return &ChangeExecutorsDisciplinesID{}
	} else {
		disciplines := strings.Split(messageText, " ")
		var disciplinesBD []int
		uniqueDisciplines := make(map[string]bool)
		for _, val := range disciplines {
			num, err := strconv.Atoi(val)
			if err != nil || (num < 1) || (num > 52) {
				state.PreviewProcess(ctc)
				return &ChangeExecutorsDisciplines{}
			}
			if uniqueDisciplines[val] {
				continue
			} else {
				disciplinesBD = append(disciplinesBD, num)
				uniqueDisciplines[val] = true
			}
		}
		_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE executors SET disciplines_id = $1 WHERE vk_id=$2", pq.Array(disciplinesBD), ExecutorChangeID)
		if err != nil {
			log.WithError(err).Error("cant change Executors disciplines on state ChangeExecutorsDisciplines")
			state.PreviewProcess(ctc)
			return &ChangeExecutorsDisciplines{}
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(ctc.User.VkID)
		vkID := strconv.Itoa(ExecutorChangeID)
		b.Message("Предметы изменены\nID исполнителя: " + vkID + "\nНомера выбранных предметов: " + strings.Trim(strings.Replace(fmt.Sprint(disciplinesBD), " ", " ", -1), "[]"))
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	}
}

func (state ChangeExecutorsDisciplines) PreviewProcess(ctc ChatContext) {
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
	b.Message("Введите новые номера предметов для исполнителя")
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ChangeExecutorsDisciplines) Name() string {
	return "ChangeExecutorsDisciplines"
}

// ////////////////////////////////////////////////////////
type ChangeExecutorsCommissionID struct {
}

func (state ChangeExecutorsCommissionID) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		ManageExecutors{}.PreviewProcess(ctc)
		return &ManageExecutors{}
	} else {
		ExecutorChangeCommissionID = 0
		execID, err := strconv.Atoi(messageText)
		ExecutorChangeCommissionID = execID
		if err != nil {
			state.PreviewProcess(ctc)
			return &ChangeExecutorsCommissionID{}
		} else {
			check, _ := repository.IsExecutor(ctc.Db, ExecutorChangeCommissionID)
			log.Println(check)
			if check {
				ChangeExecutorsCommission{}.PreviewProcess(ctc)
				return &ChangeExecutorsCommission{}
			} else {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Такого исполнителя не существует. Проверьте ID ВК")
				b.PeerID(ctc.User.VkID)
				_, err := ctc.Vk.MessagesSend(b.Params)
				if err != nil {
					log.Println("Failed to send message on state DeleteExecutorID")
					log.Error(err)
				}
				state.PreviewProcess(ctc)
				return &ChangeExecutorsCommissionID{}
			}
		}
	}
}

func (state ChangeExecutorsCommissionID) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID ВК исполнителя")
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

func (state ChangeExecutorsCommissionID) Name() string {
	return "ChangeExecutorsCommissionID"
}

// ////////////////////////////////////////////////////////
type ChangeExecutorsCommission struct {
}

func (state ChangeExecutorsCommission) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		ChangeExecutorsCommissionID{}.PreviewProcess(ctc)
		return &ChangeExecutorsCommissionID{}
	} else {
		newCommission, err := strconv.Atoi(messageText)
		if err != nil || newCommission < 0 || newCommission > 100 {
			state.PreviewProcess(ctc)
			return &ChangeExecutorsCommission{}
		}
		var oldCommission *int
		err = ctc.Db.QueryRow("SELECT percent_executor FROM executors WHERE vk_id=$1", ExecutorChangeCommissionID).Scan(&oldCommission)
		if err != nil {
			log.WithError(err).Error("cant change Executors commission_service on state ChangeExecutorsCommission")
			state.PreviewProcess(ctc)
			return &ChangeExecutorsCommission{}
		}
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE executors SET percent_executor = $1 WHERE vk_id=$2", newCommission, ExecutorChangeCommissionID)
		if err != nil {
			log.WithError(err).Error("cant change Executors commission_service on state ChangeExecutorsCommission")
			state.PreviewProcess(ctc)
			return &ChangeExecutorsCommission{}
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(ctc.User.VkID)
		var oldComm string
		if oldCommission == nil {
			oldComm = "0"
		} else {
			oldComm = strconv.Itoa(*oldCommission)
		}
		newComm := strconv.Itoa(newCommission)
		b.Message("Комиссия сервиса изменена с " + oldComm + " до " + newComm)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		b.PeerID(ExecutorChangeCommissionID)
		b.Message("Комиссия сервиса изменена с " + oldComm + " до " + newComm)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	}
}

func (state ChangeExecutorsCommission) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	b.Message("Введите новую комиссию для исполнителя")
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ChangeExecutorsCommission) Name() string {
	return "ChangeExecutorsCommission"
}
