package state

import (
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func GetLink(ctc ChatContext, msg object.MessagesMessage, vkID int, isExec bool) (link api.MessagesGetInviteLinkResponse, err error) {
	order, err := ds.Unmarshal(msg.Payload)
	if err != nil {
		return api.MessagesGetInviteLinkResponse{}, fmt.Errorf("failed to unmarshal message payload, err: %s", err.Error())
	}

	chatID, err := ctc.Vk.MessagesCreateChat(api.Params{
		"title":    "Заказ номер №" + strconv.Itoa(order.OrderID) + "_" + strconv.Itoa(vkID) + "⚠",
		"user_ids": vkID,
		"group_id": msg.PeerID,
	})
	if err != nil {
		return api.MessagesGetInviteLinkResponse{}, fmt.Errorf("failed to create chat, err: %s", err.Error())
	}
	link, err = ctc.Vk.MessagesGetInviteLink(api.Params{
		"peer_id":  2000000000 + chatID,
		"reset":    0,
		"group_id": msg.PeerID,
	})
	_ = repository.AddChatID(ctc.Db, 2000000000+chatID, order.OrderID, isExec)
	if err != nil {
		return api.MessagesGetInviteLinkResponse{}, fmt.Errorf("failed to create link, err: %s", err.Error())
	}
	return link, nil
}

func SendLinkExecutor(ctc ChatContext, msg object.MessagesMessage) (err error) {

	order, err := ds.Unmarshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message payload, err: %s", err.Error())
	}
	executor, err := repository.GetExecutorByID(ctc.Db, order.ExecutorID)
	if err != nil {
		return fmt.Errorf("failed to get executor, err: %s", err.Error())
	}
	link, err := GetLink(ctc, msg, executor.VkID, true)
	if err != nil {
		return fmt.Errorf("failed to get link, err: %s", err.Error())
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Ссылка-приглашение в анонимную беседу: " + link.Link)
	b.PeerID(executor.VkID)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Println(err)
		return err
	}
	return nil
}

type ChoosingExecutor struct {
}

func (state ChoosingExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Сейчас будут создана анонимная беседа, позволяющая общаться через бота-посредника с исполнителем!")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
		StartState{}.PreviewProcess(ctc)
		return StartState{}
	}
	//создаем анонимную беседу с заказчиком
	order, err := ds.Unmarshal(msg.Payload)
	_ = repository.CreateConversations(ctc.Db, order.OrderID)
	link, err := GetLink(ctc, msg, ctc.User.VkID, false)
	if err != nil {
		log.Println(err)
		StartState{}.PreviewProcess(ctc)
		return StartState{}
	}
	b = params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Ссылка-приглашение в анонимную беседу: " + link.Link)
	b.PeerID(ctc.User.VkID)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	err = SendLinkExecutor(ctc, msg)
	if err != nil {
		b = params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Не удалось отправить приглашение в анонимную беседу для исполнителя, попробуйте выбрать другого исполнителя")
		b.PeerID(ctc.User.VkID)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
		}
	}
	StartState{}.PreviewProcess(ctc)
	return StartState{}
}

func (state ChoosingExecutor) PreviewProcess(ctc ChatContext) {

}
func (state ChoosingExecutor) Name() string {
	return "ChoosingExecutor"
}

// ///////////////////////////////////////////////////////
type ReselectingExecutor struct {
}

func (state ReselectingExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Нельзя выбрать исполнителя повторно, пожалуйста, создайте новый заказ")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	StartState{}.PreviewProcess(ctc)
	return StartState{}
}

func (state ReselectingExecutor) PreviewProcess(ctc ChatContext) {

}
func (state ReselectingExecutor) Name() string {
	return "ReselectingExecutor"
}

// ///////////////////////////////////////////
type ChoosingExecutorError struct {
}

func (state ChoosingExecutorError) Process(ctc ChatContext, msg object.MessagesMessage) State {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Что-то пошло не так, мы не смогли выбрать для Вас этого исполнителя")
	b.PeerID(ctc.User.VkID)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	StartState{}.PreviewProcess(ctc)
	return StartState{}
}

func (state ChoosingExecutorError) PreviewProcess(ctc ChatContext) {

}
func (state ChoosingExecutorError) Name() string {
	return "ChoosingExecutorError"
}
