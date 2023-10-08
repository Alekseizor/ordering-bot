package state

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/conversion"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	"unicode/utf8"
)

// ////////////////////////////////////////////////////////
type OrderType struct {
}

func (state OrderType) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад в главное меню" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else if messageText == "Рубежный контроль" || messageText == "Домашнее задание" || messageText == "Консультация" || messageText == "Курсовая работа" || messageText == "Экзамен" {
		_, err := ctc.Db.ExecContext(*ctc.Ctx, "INSERT INTO orders(customer_vk_id,type_order,date_order) VALUES ($1, $2,$3)", ctc.User.VkID, messageText, time.Now().UTC().Add(time.Hour*3))
		if err != nil {
			log.WithError(err).Error("cant set order on state OrderType")
			state.PreviewProcess(ctc)
			return &OrderType{}
		}
		ChoiceDiscipline{}.PreviewProcess(ctc)
		return &ChoiceDiscipline{}
	} else {
		state.PreviewProcess(ctc)
		return &OrderType{}
	}
}

func (state OrderType) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("✏Выберите вид работы:")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Рубежный контроль", "", "secondary")
	k.AddTextButton("Домашнее задание", "", "secondary")
	k.AddRow()
	k.AddTextButton("Консультация", "", "secondary")
	k.AddTextButton("Курсовая работа", "", "secondary")
	k.AddRow()
	k.AddTextButton("Экзамен", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад в главное меню", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed send on state OrderType")
		log.Error(err)
	}
}
func (state OrderType) Name() string {
	return "OrderType"
}

// ////////////////////////////////////////////////////////
type OrderState struct {
}

func (state OrderState) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выбор дисциплины" {
		ChoiceDiscipline{}.PreviewProcess(ctc)
		return &ChoiceDiscipline{}
	} else if messageText == "Назад в главное меню" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		state.PreviewProcess(ctc)
		return &OrderState{}
	}
}

func (state OrderState) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("📌Выберите дисциплину, нажав на команду «Выбор дисциплины»")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Выбор дисциплины", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад в главное меню", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state OrderState) Name() string {
	return "OrderState"
}

// ////////////////////////////////////////////////////////
type ChoiceDiscipline struct {
}

func (state ChoiceDiscipline) Process(ctc ChatContext, msg object.MessagesMessage) State {
	ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println(err)
		state.PreviewProcess(ctc)
		return &ChoiceDiscipline{}
	}
	messageText := msg.Text
	if messageText == "Назад" {
		err = repository.DeleteOrder(ctc.Db, ID)
		if err != nil {
			log.Println(err)
			OrderType{}.PreviewProcess(ctc)
			return &OrderType{}
		}
		OrderType{}.PreviewProcess(ctc)
		return &OrderType{}
	} else {
		messageInt, err := strconv.Atoi(messageText)
		if err != nil {
			state.PreviewProcess(ctc)
			return &ChoiceDiscipline{}
		} else if (messageInt < 1) || (messageInt > 52) {
			state.PreviewProcess(ctc)
			return &ChoiceDiscipline{}
		} else {
			_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET discipline_id =$1 WHERE id=$2", messageInt, ID)
			if err != nil {
				log.WithError(err).Error("cant set order on state ChoiceDiscipline")
				state.PreviewProcess(ctc)
				return &ChoiceDiscipline{}
			}
			ChoiceDate{}.PreviewProcess(ctc)
			return &ChoiceDate{}
		}

	}
}

func (state ChoiceDiscipline) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("📌Отправь номер нужной дисциплины")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	b = params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("1. MATLAB\n2. MS Office(Word, excel, Access)\n3. Mathcad\n4. Аналитическая геометрия\n5. Английский язык\n6. Детали машин\n7. Дискретная математика\n8. Инженерная и компьютерная графика\n9. Интегралы и дифференциальные уравнения\n10. Информатика\n11. История\n12. Кратные интегралы и ряды\n13. Культурология\n14. Линейная алгебра\n15. Математика\n16. Математический анализ\n17. Материаловедение\n18. Менеджмент\n19. Метрология\n20. Механика жидкости и газа\n21. Начертательная геометрия\n22. Организация производства\n23. Основы конструирования приборов\n24. Основы теории цепей\n25. Основы технологии приборостроения\n26. Политология\n27. Правоведение\n28. Практика\n29. Прикладная статистика\n30. Психология\n31. Системный анализ и принятие решений\n32. Сопротивление материалов\n33. Социология\n34. Теоретическая механика\n35. Теоретические основы электротехники\n36. Теория вероятностей\n37. Теория механизмов и машин\n38. Теория поля\n39. Теория функции комплексных переменных и операционное исчисление\n40. Теория функции нескольких переменных\n41. Термодинамика\n42. Технология конструкционных материалов\n43. Уравнения математической физики\n44. Физика\n45. Физкультура\n46. Философия\n47. Финансирование инновационной деятельности\n48. Химия\n49. Цифровые устройства и микропроцессоры\n50. Экономика\n51. Электроника\n52. Электротехника")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ChoiceDiscipline) Name() string {
	return "ChoiceDiscipline"
}

