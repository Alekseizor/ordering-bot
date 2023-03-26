package anonymous_conversation

import (
	"github.com/Alekseizor/ordering-bot/internal/app/state"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
)

// ////////////////////////////////////////////////////////
type ForwardMessage struct {
}

func (state ForwardMessage) Process(ctc state.ChatContext, msg object.MessagesMessage) state.State {
	state.PreviewProcess(ctc)
	return ForwardMessage{}
}

func (state ForwardMessage) PreviewProcess(ctc state.ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Привет! Теперь ты в анонимном чате! Напиши сообщение и оно сразу же придет заказчику (исполнителю)")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Проблемы с заказом", "", "secondary")
	k.AddRow()
	k.AddTextButton("Завершить заказ", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ForwardMessage) Name() string {
	return "ForwardMessage"
}
