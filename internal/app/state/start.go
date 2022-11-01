package state

import (
	"context"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
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

//////////////////////////////////////////////////////////
type StartState struct {
}

func (state StartState) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Сделать заказ" || messageText == "1" {
		OrderState{}.PreviewProcess(ctc)
		return &OrderState{}
	} else if messageText == "Связаться с исполнителем" || messageText == "2" {

		return &StartState{}
	} else if messageText == "Оставить отзыв" || messageText == "3" {
		return &StartState{}
	} else if messageText == "Сделать заказ через посредника" || messageText == "4" {

		return &StartState{}
	} else if messageText == "Стать исполнителем" || messageText == "5" {

		return &StartState{}
	} else if messageText == "Мои заказы" || messageText == "6" {

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
	b.Message("1. Сделать заказ\n2. Связаться с исполнителем \n3. Оставить отзыв\n4. Сделать заказ через посредника\n5. Стать исполнителем\n6. Мои заказы")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Сделать заказ", "", "secondary")
	k.AddRow()
	k.AddTextButton("Связаться с исполнителем", "", "secondary")
	k.AddRow()
	k.AddTextButton("Оставить отзыв", "", "secondary")
	k.AddRow()
	k.AddTextButton("Сделать заказ через посредника", "", "secondary")
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

/////////////////////////////////////////////////////////
