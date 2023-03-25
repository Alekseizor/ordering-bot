package state

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
)

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
