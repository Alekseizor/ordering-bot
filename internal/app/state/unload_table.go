package state

import (
	"bytes"
	"github.com/Alekseizor/ordering-bot/internal/app/excel"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

//////////////////////////////////////////////////////////
type SelectionDateUnload struct {
}

func (state SelectionDateUnload) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "За весь период" {
		firstDateStr = "2016-02-12 15:04:05+03"
		secondDateStr = time.Now().UTC().Add(time.Hour * 3).Format(layout2)
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else if messageText == "По дате(с ДД.ММ.ГГ до ДД.ММ.ГГ)" {
		InputFirstDateUnload{}.PreviewProcess(ctc)
		return &InputFirstDateUnload{}
	} else if messageText == "За месяц" {
		secondDateStr = time.Now().UTC().Add(time.Hour * 3).Format(layout2)
		firstDateStr = time.Now().UTC().Add(time.Hour*3).AddDate(0, -1, 0).Format(layout2)
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else if messageText == "За 7 дней" {
		secondDateStr = time.Now().UTC().Add(time.Hour * 3).Format(layout2)
		firstDateStr = time.Now().UTC().Add(time.Hour*3).AddDate(0, 0, -7).Format(layout2)
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else if messageText == "За день" {
		secondDateStr = time.Now().UTC().Add(time.Hour * 3).Format(layout2)
		firstDateStr = time.Now().UTC().Add(time.Hour*3).AddDate(0, 0, -1).Format(layout2)
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else if messageText == "Назад" {
		SelectionUnload{}.PreviewProcess(ctc)
		return &SelectionUnload{}
	} else {
		state.PreviewProcess(ctc)
		return &SelectionDateUnload{}
	}
}

func (state SelectionDateUnload) PreviewProcess(ctc ChatContext) {
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

func (state SelectionDateUnload) Name() string {
	return "SelectionDateUnload"
}

//////////////////////////////////////////////////////////
type InputFirstDateUnload struct {
}

func (state InputFirstDateUnload) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		SelectionDateUnload{}.PreviewProcess(ctc)
		return &SelectionDateUnload{}
	}
	firstDate, err := time.Parse(layout, messageText)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		state.PreviewProcess(ctc)
		return &InputFirstDateUnload{}
	}
	firstDateStr = firstDate.Format(layout2)
	InputSecondDateUnload{}.PreviewProcess(ctc)
	return &InputSecondDateUnload{}
}

func (state InputFirstDateUnload) PreviewProcess(ctc ChatContext) {
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

func (state InputFirstDateUnload) Name() string {
	return "InputFirstDateUnload"
}

//////////////////////////////////////////////////////////
type InputSecondDateUnload struct {
}

func (state InputSecondDateUnload) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		InputFirstDateUnload{}.PreviewProcess(ctc)
		return &InputFirstDateUnload{}
	}
	secondDate, err := time.Parse(layout, messageText)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		state.PreviewProcess(ctc)
		return &InputSecondDateUnload{}
	}
	firstDate, err := time.Parse(layout2, firstDateStr)
	if err != nil {
		log.WithError(err).Error("the string is not formatted per date")
		state.PreviewProcess(ctc)
		return &InputSecondDateUnload{}
	}
	if firstDate.Before(secondDate) {
		secondDateStr = secondDate.Format(layout2)
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
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
		return &InputSecondDateUnload{}
	}
}

func (state InputSecondDateUnload) PreviewProcess(ctc ChatContext) {
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

func (state InputSecondDateUnload) Name() string {
	return "InputSecondDateUnload"
}

//////////////////////////////////////////////////////////
type SelectionAllOrPersonalUnload struct {
}

func (state SelectionAllOrPersonalUnload) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Да" {
		PersonalUnload{}.PreviewProcess(ctc)
		return &PersonalUnload{}
	} else if messageText == "Нет" {
		table, err := excel.CreateRespTable(ctc.Db, firstDateStr, secondDateStr, close)
		if err != nil {
			log.Println(err)
			state.PreviewProcess(ctc)
			return &SelectionAllOrPersonalUnload{}
		}
		tableBuffer, err := table.WriteToBuffer()
		if err != nil {
			log.Println(err)
			state.PreviewProcess(ctc)
			return &SelectionAllOrPersonalUnload{}
		}
		tableBytes := tableBuffer.Bytes()
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Ваша таблица:")
		b.PeerID(ctc.User.VkID)
		doc, err := ctc.Vk.UploadMessagesDoc(ctc.User.VkID, "doc", "Book1.xlsx", "", bytes.NewReader(tableBytes))
		if err != nil {
			log.Error(err)
			state.PreviewProcess(ctc)
			return &SelectionAllOrPersonalUnload{}
		}
		b.Attachment(doc.Type + strconv.Itoa(doc.Doc.OwnerID) + "_" + strconv.Itoa(doc.Doc.ID))
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else if messageText == "Назад" {
		SelectionDateUnload{}.PreviewProcess(ctc)
		return &SelectionDateUnload{}
	} else {
		state.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	}
}

func (state SelectionAllOrPersonalUnload) PreviewProcess(ctc ChatContext) {
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

func (state SelectionAllOrPersonalUnload) Name() string {
	return "SelectionAllOrPersonalUnload"
}

//////////////////////////////////////////////////////////
type PersonalUnload struct {
}

func (state PersonalUnload) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else if ID, err := strconv.Atoi(messageText); err == nil {
		table, err := excel.CreateRespTable(ctc.Db, firstDateStr, secondDateStr, close, ID)
		if err != nil {
			log.Println(err)
			state.PreviewProcess(ctc)
			return &PersonalUnload{}
		}
		tableBuffer, err := table.WriteToBuffer()
		if err != nil {
			log.Println(err)
			state.PreviewProcess(ctc)
			return &PersonalUnload{}
		}
		tableBytes := tableBuffer.Bytes()
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Ваша таблица:")
		b.PeerID(ctc.User.VkID)
		log.Println(table)
		doc, err := ctc.Vk.UploadMessagesDoc(ctc.User.VkID, "doc", "Book1.xlsx", "", bytes.NewReader(tableBytes))
		if err != nil {
			log.Error(err)
			state.PreviewProcess(ctc)
			return &PersonalUnload{}
		}
		b.Attachment(doc.Type + strconv.Itoa(doc.Doc.OwnerID) + "_" + strconv.Itoa(doc.Doc.ID))
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
			state.PreviewProcess(ctc)
			return &PersonalUnload{}
		}
		state.PreviewProcess(ctc)
		return &PersonalUnload{}
	} else if messageText == "Назад" {
		SelectionAllOrPersonalUnload{}.PreviewProcess(ctc)
		return &SelectionAllOrPersonalUnload{}
	} else {
		state.PreviewProcess(ctc)
		return &PersonalUnload{}
	}
}

func (state PersonalUnload) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите ID исполнителя")
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

func (state PersonalUnload) Name() string {
	return "PersonalUnload"
}
