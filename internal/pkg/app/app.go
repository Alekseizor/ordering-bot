package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alekseizor/ordering-bot/internal/app/config"
	"github.com/Alekseizor/ordering-bot/internal/app/ds"
	"github.com/Alekseizor/ordering-bot/internal/app/dsn"
	"github.com/Alekseizor/ordering-bot/internal/app/state"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
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
		log.Println("nen")
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
		//if the user writes for the first time, add to the database
		if BotUsers == nil {
			BotUser = &ds.User{}
			BotUser.VkID = obj.Message.PeerID
			BotUser.State = "StartState"
			_, err := a.db.ExecContext(a.ctx, "INSERT INTO users VALUES ($1, $2)", BotUser.VkID, BotUser.State)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
			b := params.NewMessagesSendBuilder()
			b.RandomID(0)
			b.Message("Привет! Добро пожаловать в главное меню бота. Пришли мне номер нужной команды или воспользуйся кнопками")
			b.PeerID(BotUser.VkID)
			_, err = a.vk.MessagesSend(b.Params)
			if err != nil {
				log.Println("Failed to get record")
				log.Error(err)
			}
			b.Message("1. Сделать заказ\n2. Связаться с исполнителем \n3. Оставить отзыв\n4. Сделать заказ через посредника\n5. Стать исполнителем\n6. Мои заказы")
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
			_, err = a.vk.MessagesSend(b.Params)
			if err != nil {
				log.Println("Failed to get record")
				log.Error(err)
			}
			return
		} else {
			BotUser = BotUsers[0]
		}
		strInState := map[string]state.State{
			(&(state.StartState{})).Name():         &(state.StartState{}),
			(&(state.OrderType{})).Name():          &(state.OrderType{}),
			(&(state.OrderState{})).Name():         &(state.OrderState{}),
			(&(state.ChoiceDiscipline{})).Name():   &(state.ChoiceDiscipline{}),
			(&(state.ChoiceDate{})).Name():         &(state.ChoiceDate{}),
			(&(state.ChoiceTime{})).Name():         &(state.ChoiceTime{}),
			(&(state.ConfirmationOrder{})).Name():  &(state.ConfirmationOrder{}),
			(&(state.CommentOrder{})).Name():       &(state.CommentOrder{}),
			(&(state.TaskOrder{})).Name():          &(state.TaskOrder{}),
			(&(state.OrderCompleted{})).Name():     &(state.OrderCompleted{}),
			(&(state.OrderCancel{})).Name():        &(state.OrderCancel{}),
			(&(state.OrderChange{})).Name():        &(state.OrderChange{}),
			(&(state.EditType{})).Name():           &(state.EditType{}),
			(&(state.EditDiscipline{})).Name():     &(state.EditDiscipline{}),
			(&(state.EditDate{})).Name():           &(state.EditDate{}),
			(&(state.EditTime{})).Name():           &(state.EditTime{}),
			(&(state.EditTaskOrder{})).Name():      &(state.EditTaskOrder{}),
			(&(state.EditCommentOrder{})).Name():   &(state.EditCommentOrder{}),
			(&(state.BecomeExecutor{})).Name():     &(state.BecomeExecutor{}),
			(&(state.ExecHistoryOrders{})).Name():  &(state.ExecHistoryOrders{}),
			(&(state.WriteAdmin{})).Name():         &(state.WriteAdmin{}),
			(&(state.CabinetAdmin{})).Name():       &(state.CabinetAdmin{}),
			(&(state.UnloadTable{})).Name():        &(state.UnloadTable{}),
			(&(state.AddExecutor{})).Name():        &(state.AddExecutor{}),
			(&(state.AddExecID{})).Name():          &(state.AddExecID{}),
			(&(state.AddExecDisciplines{})).Name(): &(state.AddExecDisciplines{}),
			(&(state.ManageExecutors{})).Name():    &(state.ManageExecutors{}),
			(&(state.DeleteExecutorID{})).Name():   &(state.DeleteExecutorID{}),
			(&(state.DeleteExecutor{})).Name():     &(state.DeleteExecutor{}),
		}
		ctc := state.ChatContext{
			User: BotUser,
			Vk:   a.vk,
			Db:   a.db,
			Ctx:  &ctx,
		}
		//cfg := config.FromContext(*ctc.Ctx).Bot
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
