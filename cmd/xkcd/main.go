package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/words"
	"xkcd/pkg/worker"
	"xkcd/pkg/xkcd"
)

func main() {
	n := 100
	var c bool
	var s string
	flag.BoolVar(&c, "c", false, "Use -c")
	flag.StringVar(&s, "s", "", "string")
	flag.Parse()
	if c {
		if err := initConfig(); err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
	}
	db := database.NewJsonDatabase(viper.GetString("db_file"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	words := words.NewWordsStremming()

	cl := xkcd.NewClient(viper.GetString("source_url"), words)

	data := db.Database.ReadDatabase()
	exitChan := make(chan bool, 1)
	isWriteChan := make(chan bool, 1)

	go func(db *database.JsonDatabase) {
		for {
			if <-exitChan {
				fmt.Println("Genrich")
				db.Database.CreateEmptyDatabase()
				db.Database.WriteAllOnDatabase(data, false)
				fmt.Println("go func")
				// stop()
				isWriteChan <- true
				return
			}
		}
	}(db)

	worker.WorkerPool(cl, n, viper.GetInt("parallel"), data, ctx, stop, exitChan, isWriteChan)

	fmt.Println("Egaspotamo")
	db.Database.CreateEmptyDatabase()
	db.Database.WriteAllOnDatabase(data, false)

	fmt.Println("after all")

	index := indexbase.NewJsonIndex(viper.GetString("index_file"))
	// index.CreateEmptyDatabase()
	indexes := index.IndexBase.ReadBase()

	index.IndexBase.BuildIndexFromDB(data, indexes)
	// fmt.Println(indexes)
	index.IndexBase.SaveIndexToFile(indexes)

	inputDataSFlag, err := words.Stremming.Normalization(s)
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())
	}
	fmt.Println("Robert")
	fmt.Println("input string ", inputDataSFlag)
	firstFind := db.FindInDB.Find(inputDataSFlag)
	fmt.Println("first find ", firstFind)
	secondFind := index.IndexFind.Find(inputDataSFlag)
	fmt.Println("second find ", secondFind)

}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
