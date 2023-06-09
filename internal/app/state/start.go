package state

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type ChatContext struct {
	User *ds.User
	Vk   *api.VK
	Ctx  *context.Context
	Db   *sqlx.DB
}

type State interface {
	Name() string                                      //получаем название состояния в виде строки, чтобы в дальнейшем куда-то записать(БД)
	Process(ChatContext, object.MessagesMessage) State //нужно взять контекст, посмотреть на каком состоянии сейчас пользователь, метод должен вернуть состояние
	PreviewProcess(ctc ChatContext)
}

// ////////////////////////////////////////////////////////
type StartState struct {
}

func (state StartState) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Сделать заказ" || messageText == "1" {
		OrderType{}.PreviewProcess(ctc)
		return &OrderType{}
	} else if messageText == "Стать исполнителем" || messageText == "2" {
		if ctc.User.VkID == config.FromContext(*ctc.Ctx).AdminID {
			CabinetAdmin{}.PreviewProcess(ctc)
			return &CabinetAdmin{}
		}
		BecomeExecutor{}.PreviewProcess(ctc)
		return &BecomeExecutor{}
	} else if messageText == "Мои заказы" || messageText == "3" {
		MyOrderState{}.PreviewProcess(ctc)
		return &StartState{}
	} else {
		state.PreviewProcess(ctc)
		return &StartState{}
	}
}

func (state StartState) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Пришли нужный номер команды или воспользуйся кнопками")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	b = params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("1. Сделать заказ\n2. Стать исполнителем\n3. Мои заказы")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Сделать заказ", "", "secondary")
	k.AddRow()
	k.AddTextButton("Стать исполнителем", "", "secondary")
	k.AddRow()
	k.AddTextButton("Мои заказы", "", "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state StartState) Name() string {
	return "StartState"
}

// ///////////////////////////////////////////////////////
type MyOrderState struct {
}

func (state MyOrderState) Process(ctc ChatContext, msg object.MessagesMessage) State {
	return MyOrderState{}
}

func (state MyOrderState) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	ordersID, err := repository.GetOrdersIDUser(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("Failed to GetOrdersIDUser")
		return
	}
	b.PeerID(ctc.User.VkID)
	if len(ordersID) == 0 {
		b.Message("Вы еще не оформили ни одного заказа")
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
	}
	for _, orderID := range ordersID {
		output, err := repository.GetCompleteOrders(ctc.Db, orderID)
		if err != nil {
			log.Println("Failed to get orders output")
			log.Error(err)
		}
		b.Message(output)
		attachment, err := repository.GetAttachmentsMyOrder(ctc.Vk, ctc.Db, orderID, ctc.User.VkID)
		if err != nil {
			log.Println("Failed to GetAttachmentsMyOrder")
		}
		b.Attachment(attachment)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
	}
}

func (state MyOrderState) Name() string {
	return "MyOrderState"
}

/////////////////////////////////////////////////////////
