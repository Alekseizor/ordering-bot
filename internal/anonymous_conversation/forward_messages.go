package anonymous_conversation

import (
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/Alekseizor/ordering-bot/internal/app/state"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func GetOrderAndVkID(ctc state.ChatContext, msg object.MessagesMessage) (order, vkId string, err error) {

	chat, err := ctc.Vk.MessagesGetConversationsByID(api.Params{
		"peer_ids": msg.PeerID,
	})
	title := strings.Split(chat.Items[0].ChatSettings.Title, "_")
	return title[0], title[1], nil
}

// ////////////////////////////////////////////////////////
type ForwardMessage struct {
}

func (state ForwardMessage) Process(ctc state.ChatContext, msg object.MessagesMessage) state.State {
	state.PreviewProcess(ctc)
	return ConversationSend{}

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

// ////////////////////////////////////////////////////////
type ConversationSend struct {
}

func (state ConversationSend) Process(ctc state.ChatContext, msg object.MessagesMessage) state.State {
	messageText := msg.Text
	var chatID int
	var err error

	if messageText == "Проблемы с заказом" || strings.Contains(messageText, "] Проблемы с заказом") {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Напишите администратору: https://vk.com/bitchpart")
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to send message")
			log.Error(err)
		}
		return &ConversationSend{}
	} else if messageText == "Завершить заказ" || strings.Contains(messageText, "] Завершить заказ") {
		return ConversationSend{}
	} else {
		orderID, vkId, _ := GetOrderAndVkID(ctc, msg)
		vkid, _ := strconv.Atoi(vkId)
		orderId, _ := strconv.Atoi(orderID)
		order, _ := repository.GetOrder(ctc.Db, orderId)
		isExecutor, _ := repository.IsExecutorInOrder(ctc.Db, orderId, vkid)
		chatID, err = repository.GetConversationID(ctc.Db, orderId, isExecutor)
		if err != nil {
			log.Println("Can`t find conversation ID")
		}
		var PeerID uint
		if isExecutor {
			PeerID = order.CustomerVkID
		} else {
			PeerID = *order.ExecutorVkID
		}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message(messageText)
		if len(msg.Attachments) > 0 {
			docsURL, docsTitle, imagesURL := repository.ConversationWriteUrl(msg.Attachments)
			log.Println(ctc.User.VkID)
			log.Println(imagesURL)
			attachment, _ := repository.ConversationGetAttachments(ctc.Vk, int(PeerID), docsURL, docsTitle, imagesURL)
			log.Println(attachment)
			b.Attachment(attachment)

		}
		b.PeerID(chatID)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
		return ConversationSend{}
	}
}

func (state ConversationSend) PreviewProcess(ctc state.ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
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

func (state ConversationSend) Name() string {
	return "ConversationSend"
}
