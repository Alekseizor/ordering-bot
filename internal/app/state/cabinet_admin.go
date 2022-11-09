package state

import (
	"github.com/Alekseizor/ordering-bot/internal/app/excel"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
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
	messageText := msg.Text
	if messageText == "Да" {
		table := excel.CreateRespTableOneExec(ctc.Db, "12.11.2006", "12.11.2099", ctc.User.VkID)

	} else if messageText == "Нет" {
		table := excel.CreateRespTable(ctc.Db, "12.11.2006", "12.11.2099")
		UnloadingUnclose{}.PreviewProcess(ctc)
		return &UnloadingUnclose{}
	} else if messageText == "Назад" {
		UnloadingUnclose{}.PreviewProcess(ctc)
		return &UnloadingUnclose{}
	} else {
		state.PreviewProcess(ctc)
		return &AllUnloadingUnclose{}
	}
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
