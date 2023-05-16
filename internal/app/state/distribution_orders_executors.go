package state

import (
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func DistributionOrderExecutors(ctc ChatContext) error {
	OrderID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		return err
	}
	Order, err := repository.GetOrder(ctc.Db, OrderID)
	if err != nil {
		return err
	}
	executorsDiscipline, err := repository.ExecutorsDiscipline(ctc.Db, Order.DisciplineID)
	if err != nil {
		return err
	}
	for _, executor := range executorsDiscipline {
		//todo: убрать комментарий ниже
		//if executor == ctc.User.VkID {
		//	continue
		//}
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		output, err := repository.GetCompleteOrder(ctc.Db, ctc.User.VkID)
		if err != nil {
			log.Println("Failed to get orders output")
			log.Error(err)
			return err
		}
		outputRune := []rune(output)
		outputRune = outputRune[17:]
		b.Message("Новый заказ:\n" + string(outputRune))
		attachment, err := repository.GetAttachments(ctc.Vk, ctc.Db, ctc.User.VkID)
		if err != nil {
			log.Println("Failed to get attachments")
			log.Error(err)
			return err
		}
		b.Attachment(attachment)
		b.PeerID(executor)
		k := &object.MessagesKeyboard{}
		k.Inline = true
		k.AddRow()
		k.AddTextButton("Принять", OrderID, "secondary")
		b.Keyboard(k)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Println("Failed to get record")
			log.Error(err)
			return err
		}
	}
	return nil
}

func DirectDistribution(ctc ChatContext, execID int) error {
	OrderID, err := repository.GetIDOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		return err
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	output, err := repository.GetCompleteOrder(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("Failed to get orders output")
		log.Error(err)
		return err
	}
	outputRune := []rune(output)
	outputRune = outputRune[17:]
	b.Message("Новый заказ:\n" + string(outputRune))
	attachment, err := repository.GetAttachments(ctc.Vk, ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("Failed to get attachments")
		log.Error(err)
		return err
	}
	executor, err := repository.GetExecutorByID(ctc.Db, execID)
	if err != nil {
		log.Println("Failed to get executor")
		log.Error(err)
		return err
	}
	b.Attachment(attachment)
	b.PeerID(executor.VkID)
	k := &object.MessagesKeyboard{}
	k.Inline = true
	k.AddRow()
	k.AddTextButton("Принять", OrderID, "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
		return err
	}

	return nil
}

func SendOffer(ctc ChatContext) error {
	offerID, err := repository.GetOfferID(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("couldn't find an offerID")
		return err
	}
	offer, err := repository.GetOffer(ctc.Db, offerID)
	if err != nil {
		log.Println("couldn't find an offer")
		return err
	}
	order, err := repository.GetOrder(ctc.Db, offer.OrderID)
	if err != nil {
		log.Println("couldn't find the order")
		return err
	}
	executor, err := repository.GetExecutor(ctc.Db, ctc.User.VkID)
	if err != nil {
		log.Println("couldn't find the executorID")
		return err
	}
	output := "По Вашему заказу №" + strconv.Itoa(order.Id) + " есть новое предложение:\nИсполнитель: №" + strconv.Itoa(executor.Id) + "\nЦена: " + strconv.Itoa(offer.Price) + "\nРейтинг исполнителя: "
	if executor.AmountRating == 0 {
		output += "Исполнитель ещё не получал оценок"
	} else {
		output += fmt.Sprint(executor.Rating)
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message(output)
	b.PeerID(int(order.CustomerVkID))
	k := &object.MessagesKeyboard{}
	k.Inline = true
	k.AddRow()
	payload := ds.ExecutorOrder{
		ExecutorID: executor.Id,
		OrderID:    order.Id,
		Price:      offer.Price,
	}
	k.AddTextButton("Выбрать этого исполнителя", payload, "secondary")
	b.Keyboard(k)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	return nil
}

type ConfirmationExecutor struct {
}

func (state ConfirmationExecutor) Process(ctc ChatContext, msg object.MessagesMessage) State {
	messageText := msg.Text
	if messageText == "Назад" {
		err := repository.DeleteOffer(ctc.Db, ctc.User.VkID)
		if err != nil {
			log.Println("Offer doesn't delete")
		}
		StartState{}.PreviewProcess(ctc)
		return &StartState{}
	}
	price, err := strconv.Atoi(messageText)
	if err != nil {
		log.Println("couldn't convert to a number")
		ConfirmationExecutor{}.PreviewProcess(ctc)
		return &ConfirmationExecutor{}
	}
	err = repository.WritePriceOffer(ctc.Db, ctc.User.VkID, price)
	if err != nil {
		log.Println("Price doesn't write on offer")
		ConfirmationExecutor{}.PreviewProcess(ctc)
		return &ConfirmationExecutor{}
	}
	err = SendOffer(ctc)
	if err != nil {
		log.Println("Couldn't send the offer")
		ConfirmationExecutor{}.PreviewProcess(ctc)
		return &ConfirmationExecutor{}
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Спасибо за отклик! Ваше предложение отправлено заказчику!")
	b.PeerID(ctc.User.VkID)
	_, err = ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
	StartState{}.PreviewProcess(ctc)
	return StartState{}
}

func (state ConfirmationExecutor) PreviewProcess(ctc ChatContext) {
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.Message("Введите цену заказа")
	b.PeerID(ctc.User.VkID)
	k := &object.MessagesKeyboard{}
	k.AddRow()
	k.AddTextButton("Назад", "", "secondary")
	b.Keyboard(k)
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Println("Failed to get record")
		log.Error(err)
	}
}
func (state ConfirmationExecutor) Name() string {
	return "ConfirmationExecutor"
}
