package main

import (
	"AuthtggO/db"
	"AuthtggO/logHelper"
	"AuthtggO/modules/deleteKey"
	"AuthtggO/modules/generateKey"
	"AuthtggO/modules/getInfo"
	"AuthtggO/modules/help"
	"AuthtggO/modules/registerUser"
	"AuthtggO/modules/resetadmin"
	"AuthtggO/modules/resethwid"
	"AuthtggO/modules/start"
	"AuthtggO/modules/stats"
	"AuthtggO/modules/test"
	"AuthtggO/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup


func RegisterAllHandlers(updater *ext.Updater)  {
	start.LoadStartHandler(updater)
	getInfo.LoadgetInfoHandler(updater)
	generateKey.LoadgenerateKeyHandler(updater)
	deleteKey.LoaddeleteKeyHandler(updater)
	test.LoadTestHandler(updater)
	resethwid.LoadResetUserHandler(updater)
	registerUser.LoadRegisterUserHandler(updater)
	stats.LoadStatsHandler(updater)
	help.LoadHelpHandler(updater)
	resetadmin.LoadResethwidAdminHandler(updater)

}
func periodicCheck()  {
	l := logHelper.GetLogger()
	err := db.UpdateAlllicenses()
	if err != nil {
		l.Info(err)
	}
	l.Info("Initialising update sequence")
	defer wg.Done()
	for {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Info("Checking for new licenses")
			err := db.UpdateAlllicenses()
			if err != nil {
				l.Info(err)
			}
			l.Info("Checking for new users")
			err = db.UpdateAllUsers()
			if err != nil {
				l.Info(err)
			}
		}()
		time.Sleep(time.Minute*3)
	}
}


func main()  {
	logHelper.InitLogHelper("log.log")
	l := logHelper.GetLogger()
	token := utils.GetBotToken()
	l.Info("Starting Bot.")
	l.Debug("token: ", token)
	err := db.InitDbClient()
	if err != nil {
		return 
	}
	wg.Add(1)
	go periodicCheck()
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