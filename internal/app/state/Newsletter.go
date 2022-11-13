package state

import (
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
)

// ////////////////////////////////////////////////////////
type Newsletter struct {
}

func (state Newsletter) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	//todo: Проверка - в прикрепленных только файлы или картинки
	fullMSG, _ := ctc.Vk.MessagesGetByID(api.Params{
		"message_ids": msg.ID,
	})

	attachments := fullMSG.Items[0].Attachments

	if messageText == "Назад" {
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else if messageText == "" {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Введите текст рассылки!")
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to send message on state Newsletter")
			log.Error(err)
		}
		state.PreviewProcess(ctc)
		return &Newsletter{}
	} else {
		repository.SetMessage(ctc.Db, messageText)
		if attachments != nil {
			repository.NewsletterWriteUrl(ctc.Db, attachments)
		}
		NewsletterConfirmation{}.PreviewProcess(ctc)
		return &NewsletterConfirmation{}
	}
}

func (state Newsletter) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите текст рассылки и прикрепите до 10 фотографий в одном сообщении")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to send message on state Newsletter")
		log.Error(err)
	}
}

func (state Newsletter) Name() string {
	return "Newsletter"
}

// ////////////////////////////////////////////////////////
type NewsletterConfirmation struct {
}

func (state NewsletterConfirmation) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text

	if messageText == "Нет" {
		Newsletter{}.PreviewProcess(ctc)
		return &Newsletter{}
	} else if messageText == "Да" {
		NewsletterPeerID{}.PreviewProcess(ctc)
		return &NewsletterPeerID{}
	} else {
		state.PreviewProcess(ctc)
		return &NewsletterConfirmation{}
	}
}

func (state NewsletterConfirmation) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	message := repository.GetNewsletterMessage(ctc.Db)
	attachment, _ := repository.NewsletterGetAttachments(ctc.Ctx, ctc.Vk, ctc.Db)
	b.Attachment(attachment)
	b.Message("Вы уверены, что хотите отправить сообщение?\n" + message)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to send message on state Newsletter")
		log.Error(err)
	}
}

func (state NewsletterConfirmation) Name() string {
	return "NewsletterConfirmation"
}

// ////////////////////////////////////////////////////////
type NewsletterPeerID struct {
}

func (state NewsletterPeerID) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text

	if messageText == "Клиентам сервиса (заказчикам)" {
		repository.SetPeerIDsOrders(ctc.Db)
		NewsletterSend{}.PreviewProcess(ctc)
		return &NewsletterSend{}
	} else if messageText == "Исполнителям" {
		repository.SetPeerIDsExecutors(ctc.Db)
		NewsletterSend{}.PreviewProcess(ctc)
		return &NewsletterSend{}
	} else if messageText == "Всем пользователям сервиса" {
		repository.SetPeerIDsAll(ctc.Db)
		NewsletterSend{}.PreviewProcess(ctc)
		return &NewsletterSend{}
	} else {
		state.PreviewProcess(ctc)
		return &NewsletterPeerID{}
	}
}

func (state NewsletterPeerID) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Кому предназначена рассылка?")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Клиентам сервиса (заказчикам)", "", "secondary")
	k.AddRow()
	k.AddTextButton("Исполнителям ", "", "secondary")
	k.AddRow()
	k.AddTextButton("Всем пользователям сервиса", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to send message on state NewsletterPeerID")
		log.Error(err)
	}
}

func (state NewsletterPeerID) Name() string {
	return "NewsletterPeerID"
}

// ////////////////////////////////////////////////////////
type NewsletterSend struct {
}

func (state NewsletterSend) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text

	if messageText == "Нет" {
		NewsletterPeerID{}.PreviewProcess(ctc)
		return &NewsletterPeerID{}
	} else if messageText == "Да" {
		message := repository.GetNewsletterMessage(ctc.Db)
		PeerIDS := repository.GetNewsletterPeerIDs(ctc.Db)
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(message)
		b.PeerIDs(PeerIDS)
		attachment, _ := repository.NewsletterGetAttachments(ctc.Ctx, ctc.Vk, ctc.Db)
		b.Attachment(attachment)
		log.Println(b)
		_, err := ctc.Vk.MessagesSendPeerIDs(b.Params)
		if err != nil {
			log.Println("Failed to send message on state NewsletterPeerID")
			log.Error(err)
		}
		repository.ClearNewsletter(ctc.Db)
		CabinetAdmin{}.PreviewProcess(ctc)
		return &CabinetAdmin{}
	} else {
		state.PreviewProcess(ctc)
		return &NewsletterSend{}
	}
}

func (state NewsletterSend) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Подвердите отправку рассылки")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "secondary")
	k.AddRow()
	k.AddTextButton("Нет", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to send message on state NewsletterPeerID")
		log.Error(err)
	}
}

func (state NewsletterSend) Name() string {
	return "NewsletterSend"
}