// ////////////////////////////////////////////////////////
const (
	layout  = "02.01.2006"
	layout2 = "2006-01-02 15:04:05-07"
)

type ChoiceDate struct {
}

func (state ChoiceDate) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		state.PreviewProcess(ctc)
		return &ChoiceDate{}
	}
	if messageText == "Предыдущий шаг" {
		ChoiceDiscipline{}.PreviewProcess(ctc)
		return &ChoiceDiscipline{}
	} else if messageText == "Свой вариант" {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Введите время выполнения заказа в формате ДД.ММ.ГГГГ")
		b.PeerID(ctc.User.VkID)
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Предыдущий шаг", "", "secondary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		//state.PreviewProcess(ctc)
		return &ChoiceDate{}
	} else if messageText == "Сегодня" || messageText == "Сейчас" {
		_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET date_finish =$1 WHERE id=$2", time.Now().UTC().Add(time.Hour*3), ID)
		if err != nil {
			log.WithError(err).Error("cant set date_finish")
			state.PreviewProcess(ctc)
			return &ChoiceDate{}
		}
		ChoiceTime{}.PreviewProcess(ctc)
		return &ChoiceTime{}
	} else if messageText == "Завтра" {
		_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET date_finish =$1 WHERE id=$2", time.Now().UTC().Add(time.Hour*3).AddDate(0, 0, 1), ID)
		if err != nil {
			log.WithError(err).Error("cant set date_finish")
			state.PreviewProcess(ctc)
			return &ChoiceDate{}
		}
		ChoiceTime{}.PreviewProcess(ctc)
		return &ChoiceTime{}
	} else if messageText == "Через 2 недели" {
		_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET date_finish =$1 WHERE id=$2", time.Now().UTC().Add(time.Hour*3).AddDate(0, 0, 14), ID)
		if err != nil {
			log.WithError(err).Error("cant set date_finish")
			state.PreviewProcess(ctc)
			return &ChoiceDate{}
		}
		ChoiceTime{}.PreviewProcess(ctc)
		return &ChoiceTime{}
	} else if utf8.RuneCountInString(messageText) > 7 {
		if messageText[2] == '.' && messageText[5] == ' ' {
			day, err := strconv.Atoi(messageText[0:2])
			if err != nil {
				log.WithError(err).Error("the string is not formatted per day")
				state.PreviewProcess(ctc)
				return &ChoiceDate{}
			}
			month, err := strconv.Atoi(messageText[3:5])
			if err != nil {
				log.WithError(err).Error("the string is not formatted per month")
				state.PreviewProcess(ctc)
				return &ChoiceDate{}
			}
			weekday := messageText[6:]
			log.Println(day, month, weekday)
			today := time.Now().UTC().Add(time.Hour * 3)
			today = today.AddDate(0, 0, 1) //сместили дату на завтра
			for i := 0; i < 8; i++ {
				today = today.AddDate(0, 0, 1) //смещаем поэтапно на каждый из пяти дней
				if day == today.Day() && month == int(today.Month()) && weekday == conversion.GetWeekDayStr(today) {
					_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET date_finish =$1 WHERE id=$2", today, ID)
					if err != nil {
						log.WithError(err).Error("cant set date_finish")
						state.PreviewProcess(ctc)
						return &ChoiceDate{}
					}
					ChoiceTime{}.PreviewProcess(ctc)
					return &ChoiceTime{}
				}
			}
			state.PreviewProcess(ctc)
			return &ChoiceDate{}
		} else if messageText[2] == '.' && messageText[5] == '.' {
			date, err := time.Parse(layout, messageText)
			if err != nil {
				log.WithError(err).Error("the string is not formatted per date")
				state.PreviewProcess(ctc)
				return &ChoiceDate{}
			}
			if date.After(time.Now().UTC().Add(time.Hour*3).AddDate(0, 0, -1)) {
				_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET date_finish =$1 WHERE id=$2", date, ID)
				if err != nil {
					log.WithError(err).Error("cant set date_finish")
					state.PreviewProcess(ctc)
					return &ChoiceDate{}
				}
				ChoiceTime{}.PreviewProcess(ctc)
				return &ChoiceTime{}
			} else {
				b := params.NewMessagesSendBuilder()
				b.RandomID(0)
				b.Message("Попробуйте ввести время выполнения заказа в формате ДД.ММ.ГГГГ")
				b.PeerID(ctc.User.VkID)
				k := &object.MessagesKeyboard{}
				k.AddRow()
				k.AddTextButton("Предыдущий шаг", "", "secondary")
				_, err := ctc.Vk.MessagesSend(b.Params)
				if err != nil {
					log.Println("Failed to get record")
					log.Error(err)
				}
				return &ChoiceDate{}
			}
		} else {
			state.PreviewProcess(ctc)
			return &ChoiceDate{}
		}
	} else {
		state.PreviewProcess(ctc)
		return &ChoiceDate{}
	}
}

