package getInfo

import (
	"AuthtggO/helper/InfoPackage"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strings"
)

func checkMods(b *gotgbot.Bot, person int64) bool{
	allAdmins, err := b.GetChatAdministrators(utils.GetModChatId())
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, v := range allAdmins{
		if person == v.GetUser().Id{
			return true
		}
	}
	return false
}

func getInfoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserSudo(ctx.EffectiveUser.Id){
		if !utils.IsUserOwner(ctx.EffectiveUser.Id){
			if !checkMods(b, ctx.EffectiveUser.Id){
				return nil
			}

		}
	}

	if !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) == 3){
		_, err := ctx.EffectiveMessage.Reply(b, "I need you to send me 2 arguments", nil)
		if err != nil {
			return err
		}
		return nil
	}
	firstArgument := strings.Split(ctx.EffectiveMessage.Text, " ")[1]
	secondArgument := strings.Split(ctx.EffectiveMessage.Text, " ")[2]
	_, err := b.SendMessage(ctx.EffectiveChat.Id,InfoPackage.GetInfo(firstArgument, secondArgument),
		&gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
		DisableWebPagePreview: true},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func LoadgetInfoHandler(updater *ext.Updater)  {
	defer fmt.Println("getInfo module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("getinfo", getInfoHandler))
}
