package state

import (
	"github.com/Alekseizor/ordering-bot/internal/app/conversion"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
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
			_, err := ctc.Db.ExecContext(*ctc.Ctx, "INSERT INTO orders(customer_vk_id,discipline_id) VALUES ($1, $2)", ctc.User.VkID, messageInt)
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
type ChoiceDate struct {
}

func (state ChoiceDate) Process(ctc ChatContext, messageText string) State {
	if messageText == "Выбор дисциплины" {
		ChoiceDiscipline{}.PreviewProcess(ctc)
		return &ChoiceDiscipline{}
	} else if messageText == "Назад в главное меню" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
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
	k.AddRow()
	k.AddTextButton("Завтра", "", "secondary")
	//взял московское время
	today := time.Now().UTC().Add(time.Hour * 3)
	today = today.AddDate(0, 0, 1) //сместили дату на завтра
	for i := 0; i < 5; i++ {
		today = today.AddDate(0, 0, 1) //смещаем поэтапно на каждый из пяти дней
		k.AddRow()
		k.AddTextButton(conversion.GetDateStr(today), "", "secondary")
	}
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