func (state ChoiceDate) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("📅Выберите дату выполнения заказа")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Сегодня", "", "secondary")
	k.AddTextButton("Завтра", "", "secondary")
	//взял московское время
	today := time.Now().UTC().Add(time.Hour * 3)
	today = today.AddDate(0, 0, 1) //сместили дату на завтра
	k.AddRow()
	for i := 0; i < 8; i++ {
		if i == 4 {
			k.AddRow()
		}
		today = today.AddDate(0, 0, 1) //смещаем поэтапно на каждый из пяти дней
		k.AddTextButton(conversion.GetDateStr(today), "", "secondary")
	}
	k.AddRow()
	k.AddTextButton("Через 2 недели", "", "secondary")
	k.AddTextButton("Сейчас", "", "secondary")
	k.AddTextButton("Свой вариант", "", "secondary")
	k.AddRow()
	k.AddTextButton("Предыдущий шаг", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state ChoiceDate) Name() string {
	return "ChoiceDate"
}

// ////////////////////////////////////////////////////////
type ChoiceTime struct {
}

func (state ChoiceTime) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		state.PreviewProcess(ctc)
		return &ChoiceTime{}
	}
	if messageText == "Вернуться к выбору дня" {
		ChoiceDate{}.PreviewProcess(ctc)
		return &ChoiceDate{}
	} else if utf8.RuneCountInString(messageText) == 5 {
		if messageText[2] == ':' {
			hour, err := strconv.Atoi(messageText[0:2])
			if err != nil || hour < 0 || hour > 23 {
				log.Println("z")
				log.WithError(err).Error("the string is not formatted per day")
				state.PreviewProcess(ctc)
				return &ChoiceTime{}
			}
			minute, err := strconv.Atoi(messageText[3:5])
			if err != nil || minute < 0 || minute > 60 {
				log.WithError(err).Error("the string is not formatted per month")
				state.PreviewProcess(ctc)
				return &ChoiceTime{}
			}
			var date time.Time
			err = ctc.Db.QueryRow("SELECT date_finish from orders WHERE customer_vk_id =$1 ORDER BY id DESC LIMIT 1", ctc.User.VkID).Scan(&date)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Println("Row with customer_vk_id unknown")
				} else {
					log.Println("Couldn't find the line with the order")
				}
				log.Error(err)
				state.PreviewProcess(ctc)
				return &ChoiceTime{}
			}
			date = time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.UTC)
			_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET date_finish =$1 WHERE id=$2", date, ID)
			if err != nil {
				log.WithError(err).Error("cant set date_finish")
				state.PreviewProcess(ctc)
				return &ChoiceTime{}
			}
			ConfirmationOrder{}.PreviewProcess(ctc)
			return &ConfirmationOrder{}
		} else {
			state.PreviewProcess(ctc)
			return &ChoiceTime{}
		}
	} else {
		state.PreviewProcess(ctc)
		return &ChoiceTime{}
	}
}

func (state ChoiceTime) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("⏰Введите время выполнения заказа в формате ЧЧ:ММ")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Вернуться к выбору дня", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state ChoiceTime) Name() string {
	return "ChoiceTime"
}

