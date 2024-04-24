package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/database"
	"xkcd/pkg/indexbase"
	"xkcd/pkg/words"
)

func main() {
	// n := 40
	var c, i bool
	var s string
	flag.BoolVar(&c, "c", false, "Use -c")
	flag.BoolVar(&i, "i", false, "Use -i")
	flag.StringVar(&s, "s", "", "string")
	flag.Parse()
	if c {
		if err := initConfig(); err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
	}
	db := database.NewJsonDatabase(viper.GetString("db_file"))

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()

	words := words.NewWordsStremming()

	// cl := xkcd.NewClient(viper.GetString("source_url"), words)

	// data := db.Database.ReadDatabase()
	// exitChan := make(chan bool, 1)
	// isWriteChan := make(chan bool, 1)

	// go func(db *database.JsonDatabase) {
	// 	for {
	// 		if <-exitChan {
	// 			// fmt.Println("Genrich")
	// 			db.Database.CreateEmptyDatabase()
	// 			db.Database.WriteAllOnDatabase(data, false)
	// 			// fmt.Println("go func")
	// 			// stop()
	// 			isWriteChan <- true
	// 			return
	// 		}
	// 	}
	// }(db)

	// worker.WorkerPool(cl, n, viper.GetInt("parallel"), data, ctx, stop, exitChan, isWriteChan)

	// fmt.Println("Egaspotamo")
	// db.Database.CreateEmptyDatabase()
	// db.Database.WriteAllOnDatabase(data, false)

	// fmt.Println("after all")

	index := indexbase.NewJsonIndex(viper.GetString("index_file"))
	// // index.CreateEmptyDatabase()
	// indexes := index.IndexBase.ReadBase()

	// index.IndexBase.BuildIndexFromDB(data, indexes)
	// // fmt.Println(indexes)
	// index.IndexBase.SaveIndexToFile(indexes)

	if s != "" {
		inputDataSFlag, err := words.Stremming.Normalization(s)
		if err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
		fmt.Println("Robert")
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

}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
