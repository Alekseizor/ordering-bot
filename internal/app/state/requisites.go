package state

import (
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
)

// ////////////////////////////////////////////////////////
type ChangeRequisites struct {
}

func (state ChangeRequisites) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else {
		repository.ChangeRequisites(ctc.Db, messageText)
		message := repository.GetRequisites(ctc.Db)
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Данные изменены\n" + message)
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to send message on state ChangeRequisites")
			log.Error(err)
		}
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	}
}

func (state ChangeRequisites) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите новые реквизиты")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state ChangeRequisites) Name() string {
	return "ChangeRequisites"
}
