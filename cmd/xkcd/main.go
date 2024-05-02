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
	var i bool
	var s, c string
	flag.StringVar(&c, "c", "", "Use -c")
	flag.BoolVar(&i, "i", false, "Use -i")
	// flag.BoolVar(&u, "u", false, "update db and index")
	flag.StringVar(&s, "s", "", "string")
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
	service := service.NewService(db, index, n, cl, ctx, stop)
	defer stop()
	if s != "" {
		inputDataSFlag, err := words.Stremming.Normalization(s)
		if err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
		// fmt.Println("Robert")
		fmt.Println("input string ", inputDataSFlag)
		a1 := time.Now()
		firstFind := db.FindInDB.Find(inputDataSFlag, viper.GetInt("serch_limit"))
		a2 := time.Now()
		aDelta := a2.Sub(a1)
		fmt.Println("first find ", firstFind)
		if i {
			b1 := time.Now()
			secondFind := index.IndexFind.Find(inputDataSFlag, viper.GetInt("serch_limit"))
			b2 := time.Now()
			bDelta := b2.Sub(b1)
			fmt.Println("second find ", secondFind)
			fmt.Println("time ", aDelta, bDelta, aDelta/bDelta)
		}
	}

	handler := handler.NewHandler(service)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}

	}()
	logrus.Print("RNEB Started")

	<-ctx.Done()
	time.Sleep(5 * time.Second)
	logrus.Print("RNEB Shutting Down")

	if err := srv.ShutDown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}

func initConfig(c string) error {
	viper.AddConfigPath(c)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
