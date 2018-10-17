package main

import (
    "os"
    "fmt"
	"github.com/spf13/viper"
	sv "github.com/jingwu15/composer_pack/service"
)

func main() {
    sv.Run()
}

func init() {
	viper.SetConfigFile("./composer_pack.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//go logchan.LogWrite()
	//log.SetFormatter(&log.JSONFormatter{})
	////log.SetOutput(ioutil.Discard)
	//log.SetLevel(log.DebugLevel)
	//config := map[string]string{
	//	"error":      viper.GetString("log_error"),
	//	"info":       viper.GetString("log_info"),
	//	"writeDelay": "1",
	//	"cutType":    "day",
	//}
	//logChanHook := logchan.NewLogChanHook(config)
	//log.AddHook(&logChanHook)
}

