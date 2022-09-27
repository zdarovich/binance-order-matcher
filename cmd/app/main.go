package main

import (
	"binance-order-matcher/config"
	"binance-order-matcher/internal/model"
	"binance-order-matcher/internal/service"
	"binance-order-matcher/internal/service/api"
	"binance-order-matcher/internal/service/repo"
	"binance-order-matcher/pkg/sqlite"
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	conf, err := config.Init()
	if err != nil {
		log.Panic("while getting config", err)
		os.Exit(1)
	}
	err = sqlite.Create(conf.Database.URL)
	db, err := sqlite.NewConnection(conf.Database.URL)
	if err != nil {
		log.Panic("while connecting to database. ", err)
		os.Exit(1)
	}
	m, err := migrate.New("file://migrations", conf.Database.URL)
	if err != nil {
		log.Panic("while connecting to database. ", err)
		os.Exit(1)
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("migrate: no change")
		} else {
			log.Panic("while migrating database. ", err)
			os.Exit(1)
		}
	}
	orderRepo := repo.NewOrderRepo(db)
	bookRepo := repo.NewBookRepo(db)
	subscriptionAPI := api.New(conf.OrderMatcherService.URL)
	orderService := service.NewOrderMatcher(bookRepo, subscriptionAPI)

	id, _ := uuid.NewUUID()
	order := model.CreateOrder(id.String(), conf.Order.Type, conf.Order.Symbol, conf.Order.Size, conf.Order.Price)
	_, err = orderRepo.Save(order)
	if err != nil {
		log.Error("while trying to save the order: ", err)
		os.Exit(1)
	}
	isMatch, err := orderService.Match(ctx, 60*time.Second, order)
	if err != nil {
		log.Error("while trying to match the order: ", err)
		os.Exit(1)
	}

	if isMatch {
		log.Info("limit order was executed")
		books, err := bookRepo.GetByOrderId(order.Id)
		if err != nil {
			log.Error("while trying to match the order: ", err)
			os.Exit(1)
		}
		log.Infof("order: %+v", order)
		for _, el := range books {
			log.Infof("book: %+v", el)
		}
	}
}
