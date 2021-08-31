package help

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func LoadHelpHandler(updater *ext.Updater)  {
	defer fmt.Println("Help module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("help", helpHandler))
}

func helpHandler(b *gotgbot.Bot, ctx *ext.Context) error  {
	if ctx.EffectiveChat.Id != ctx.EffectiveUser.Id{
		return nil
	}
	message := fmt.Sprintf("<b>Help Section of this management bot</b>\n\n"+
		"Available commands are as follows: \n\n"+
		"/start : <i>The all boring start command.</i>\n\n"+
		"/register : <i>Registers you as valid user, Tap register to know more</i>\n\n"+
		"/stats : <i>The command gives you info about you</i>\n\n" +
		"/resetme : Use this command to reset yourself without the need of admin, Need to send one loader username\n\n")
	_, err := ctx.EffectiveMessage.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "HTMl"})
	if err != nil {
		return err
	}

	return nil

}
