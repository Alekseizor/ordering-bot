package state

import (
	"database/sql"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
)

//////////////////////////////////////////////////////////
type BecomeExecutor struct {
}

func (state BecomeExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "История заказов" {
		ExecHistoryOrders{}.PreviewProcess(ctc)
		return &ExecHistoryOrders{}
	} else if messageText == "Написать администратору" {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Напишите администратору проекта по ссылке: https://vk.com/bitchpart")
		b.PeerID(ctc.User.VkID)
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Назад", "", "secondary")
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		return &BecomeExecutor{}
	} else if messageText == "Назад" {
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		state.PreviewProcess(ctc)
		return &BecomeExecutor{}
	}
}

func (state BecomeExecutor) PreviewProcess(ctc ChatContext) {
	var ID int
	err := ctc.Db.QueryRow("SELECT vk_id from executors WHERE vk_id =$1", ctc.User.VkID).Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Чтобы стать исполнителем, напишите администратору проекта по ссылке: https://vk.com/bitchpart")
			b.PeerID(ctc.User.VkID)
			k := &object.MessagesKeyboard{}
			k.AddRow()
			k.AddTextButton("Назад", "", "secondary")
			b.Keyboard(k)
			_, err = ctc.Vk.MessagesSend(b.Params)
			if err != nil {
				log.Println("Failed to get record")
				log.Error(err)
			}
		} else {
			log.Println("Couldn't find the line with the order")
		}
		log.Error(err)
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выбери нужный пункт")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("История заказов", "", "secondary")
	k.AddRow()
	k.AddTextButton("Написать администратору", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state BecomeExecutor) Name() string {
	return "BecomeExecutor"
}

//////////////////////////////////////////////////////////
type ExecHistoryOrders struct {
}

func (state ExecHistoryOrders) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text

	if messageText == "Войти в личный кабинет" {
		OrderState{}.PreviewProcess(ctc)
		return &OrderState{}
	} else if messageText == "Назад" {
		BecomeExecutor{}.PreviewProcess(ctc)
		return &BecomeExecutor{}
	} else {
		state.PreviewProcess(ctc)
		return &ExecHistoryOrders{}
	}
}

func (state ExecHistoryOrders) PreviewProcess(ctc ChatContext) {
	var numberLines, rating int
	err := ctc.Db.QueryRow("SELECT COUNT(*) from orders WHERE executor_vk_id =$1", ctc.User.VkID).Scan(&numberLines)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with id unknown")
		} else {
			log.Println("Couldn't find the line with the order")
		}
		log.Error(err)
	}
	err = ctc.Db.QueryRow("SELECT rating from executors WHERE vk_id =$1", ctc.User.VkID).Scan(&rating)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with id unknown")
		} else {
			log.Println("Couldn't find the line with the order")
		}
		log.Error(err)
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	numberLinesStr := strconv.Itoa(numberLines)
	ratingStr := strconv.Itoa(rating)
	resp := "Количество заказов: " + numberLinesStr + "\nСредняя оценка: " + ratingStr + "\nДоход за всё время: \n"
	b.Message(resp)
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ExecHistoryOrders) Name() string {
	return "ExecHistoryOrders"
}
