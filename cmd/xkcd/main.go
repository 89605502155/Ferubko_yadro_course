package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/xkcd"
)

func main() {
	var n int
	flag.IntVar(&n, "n", 0, "max length commics")
	flag.Parse()

	if err := initConfig(); err != nil {
		logrus.Fatalf("you have error %s", err.Error())
	}

	cl := xkcd.NewClient(viper.GetString("source_url"))
	maxCommicsIndex, err := cl.GetLatestComicNumber()
	if err != nil {
		logrus.Fatalf("you have error %s", err.Error())

	}
	fmt.Println("b", maxCommicsIndex)

	if n > maxCommicsIndex {
		panic("Your n is bigger than max index of the commics.")
	}

}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
