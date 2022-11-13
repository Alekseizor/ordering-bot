package state

import (
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"time"
)

//////////////////////////////////////////////////////////
type SelectionDateClear struct {
}

func (state SelectionDateClear) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "За весь период" {
		firstDateStr = "2016-02-12 15:04:05+00"
		secondDateStr = time.Now().UTC().Add(time.Hour * 3).Format(layout2)
		err := repository.ClearTable(ctc.Db, firstDateStr, secondDateStr, close)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Очистка не выполнена")
			b.PeerID(ctc.User.VkID)
			_, err = ctc.Vk.MessagesSend(b.Params)
			if err != nil {
				log.Println("Failed to get record")
				log.Error(err)
				SelectionDateClear{}.PreviewProcess(ctc)
				return &SelectionDateClear{}
			}
			log.Error(err)
			SelectionDateClear{}.PreviewProcess(ctc)
			return &SelectionDateClear{}
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Очистка выполнена успешно!")
		b.PeerID(ctc.User.VkID)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
			SelectionDateClear{}.PreviewProcess(ctc)
			return &SelectionDateClear{}
		}
		SelectionDateClear{}.PreviewProcess(ctc)
		return &SelectionDateClear{}
	} else if messageText == "По дате(с ДД.ММ.ГГ до ДД.ММ.ГГ)" {
		InputFirstDateClear{}.PreviewProcess(ctc)
		return &InputFirstDateClear{}
	} else if messageText == "Назад" {
		SelectionUnload{}.PreviewProcess(ctc)
		return &SelectionUnload{}
	} else {
		state.PreviewProcess(ctc)
		return &SelectionDateClear{}
	}
}

func (state SelectionDateClear) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Выберите период очистки")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("За весь период", "", "secondary")
	k.AddRow()
	k.AddTextButton("По дате(с ДД.ММ.ГГ до ДД.ММ.ГГ)", "", "secondary")
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state SelectionDateClear) Name() string {
	return "SelectionDateClear"
}

//////////////////////////////////////////////////////////
type InputFirstDateClear struct {
}

func (state InputFirstDateClear) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		SelectionDateClear{}.PreviewProcess(ctc)
		return &SelectionDateClear{}
	}
	firstDate, err := time.Parse(layout, messageText)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		state.PreviewProcess(ctc)
		return &InputFirstDateClear{}
	}
	firstDateStr = firstDate.Format(layout2)
	InputSecondDateClear{}.PreviewProcess(ctc)
	return &InputSecondDateClear{}
}

func (state InputFirstDateClear) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите дату, с которой надо начать выборку в формате ДД.ММ.ГГ")
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

func (state InputFirstDateClear) Name() string {
	return "InputFirstDateClear"
}

//////////////////////////////////////////////////////////
type InputSecondDateClear struct {
}

func (state InputSecondDateClear) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		InputFirstDateClear{}.PreviewProcess(ctc)
		return &InputFirstDateClear{}
	}
	secondDate, err := time.Parse(layout, messageText)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		state.PreviewProcess(ctc)
		return &InputSecondDateClear{}
	}
	firstDate, err := time.Parse(layout2, firstDateStr)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		state.PreviewProcess(ctc)
		return &InputSecondDateClear{}
	}
	if firstDate.Before(secondDate) {
		secondDateStr = secondDate.Format(layout2)
		err := repository.ClearTable(ctc.Db, firstDateStr, secondDateStr, close)
		if err != nil {
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Очистка не выполнена")
			b.PeerID(ctc.User.VkID)
			_, err = ctc.Vk.MessagesSend(b.Params)
			if err != nil {
				log.Println("Failed to get record")
				log.Error(err)
				SelectionDateClear{}.PreviewProcess(ctc)
				return &SelectionDateClear{}
			}
			log.Error(err)
			SelectionDateClear{}.PreviewProcess(ctc)
			return &SelectionDateClear{}
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Очистка выполнена успешно!")
		b.PeerID(ctc.User.VkID)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
			InputSecondDateClear{}.PreviewProcess(ctc)
			return &InputSecondDateClear{}
		}
		SelectionDateClear{}.PreviewProcess(ctc)
		return &SelectionDateClear{}
	} else {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Дата не должна предшествовать предыдущей")
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
			state.PreviewProcess(ctc)
			return &InputSecondDateUnload{}
		}
		state.PreviewProcess(ctc)
		return &InputSecondDateClear{}
	}
}

func (state InputSecondDateClear) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите дату, до которой надо делать выборку в формате ДД.ММ.ГГ")
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

func (state InputSecondDateClear) Name() string {
	return "InputSecondDateClear"
}