// ////////////////////////////////////////////////////////
type ConfirmationOrder struct {
}

func (state ConfirmationOrder) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Вернуться к выбору времени" {
		ChoiceTime{}.PreviewProcess(ctc)
		return &ChoiceTime{}
	} else if messageText == "Подтвердить" {
		ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET customers_comment =$1 WHERE id=$2", nil, ID)
		if err != nil {
			log.WithError(err).Error("cant record users comment")
			state.PreviewProcess(ctc)
			return &ConfirmationOrder{}
		}
		TaskOrder{}.PreviewProcess(ctc)
		return &TaskOrder{}
	} else if messageText == "Добавить комментарий к заказу" {
		CommentOrder{}.PreviewProcess(ctc)
		return &CommentOrder{}
	} else {
		state.PreviewProcess(ctc)
		return &ConfirmationOrder{}
	}
}

func (state ConfirmationOrder) PreviewProcess(ctc ChatContext) {
	output, err := repository.GetCompleteOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("Failed to get orders output")
		log.Error(err)
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(output)
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Подтвердить", "", "secondary")
	k.AddTextButton("Вернуться к выбору времени", "", "secondary")
	k.AddTextButton("Добавить комментарий к заказу", "", "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to send order: state ConfirmationOrder")
		log.Error(err)
	}
}
func (state ConfirmationOrder) Name() string {
	return "ConfirmationOrder"
}

// ////////////////////////////////////////////////////////
type CommentOrder struct {
}

func (state CommentOrder) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		ConfirmationOrder{}.PreviewProcess(ctc)
		return &ConfirmationOrder{}
	} else {
		if utf8.RuneCountInString(messageText) > 150 {
			log.Println("Text is to large")
			CommentOrder{}.PreviewProcess(ctc)
			return &CommentOrder{}
		}
		ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET customers_comment =$1 WHERE id=$2", messageText, ID)
		if err != nil {
			log.WithError(err).Error("cant record users comment")
		}

		TaskOrder{}.PreviewProcess(ctc)
		return &TaskOrder{}
	}
}

func (state CommentOrder) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Ограничение на комментарий - 150 символов")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	//k.AddTextButton("Отправить комментарий", "", "secondary")
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state CommentOrder) Name() string {
	return "CommentOrder"
}

// ////////////////////////////////////////////////////////
type TaskOrder struct {
}

func (state TaskOrder) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	//todo: Проверка - в прикрепленных только файлы или картинки
	fullMSG, _ := ctc.Vk.MessagesGetByID(api.Params{
		"message_ids": msg.ID,
	})

	attachments := fullMSG.Items[0].Attachments
	if attachments != nil {
		repository.WriteUrl(ctc.Db, ctc.User.VkID, attachments)
	}

	if messageText == "Назад" {
		ConfirmationOrder{}.PreviewProcess(ctc)
		return &ConfirmationOrder{}
	} else {
		ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "UPDATE orders SET order_task =$1 WHERE id=$2", messageText, ID)
		if err != nil {
			log.WithError(err).Error("cant record users comment")
			state.PreviewProcess(ctc)
			return &TaskOrder{}
		}
		OrderCompleted{}.PreviewProcess(ctc)
		return &OrderCompleted{}
	}
}

func (state TaskOrder) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("📎Отправьте фото,текстовое описание или документ задания (любой формат) одним сообщением!")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)

	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state TaskOrder) Name() string {
	return "TaskOrder"
}

// ////////////////////////////////////////////////////////
type ConfirmExecutor struct {
}

func (state ConfirmExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		OrderCompleted{}.PreviewProcess(ctc)
		return &OrderCompleted{}
	}
	executorID, err := strconv.Atoi(messageText)
	if err != nil {
		log.WithError(err).Error("the string is not number for executor id")
		state.PreviewProcess(ctc)
		return &ConfirmExecutor{}
	}
	isExec, err := repository.IsExecutorByID(ctc.Db, executorID)
	if isExec {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Исполнитель найден. Заказ отправлен.")
		b.PeerID(ctc.User.VkID)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		err = DirectDistribution(ctc, executorID)
		if err != nil {
			log.Println("Failed to send direct offer")
			log.Error(err)
		}
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Исполнитель с таким ID не найден")
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		state.PreviewProcess(ctc)
		return &ConfirmExecutor{}
	}
}

