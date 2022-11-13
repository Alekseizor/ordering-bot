package state

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
)

var firstDateStr, secondDateStr, close string

//////////////////////////////////////////////////////////
type CabinetAdmin struct {
}

func (state CabinetAdmin) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Выгрузка таблицы заказов" {
		UnloadTable{}.PreviewProcess(ctc)
		return &UnloadTable{}
	} else if messageText == "Назначить исполнителя" {
		WriteAdmin{}.PreviewProcess(ctc)
		return &WriteAdmin{}
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

//////////////////////////////////////////////////////////
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
