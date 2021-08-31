package stats

import (
	"AuthtggO/db"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadStatsHandler(updater *ext.Updater)  {
	defer fmt.Println("stats module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("stats", statsHandler))
}

func statsHandler(b *gotgbot.Bot, ctx *ext.Context) error  {
	var (
		message string
	)
	tempUser, check := db.GetOneTgUser(ctx.EffectiveUser.Id)
	if check{
		TempData := make(map[db.User]db.Licence,0)
		for _, v := range  tempUser.KeyUser{
			user, _ := db.GetOneUser(v)
			lic, _ :=  db.GetOneLicense(db.GetallKeysOfUser(user.UserName)[0])
			TempData[user] = lic
		}
		for i, v := range TempData{
			resetCount, err := db.FetchResetCount(v.Key)
			if err != nil {
				message = err.Error()
			} else {
				message += fmt.Sprintf(
					"Loader Username : `%v`\n" +
						"Key : `%v`\n" +
						"Expiry of Key : `%v`\n" +
						"Resets Done this month : `%v`\n\n", i.UserName, v.Key, i.ExpiryDate, resetCount)
			}
		}
	} else {
		message = "You have not registered with me yet, Press /register for more information"
	}
	_, err := ctx.EffectiveMessage.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "MARKDOWN"})
	if err != nil {
		return err
	}

	return nil

}