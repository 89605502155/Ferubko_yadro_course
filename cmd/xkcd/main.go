package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"xkcd/pkg/database"
	"xkcd/pkg/words"
	"xkcd/pkg/worker"
	"xkcd/pkg/xkcd"
)

func main() {
	n := 900
	var c bool
	flag.BoolVar(&c, "c", false, "Use -c")
	flag.Parse()
	if c {
		if err := initConfig(); err != nil {
			logrus.Fatalf("you have error %s", err.Error())
		}
	}
	db := database.NewJsonDatabase(viper.GetString("db_file"))

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// done := make(chan bool, 1)

	words := words.NewWordsStremming()

	cl := xkcd.NewClient(viper.GetString("source_url"), words)

	data := db.ReadDatabase()

	worker.WorkerPool(cl, n, viper.GetInt("parallel"), data)

	db.CreateEmptyDatabase()
	db.WriteAllOnDatabase(data, true)

}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