func (state ConfirmExecutor) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID исполнителя")
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state ConfirmExecutor) Name() string {
	return "ConfirmExecutor"
}

// ////////////////////////////////////////////////////////
type OrderCompleted struct {
}

func (state OrderCompleted) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Оформить заказ" {
		err := DistributionOrderExecutors(ctc)
		if err != nil {
			log.Println("the order could not be sent to the executors")
		}
		StartState{}.PreviewProcess(ctc)
		return &StartState{}

	} else if messageText == "Редактировать заказ" {
		OrderChange{}.PreviewProcess(ctc)
		return &OrderChange{}

	} else if messageText == "Отменить заказ" {
		OrderCancel{}.PreviewProcess(ctc)
		return &OrderCancel{}

	} else if messageText == "Выбрать конкретного исполнителя" {
		ConfirmExecutor{}.PreviewProcess(ctc)
		return &ConfirmExecutor{}

	} else {
		state.PreviewProcess(ctc)
		return &OrderCompleted{}
	}
}

func (state OrderCompleted) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Информация получена. Ваш заказ загружается")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record OrderCompleted")
		log.Error(err)
	}
	output, err := repository.GetCompleteOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("Failed to get orders output")
		log.Error(err)
	}
	b.Message(output)
	attachment, _ := repository.GetAttachments(ctc.Vk, ctc.Db, ctc.User.VkID)
	log.Println("вывод - " + output)
	b.Attachment(attachment)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Оформить заказ", "", "secondary")
	k.AddTextButton("Редактировать заказ", "", "secondary")
	k.AddTextButton("Отменить заказ", "", "secondary")
	k.AddRow()
	k.AddTextButton("Выбрать конкретного исполнителя", "", "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record OrderCompleted 2")
		log.Error(err)
	}
}
func (state OrderCompleted) Name() string {
	return "OrderCompleted"
}

// ////////////////////////////////////////////////////////
type OrderCancel struct {
}

func (state OrderCancel) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Да" {
		ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
		if err != nil {
			log.WithError(err).Error("cant get order id")
			state.PreviewProcess(ctc)
			return &OrderCancel{}
		}
		_, err = ctc.Db.ExecContext(*ctc.Ctx, "DELETE FROM orders WHERE id=$1", ID)
		if err != nil {
			log.WithError(err).Error("cant delete order")
			state.PreviewProcess(ctc)
			return &OrderCancel{}
		}
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		OrderCompleted{}.PreviewProcess(ctc)
		return &OrderCompleted{}
	}
}

func (state OrderCancel) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Вы действительно хотите отменить заказ?")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "positive")
	k.AddTextButton("Нет", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state OrderCancel) Name() string {
	return "OrderCancel"
}

// ////////////////////////////////////////////////////////
type OrderChange struct {
}

func (state OrderChange) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		OrderCompleted{}.PreviewProcess(ctc)
		return &OrderCompleted{}
	} else if messageText == "Вид работы" {
		EditType{}.PreviewProcess(ctc)
		return &EditType{}
	} else if messageText == "Вид дисциплины" {
		EditDiscipline{}.PreviewProcess(ctc)
		return &EditDiscipline{}
	} else if messageText == "Дата исполнения заказа" {
		EditDate{}.PreviewProcess(ctc)
		return &EditDate{}
	} else if messageText == "Время исполнения заказа" {
		EditTime{}.PreviewProcess(ctc)
		return &EditTime{}
	} else if messageText == "Информация по заказу" {
		EditTaskOrder{}.PreviewProcess(ctc)
		return &EditTaskOrder{}
	} else if messageText == "Комментарий к заказу" {
		EditCommentOrder{}.PreviewProcess(ctc)
		return &EditCommentOrder{}
	} else {
		OrderCompleted{}.PreviewProcess(ctc)
		return &OrderCompleted{}
	}
}

func (state OrderChange) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите пункт для редактирования")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Вид работы", "", "secondary")
	//k.AddRow()
	k.AddTextButton("Вид дисциплины", "", "secondary")
	k.AddRow()
	k.AddTextButton("Дата исполнения заказа", "", "secondary")
	//k.AddRow()
	k.AddTextButton("Время исполнения заказа", "", "secondary")
	k.AddRow()
	k.AddTextButton("Информация по заказу", "", "secondary")
	//k.AddRow()
	k.AddTextButton("Комментарий к заказу", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state OrderChange) Name() string {
	return "OrderChange"
}
