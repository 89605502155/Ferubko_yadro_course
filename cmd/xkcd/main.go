package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	server "xkcd"
	"xkcd/pkg/database"
	"xkcd/pkg/handler"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/service"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

func main() {
	n := 1408
	var c string
	flag.StringVar(&c, "c", "", "Use -c")
	// flag.BoolVar(&i, "i", false, "Use -i")
	// flag.BoolVar(&u, "u", false, "update db and index")
	// flag.StringVar(&s, "s", "", "string")
	flag.Parse()
	if err := initConfig(c); err != nil {
		fmt.Println(c)
		logrus.Fatalf("you have error %s", err.Error())
	}
	db := database.NewJsonDatabase(viper.GetString("db_file"))

	words := words.NewWordsStremming()

	cl := xkcd.NewClient(viper.GetString("source_url"), words)
	index := indexbase.NewJsonIndex(viper.GetString("index_file"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	service := service.NewService(db, index, n, cl, ctx, stop, words, viper.GetInt("serch_limit"))
	defer stop()

	handler := handler.NewHandler(service)
	srv := new(server.Server)
	go func() {
		fmt.Println("Rockrya")
		targetTime := time.Date(2024, time.May, 2, 2, 18, 0, 0, time.UTC)

		// Получаем текущее время
		currentTime := time.Now()
		duration := targetTime.Sub(currentTime)
		fmt.Println(duration)
		timer := time.NewTimer(duration)
		<-timer.C
		fmt.Println("Valonia")
		service.Update()
		fmt.Println("Update")
	}()
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
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
