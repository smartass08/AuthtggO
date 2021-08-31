package registerUser

import (
	"AuthtggO/db"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strings"
)

func registerUserHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	var (
		message string
		newUser bool = true
		userMessage string = ctx.EffectiveMessage.Text
	)
	if len(strings.Split(ctx.EffectiveMessage.Text, " ")) == 1 {
		message = "This command allows me to register you on TG as valid customer.\n\n" +
					"Please send command in the following order :- \n`/register --username loader_USERNAME --key ENDUTY-XXX-XXX-KEY`\n" +
					"\nwhere `loader_USERNAME` is the username of the account you created in loader"

		_, err := ctx.EffectiveMessage.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "MARKDOWN"})
		if err != nil {
			return err
		}
		return nil
	}
	if !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) == 5) || !strings.Contains(userMessage, "--username") ||!strings.Contains(userMessage, "--key"){
		message = "Wrong format of message, Please press /register to check the correct format"
		_, err := ctx.EffectiveMessage.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "MARKDOWN"})
		if err != nil {
			return err
		}
		return nil
	}
	loaderUser := strings.Split(strings.Split(userMessage, "--username")[1], " ")[1]
	loaderKey := strings.Split(strings.Split(userMessage, "--key")[1], " ")[1]
	result, valid := db.GetOneLicense(loaderKey)
	if !valid{
		message = "No such license found, Please make sure key exists and its correct"
		newUser = false
	} else if strings.ToLower(result.UsedBy) != strings.ToLower(loaderUser){
			message = "No such user found, Please make sure user exists and its valid"
			newUser = false
	}
	if db.CheckTgUserExists(ctx.EffectiveUser.Id){
			message = "You have already registered, Please contact admin."
			newUser = false
	} else if db.IstgUserScammer(loaderKey, loaderUser){
			message = "The following key and user has already been registered, Please contact admin."
			newUser = false
		}
	if newUser{
		message = "Sucessfully verified!"
		err := db.AddtgUser(db.TGUser{
			TgID: ctx.EffectiveUser.Id,
			KeyUser:       []string{loaderUser},
			BindedKeys:    []string{loaderKey},
		})
		if err != nil {
			message = err.Error()
		}
	}
	_, err := b.SendMessage(ctx.EffectiveChat.Id,message, &gotgbot.SendMessageOpts{ParseMode: "MARKDOWN"})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}


func LoadRegisterUserHandler(updater *ext.Updater)  {
	defer fmt.Println("RegisterUser module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("register", registerUserHandler))
}