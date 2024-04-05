package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/database"
	"xkcd/pkg/words"
	"xkcd/pkg/xkcd"
)

func main() {
	var n int
	flag.IntVar(&n, "n", -1, "max length commics")
	var useFlagO bool
	flag.BoolVar(&useFlagO, "o", false, "Use -o")
	flag.Parse()

	if err := initConfig(); err != nil {
		logrus.Fatalf("you have error %s", err.Error())
	}

	words := words.NewWordsStremming()

	cl := xkcd.NewClient(viper.GetString("source_url"), words)
	maxCommicsIndex, err := cl.GetLatestComicsNumber()
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())

	}

	if n > maxCommicsIndex {
		panic("Your n is bigger than max index of the commics.")
	}
	var numIter int
	if n > 0 {
		numIter = n
	} else {
		numIter = maxCommicsIndex
	}

	resoultMap := make(map[string]xkcd.ComicsInfo)
	for i := 1; i <= numIter; i++ {
		res, err := cl.GetComics(i)
		if err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
		for k, j := range *res {
			resoultMap[k] = j
		}
	}

	db := database.NewJsonDatabase(viper.GetString("db_file"))
	db.CreateEmptyDatabase()
	db.WriteAllOnDatabase(&resoultMap, useFlagO)

}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
