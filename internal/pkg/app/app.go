package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/Alekseizor/ordering-bot/internal/app/dsn"
	"github.com/Alekseizor/ordering-bot/internal/app/repository"
	"github.com/Alekseizor/ordering-bot/internal/app/state"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"time"

	"github.com/Alekseizor/ordering-bot/internal/anonymous_conversation"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/jmoiron/sqlx"
)

type App struct {
	ctx context.Context

	vk *api.VK
	lp *longpoll.LongPoll

	// db подключение к БД
	db *sqlx.DB
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.FromContext(ctx)
	vk := api.NewVK(cfg.VKToken)
	//получаем всю инфу про группу
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.WithError(err).Error("cant get groups by id")

		return nil, err
	}
	// БД
	db, err := sqlx.Connect("postgres", dsn.FromEnv())
	if err != nil {
		log.Println("nen", err)
		return nil, err
	}
	//starting long poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Println("error on request")
		log.Error(err)
	}
	app := &App{
		ctx: ctx,
		vk:  vk,
		lp:  lp,
		db:  db,
	}
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	var err error
	go func() error {
		if err = InitSysRoutes(ctx); err != nil {
			log.WithError(err).Error("can't InitSysRoute")
			return err
		}
		return nil
	}()

	var BotUser *ds.User
	var BotUsers []*ds.User
	// New message event
	a.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		//смотрим, новый ли пользователь
		query := "SELECT * FROM users WHERE vk_id=" + strconv.Itoa(obj.Message.PeerID)
		err := a.db.Select(&BotUsers, query)
		if err != nil {
			log.WithError(err).Error("cant set user")
			return
		}
		if obj.Message.PeerID > 2000000000 {
			strInState := map[string]state.State{
				(&(anonymous_conversation.ForwardMessage{})).Name():   &(anonymous_conversation.ForwardMessage{}),
				(&(anonymous_conversation.ConversationSend{})).Name(): &(anonymous_conversation.ConversationSend{}),
				(&(anonymous_conversation.FinishOrderCheck{})).Name(): &(anonymous_conversation.FinishOrderCheck{}),
				(&(anonymous_conversation.FinishOrder{})).Name():      &(anonymous_conversation.FinishOrder{}),
			}
			if obj.Message.Action.Type == "chat_invite_user_by_link" {
				BotUser = &ds.User{}
				BotUser.VkID = obj.Message.PeerID
				BotUser.State = "ForwardMessage"
				_, err := a.db.ExecContext(a.ctx, "INSERT INTO users VALUES ($1, $2)", BotUser.VkID, BotUser.State)
				if err != nil {
					log.WithError(err).Error("cant set user")
					return
				}
			} else {
				BotUser = BotUsers[0]
			}
			//a.vk.MessagesGetConversationsByID()
			step := strInState[BotUser.State]
			ctc := state.ChatContext{
				User: BotUser,
				Vk:   a.vk,
				Db:   a.db,
				Ctx:  &ctx,
			}
			nextStep := step.Process(ctc, obj.Message)
			BotUser.State = nextStep.Name()
			_, err = a.db.ExecContext(a.ctx, "UPDATE users SET State = $1 WHERE vk_id = $2", BotUser.State, BotUser.VkID)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
			return
		}
		//if the user writes for the first time, add to the database
		if len(BotUsers) == 0 {
			BotUser = &ds.User{}
			BotUser.VkID = obj.Message.PeerID
			BotUser.State = "StartState"
			_, err := a.db.ExecContext(a.ctx, "INSERT INTO users VALUES ($1, $2)", BotUser.VkID, BotUser.State)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
		} else {
			BotUser = BotUsers[0]
		}
		strInState := map[string]state.State{
			(&(state.StartState{})).Name():                   &(state.StartState{}),
			(&(state.OrderType{})).Name():                    &(state.OrderType{}),
			(&(state.OrderState{})).Name():                   &(state.OrderState{}),
			(&(state.ChoiceDiscipline{})).Name():             &(state.ChoiceDiscipline{}),
			(&(state.ChoiceDate{})).Name():                   &(state.ChoiceDate{}),
			(&(state.ChoiceTime{})).Name():                   &(state.ChoiceTime{}),
			(&(state.ConfirmationOrder{})).Name():            &(state.ConfirmationOrder{}),
			(&(state.CommentOrder{})).Name():                 &(state.CommentOrder{}),
			(&(state.TaskOrder{})).Name():                    &(state.TaskOrder{}),
			(&(state.ConfirmExecutor{})).Name():              &(state.ConfirmExecutor{}),
			(&(state.OrderCompleted{})).Name():               &(state.OrderCompleted{}),
			(&(state.OrderCancel{})).Name():                  &(state.OrderCancel{}),
			(&(state.OrderChange{})).Name():                  &(state.OrderChange{}),
			(&(state.EditType{})).Name():                     &(state.EditType{}),
			(&(state.EditDiscipline{})).Name():               &(state.EditDiscipline{}),
			(&(state.EditDate{})).Name():                     &(state.EditDate{}),
			(&(state.EditTime{})).Name():                     &(state.EditTime{}),
			(&(state.EditTaskOrder{})).Name():                &(state.EditTaskOrder{}),
			(&(state.EditCommentOrder{})).Name():             &(state.EditCommentOrder{}),
			(&(state.BecomeExecutor{})).Name():               &(state.BecomeExecutor{}),
			(&(state.ExecHistoryOrders{})).Name():            &(state.ExecHistoryOrders{}),
			(&(state.WriteAdmin{})).Name():                   &(state.WriteAdmin{}),
			(&(state.CabinetAdmin{})).Name():                 &(state.CabinetAdmin{}),
			(&(state.UnloadTable{})).Name():                  &(state.UnloadTable{}),
			(&(state.UnloadTableExec{})).Name():              &(state.UnloadTableExec{}),
			(&(state.AddExecutor{})).Name():                  &(state.AddExecutor{}),
			(&(state.AddExecID{})).Name():                    &(state.AddExecID{}),
			(&(state.AddExecDisciplines{})).Name():           &(state.AddExecDisciplines{}),
			(&(state.ManageExecutors{})).Name():              &(state.ManageExecutors{}),
			(&(state.DeleteExecutorID{})).Name():             &(state.DeleteExecutorID{}),
			(&(state.DeleteExecutor{})).Name():               &(state.DeleteExecutor{}),
			(&(state.ChangeExecutorsDisciplinesID{})).Name(): &(state.ChangeExecutorsDisciplinesID{}),
			(&(state.ChangeExecutorsDisciplines{})).Name():   &(state.ChangeExecutorsDisciplines{}),
			(&(state.ChangeExecutorsCommissionID{})).Name():  &(state.ChangeExecutorsCommissionID{}),
			(&(state.ChangeExecutorsCommission{})).Name():    &(state.ChangeExecutorsCommission{}),
			(&(state.ChangeRequisites{})).Name():             &(state.ChangeRequisites{}),
			(&(state.Newsletter{})).Name():                   &(state.Newsletter{}),
			(&(state.NewsletterConfirmation{})).Name():       &(state.NewsletterConfirmation{}),
			(&(state.NewsletterPeerID{})).Name():             &(state.NewsletterPeerID{}),
			(&(state.NewsletterSend{})).Name():               &(state.NewsletterSend{}),
			(&(state.SelectionUnload{})).Name():              &(state.SelectionUnload{}),
			(&(state.SelectionDateUnload{})).Name():          &(state.SelectionDateUnload{}),
			(&(state.InputFirstDateUnload{})).Name():         &(state.InputFirstDateUnload{}),
			(&(state.InputSecondDateUnload{})).Name():        &(state.InputSecondDateUnload{}),
			(&(state.SelectionAllOrPersonalUnload{})).Name(): &(state.SelectionAllOrPersonalUnload{}),
			(&(state.PersonalUnload{})).Name():               &(state.PersonalUnload{}),
			(&(state.SelectionDateClear{})).Name():           &(state.SelectionDateClear{}),
			(&(state.InputFirstDateClear{})).Name():          &(state.InputFirstDateClear{}),
			(&(state.InputSecondDateClear{})).Name():         &(state.InputSecondDateClear{}),
			(&(state.ChangeRequisiteExecutor{})).Name():      &(state.ChangeRequisiteExecutor{}),
			(&(state.ConfirmationExecutor{})).Name():         &(state.ConfirmationExecutor{}),
			(&(state.ChoosingExecutor{})).Name():             &(state.ChoosingExecutor{}),
			(&(state.ReselectingExecutor{})).Name():          &(state.ReselectingExecutor{}),
			(&(state.ChoosingExecutorError{})).Name():        &(state.ChoosingExecutorError{}),
			(&(state.MyOrderState{})).Name():                 &(state.MyOrderState{}),
		}
		ctc := state.ChatContext{
			User: BotUser,
			Vk:   a.vk,
			Db:   a.db,
			Ctx:  &ctx,
		}
		//cfg := config.FromContext(*ctc.Ctx).Bot
		if obj.Message.Payload != "" {
			if obj.Message.Text == "Принять" {
				orderNumber, err := strconv.Atoi(obj.Message.Payload)
				if err != nil {
					log.WithError(err).Error("cant set user")
					return
				}
				repository.WriteOffer(ctc.Db, orderNumber, ctc.User.VkID)
				BotUser.State = "ConfirmationExecutor"
			} else {
				log.Println(obj.Message.Payload)
				execOrder, err := ds.Unmarshal(obj.Message.Payload)
				if err != nil {
					log.WithError(err).Error("couldn't parse payload")
					return
				}
				err = repository.AddingExecutor(ctc.Db, execOrder)
				if err == nil {
					BotUser.State = "ChoosingExecutor"
				} else if err.Error() == "the executor has already been selected" {
					BotUser.State = "ReselectingExecutor"
				} else {
					BotUser.State = "ChoosingExecutorError"
				}
			}
		}
		step := strInState[BotUser.State]
		nextStep := step.Process(ctc, obj.Message)
		BotUser.State = nextStep.Name()
		_, err = a.db.ExecContext(a.ctx, "UPDATE users SET State = $1 WHERE vk_id = $2", BotUser.State, BotUser.VkID)
		if err != nil {
			log.WithError(err).Error("cant set user")
			return
		}
	})
	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := a.lp.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

const (
	sysHTTPDefaultTimeout = 5 * time.Minute
)

func InitSysRoutes(ctx context.Context) error {

	mux := http.NewServeMux()
	{
		mux.HandleFunc("/ready", ReadyHandler)
		mux.HandleFunc("/live", LiveHandler)
	}

	port := "8080"

	s := &http.Server{
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: sysHTTPDefaultTimeout,
		ReadTimeout:  sysHTTPDefaultTimeout,
		IdleTimeout:  sysHTTPDefaultTimeout,
		Handler:      mux,
	}
	err := s.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		fmt.Println(err)
	}
	return err
}

func ReadyHandler(w http.ResponseWriter, _ *http.Request) {
	httpStatus := http.StatusOK
	w.WriteHeader(httpStatus)
	enc := json.NewEncoder(w)
	_ = enc.Encode(map[string]bool{
		"ready": true,
	})
}

func LiveHandler(w http.ResponseWriter, _ *http.Request) {
	httpStatus := http.StatusOK
	w.WriteHeader(httpStatus)
	enc := json.NewEncoder(w)
	_ = enc.Encode(map[string]bool{
		"live": true,
	})
}
