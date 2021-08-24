package start

import (
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func startHandler(b *gotgbot.Bot, ctx *ext.Context) error   {
	if !utils.IsUserSudo(ctx.EffectiveUser.Id){
		if !utils.IsUserOwner(ctx.EffectiveUser.Id){
			return nil
		}
	}
	_, err := ctx.EffectiveMessage.Reply(b, "Hello, this is a bot to manage licences and users of auth.gg", nil)
	if err != nil {
		return err
	}
	return nil
}

func LoadStartHandler(updater *ext.Updater)  {
	defer fmt.Println("Start module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", startHandler))
}