package main

import (
	"AuthtggO/logHelper"
	"AuthtggO/modules/generateKey"
	"AuthtggO/modules/getInfo"
	"AuthtggO/modules/start"
	"AuthtggO/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"net/http"
)

func RegisterAllHandlers(updater *ext.Updater)  {
	start.LoadStartHandler(updater)
	getInfo.LoadgetInfoHandler(updater)
	generateKey.LoadgenerateKeyHandler(updater)
}

func main()  {
	logHelper.InitLogHelper("log.log")
	l := logHelper.GetLogger()
	token := utils.GetBotToken()
	l.Info("Starting Bot.")
	l.Debug("token: ", token)
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		l.Fatal(err)
	}
	updater := ext.NewUpdater(&ext.UpdaterOpts{})
	l.Info("Initializing bot updater")
	RegisterAllHandlers(&updater)
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		l.Fatalf("Error occurred at start of polling :  %s", err.Error())
		return
	}
	l.Info("Started Updater.")
	updater.Idle()
}