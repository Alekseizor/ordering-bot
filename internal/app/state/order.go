package state

import (
	"database/sql"
	"github.com/Alekseizor/ordering-bot/internal/app/conversion"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	"unicode/utf8"
)

//////////////////////////////////////////////////////////
type OrderState struct {
}

func (state OrderState) Process(ctc ChatContext, messageText string) State {
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
	b.Message("Выберите дисциплину, нажав на команду «Выбор дисциплины»")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Выбор дисциплины", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад в главное меню", "", "secondary")
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

//////////////////////////////////////////////////////////
type ChoiceDiscipline struct {
}

func (state ChoiceDiscipline) Process(ctc ChatContext, messageText string) State {
	if messageText == "Назад в главное меню" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		messageInt, err := strconv.Atoi(messageText)
		if err != nil {
			state.PreviewProcess(ctc)
			return &ChoiceDiscipline{}
		} else if (messageInt < 1) && (messageInt > 52) {
			state.PreviewProcess(ctc)
			return &ChoiceDiscipline{}
		} else {
			_, err := ctc.Db.ExecContext(*ctc.Ctx, "INSERT INTO orders(customer_vk_id,discipline_id,date_order) VALUES ($1, $2,$3)", ctc.User.VkID, messageInt, time.Now().UTC().Add(time.Hour*3))
			if err != nil {
				log.WithError(err).Error("cant set user")
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
	b.Message("Отправь номер нужной дисциплины")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	b = params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("1. MATLAB\n2. MS Office(Word, Excel, Access)\n3. Mathcad\n4. Аналитическая геометрия\n5. Английский язык\n6. Детали машин\n7. Дискретная математика\n8. Инженерная и компьютерная графика\n9. Интегралы и дифференциальные уравнения\n10. Информатика\n11. История\n12. Кратные интегралы и ряды\n13. Культурология\n14. Линейная алгебра\n15. Математика\n16. Математический анализ\n17. Материаловедение\n18. Менеджмент\n19. Метрология\n20. Механика жидкости и газа\n21. Начертательная геометрия\n22. Организация производства\n23. Основы конструирования приборов\n24. Основы теории цепей\n25. Основы технологии приборостроения\n26. Политология\n27. Правоведение\n28. Практика\n29. Прикладная статистика\n30. Психология\n31. Системный анализ и принятие решений\n32. Сопротивление материалов\n33. Социология\n34. Теоретическая механика\n35. Теоретические основы электротехники\n36. Теория вероятностей\n37. Теория механизмов и машин\n38. Теория поля\n39. Теория функции комплексных переменных и операционное исчисление\n40. Теория функции нескольких переменных\n41. Термодинамика\n42. Технология конструкционных материалов\n43. Уравнения математической физики\n44. Физика\n45. Физкультура\n46. Философия\n47. Финансирование инновационной деятельности\n48. Химия\n49. Цифровые устройства и микропроцессоры\n50. Экономика\n51. Электроника\n52. Электротехника")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад в главное меню", "", "secondary")
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

//////////////////////////////////////////////////////////
const (
	layout = "02.01.2006"
)

type ChoiceDate struct {
}

func (state ChoiceDate) Process(ctc ChatContext, messageText string) State {
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
	b.Message("Выберите дату выполнения заказа")
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

//////////////////////////////////////////////////////////
type ChoiceTime struct {
}

func (state ChoiceTime) Process(ctc ChatContext, messageText string) State {
	ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		state.PreviewProcess(ctc)
		return &ChoiceTime{}
	}
	if messageText == "Вернуться к выбору дня" {
		ChoiceDate{}.PreviewProcess(ctc)
		return &ChoiceDate{}
	} else if utf8.RuneCountInString(messageText) > 4 {
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
	b.Message("Введите время выполнения заказа в формате ЧЧ:ММ")
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

//////////////////////////////////////////////////////////
type ConfirmationOrder struct {
}

func (state ConfirmationOrder) Process(ctc ChatContext, messageText string) State {
	if messageText == "Вернуться к выбору времени" {
		ChoiceTime{}.PreviewProcess(ctc)
		return &ChoiceTime{}
	} else if messageText == "Подтвердить" {
		state.PreviewProcess(ctc)
		return &ConfirmationOrder{}
	} else if messageText == "Добавить комментарий к заказу" {
		state.PreviewProcess(ctc)
		return &ConfirmationOrder{}
	} else {
		state.PreviewProcess(ctc)
		return &ConfirmationOrder{}
	}
}

func (state ConfirmationOrder) PreviewProcess(ctc ChatContext) {
	/*ID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		state.PreviewProcess(ctc)
		return
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	order, err := repository.GetOrder(ctc.Db, ID)
	if err != nil {
		state.PreviewProcess(ctc)
		return
	}
	b.Message("Ваш заказ:\nДисциплина - " + order.DisciplineName)*/
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Подтвердите заказ")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Подтвердить", "", "secondary")
	k.AddTextButton("Вернуться к выбору времени", "", "secondary")
	k.AddTextButton("Добавить комментарий к заказу", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state ConfirmationOrder) Name() string {
	return "ConfirmationOrder"
}
