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

	title := strings.TrimPrefix(chat.Items[0].ChatSettings.Title, "Заказ номер №")
	title = strings.TrimSuffix(title, "⚠")
	res := strings.Split(title, "_")
	return res[0], res[1], nil
}

// ////////////////////////////////////////////////////////
type ForwardMessage struct {
}

func (state ForwardMessage) Process(ctc state.ChatContext, msg object.MessagesMessage) state.State {
	state.PreviewProcess(ctc)
	var message string
	orderID, vkId, _ := GetOrderAndVkID(ctc, msg)
	vkid, _ := strconv.Atoi(vkId)
	orderId, _ := strconv.Atoi(orderID)
	order, _ := repository.GetOrder(ctc.Db, orderId)
	idExec, _ := repository.GetExecutor(ctc.Db, int(*order.ExecutorVkID))
	idExecutor := strconv.Itoa(idExec.Id)
	isExecutor, _ := repository.IsExecutorInOrder(ctc.Db, orderId, vkid)
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	if isExecutor {
		message = "Начните исполнять заказ лишь после того, как получите фото с подтверждением оплаты"
	} else {
		message = "Пришлите фото с подтверждением оплаты"
	}
	b.Message(message)
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to send message")
		log.Error(err)
	}

	if !isExecutor {
		b.RandomID(0)
		req := repository.GetRequisites(ctc.Db)
		message = "Реквизиты для оплаты: " + req + "\nID вашего исполнителя - " + idExecutor
		b.Message(message)
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to send message")
			log.Error(err)
		}
	}

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
		FinishOrderCheck{}.PreviewProcess(ctc)
		return FinishOrderCheck{}
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
		state.PreviewProcess(ctc)
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

// ////////////////////////////////////////////////////////
type FinishOrderCheck struct {
}

func (state FinishOrderCheck) Process(ctc state.ChatContext, msg object.MessagesMessage) state.State {
	messageText := msg.Text
	if messageText == "Да" || strings.Contains(messageText, "] Да") {
		orderID, vkId, _ := GetOrderAndVkID(ctc, msg)
		vkid, _ := strconv.Atoi(vkId)
		orderId, _ := strconv.Atoi(orderID)
		isExecutor, _ := repository.IsExecutorInOrder(ctc.Db, orderId, vkid)
		_ = repository.FinishOrder(ctc.Db, orderId, isExecutor)
		chat, err := ctc.Vk.MessagesGetConversationsByID(api.Params{
			"peer_ids": msg.PeerID,
		})
		chatID, err := repository.GetConversationID(ctc.Db, orderId, !isExecutor)
		if err != nil {
			log.Println("Can`t find conversation ID")
		}
		title := strings.TrimSuffix(chat.Items[0].ChatSettings.Title, "⚠")
		_, err = ctc.Vk.MessagesEditChat(api.Params{
			"chat_id": chatID - 2000000000,
			"title":   title + "✅",
		})
		if err != nil {
			log.Println("Failed to edit chat title")
			log.Println(chatID)
			log.Error(err)
		}
		FinishOrder{}.PreviewProcess(ctc)
		return FinishOrder{}
	} else if messageText == "Нет" || strings.Contains(messageText, "] Нет") {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(ctc.User.VkID)
		b.Message("Завершение заказа отменено")
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
		ConversationSend{}.PreviewProcess(ctc)
		return ConversationSend{}
	} else {
		state.PreviewProcess(ctc)
		return FinishOrderCheck{}
	}
}

func (state FinishOrderCheck) PreviewProcess(ctc state.ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	b.Message("Вы уверены, что хотите завершить заказ?")
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Да", "", "positive")
	k.AddRow()
	k.AddTextButton("Нет", "", "negative")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state FinishOrderCheck) Name() string {
	return "FinishOrderCheck"
}

// ////////////////////////////////////////////////////////
type FinishOrder struct {
}

func (state FinishOrder) Process(ctc state.ChatContext, msg object.MessagesMessage) state.State {
	state.PreviewProcess(ctc)
	return FinishOrder{}
}

func (state FinishOrder) PreviewProcess(ctc state.ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	b.Message("Заказ завершён")
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}

func (state FinishOrder) Name() string {
	return "FinishOrder"
}
