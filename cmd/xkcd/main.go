package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	server "xkcd"
	"xkcd/pkg/handler"
	"xkcd/pkg/personal_limiter"
	"xkcd/pkg/rate_limiter"
	"xkcd/pkg/repository"
	"xkcd/pkg/service"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

func main() {
	n := 1408
	var c string
	var p string
	flag.StringVar(&c, "c", "", "Use -c")
	flag.StringVar(&p, "p", "", "Use -p")
	flag.Parse()
	if err := initConfig(c); err != nil {
		logrus.Debug(c)
		logrus.Fatalf("you have error %s", err.Error())
	}
	if p == "" {
		p = viper.GetString("port")
	}
	sqlite, err := repository.NewSQLiteDB(repository.Config{
		DBName: "./xkcd.db",
	})
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
	}
	words := words.NewWordsStremming()

	cl := xkcd.NewClient(viper.GetString("source_url"), words)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	srv := new(server.Server)
	repo := repository.NewRepository(sqlite)
	service := service.NewService(n, cl, ctx, stop, repo, words, viper.GetInt("serch_limit"))
	defer stop()

	rate_limiter := rate_limiter.NewSlidingLogLimiter(10, 1)
	personal_limiter := personal_limiter.NewPersonalLimiter(ctx, viper.GetInt("person_limit"), time.Minute)

	handler := handler.NewHandler(service, rate_limiter, personal_limiter)

	go func() {
		for {
			currentTime := time.Now()
			targetTime := time.Date(currentTime.Year(), currentTime.Month(),
				currentTime.Day()+1, 2, 18, 0, 0, time.UTC)
			duration := targetTime.Sub(currentTime)
			logrus.Print(duration)
			timer := time.NewTimer(duration)
			<-timer.C
			service.Update()
			logrus.Debug("Update")
		}
	}()

	go func() {
		if err := srv.Run(p, handler.InitRoutes()); err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
	}()
	logrus.Print("xkcd Started")

	<-ctx.Done()
	time.Sleep(5 * time.Second)
	logrus.Print("xkcd Shutting Down")

	if err := srv.ShutDown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig(c string) error {
	viper.AddConfigPath(c)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
