package state

import (
	"database/sql"
	"fmt"
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
	if messageText == "Изменить реквизиты" {
		ChangeRequisiteExecutor{}.PreviewProcess(ctc)
		return &ChangeRequisiteExecutor{}
	}
	if messageText == "История заказов" {
		ExecHistoryOrders{}.PreviewProcess(ctc)
		return &ExecHistoryOrders{}
	} else if messageText == "Написать администратору" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
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
	} else {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Выбери нужный пункт")
		b.PeerID(ctc.User.VkID)
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Изменить реквизиты", "", "secondary")
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
}

func (state BecomeExecutor) Name() string {
	return "BecomeExecutor"
}

//////////////////////////////////////////////////////////
type ExecHistoryOrders struct {
}

func (state ExecHistoryOrders) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text

	if messageText == "Назад" {
		BecomeExecutor{}.PreviewProcess(ctc)
		return &BecomeExecutor{}
	} else {
		state.PreviewProcess(ctc)
		return &ExecHistoryOrders{}
	}
}

func (state ExecHistoryOrders) PreviewProcess(ctc ChatContext) {
	var numberLines int
	var rating, profit float32
	err := ctc.Db.QueryRow("SELECT amount_orders from executors WHERE vk_id =$1", ctc.User.VkID).Scan(&numberLines)
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
	err = ctc.Db.QueryRow("SELECT profit from executors WHERE vk_id =$1", ctc.User.VkID).Scan(&profit)
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
	ratingStr := fmt.Sprint(rating)
	profitStr := fmt.Sprint(profit)
	resp := "Количество заказов: " + numberLinesStr + "\nСредняя оценка: " + ratingStr + "\nДоход за всё время: " + profitStr
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

//////////////////////////////////////////////////////////
type WriteAdmin struct {
}

func (state WriteAdmin) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text

	if messageText == "Назад" {
		BecomeExecutor{}.PreviewProcess(ctc)
		return &BecomeExecutor{}
	} else {
		state.PreviewProcess(ctc)
		return &WriteAdmin{}
	}
}

func (state WriteAdmin) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите администратору проекта по ссылке: https://vk.com/bitchpart")
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

func (state WriteAdmin) Name() string {
	return "WriteAdmin"
}

//////////////////////////////////////////////////////////
type ChangeRequisiteExecutor struct {
}

func (state ChangeRequisiteExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		BecomeExecutor{}.PreviewProcess(ctc)
		return &BecomeExecutor{}
	} else {
		_, err := ctc.Db.ExecContext(*ctc.Ctx, "UPDATE executors SET requisite =$1 WHERE vk_id=$2", messageText, ctc.User.VkID)
		if err != nil {
			log.WithError(err).Error("cant set order on state ChangeRequisiteExecutor")
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Реквизиты успешно изменены!")
		b.PeerID(ctc.User.VkID)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		BecomeExecutor{}.PreviewProcess(ctc)
		return &BecomeExecutor{}
	}
}

func (state ChangeRequisiteExecutor) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Напишите одним сообщением свои реквизиты.\n Например:\nСбербанк:1111222233334444\nТинькофф:5555666677778888")
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

func (state ChangeRequisiteExecutor) Name() string {
	return "ChangeRequisiteExecutor"
}
