package resethwid

import (
	"AuthtggO/helper/authGG"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strings"
)

func ResetUserHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserSudo(ctx.EffectiveUser.Id) {
		return nil
	}
	if !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) == 2) {
		_, err := ctx.EffectiveMessage.Reply(b,"Please send me 1 username to whom you wanna reset", nil)
		if err != nil {
			return err
		}
		return nil
	}

	user := strings.Split(ctx.EffectiveMessage.Text, " ")[1]
	result, err := authGG.ResetHwid(user)
	if err != nil {
		return err
	}
	_, err = b.SendMessage(ctx.EffectiveChat.Id, result,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
			DisableWebPagePreview: true},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func LoadResetUserHandler(updater *ext.Updater)  {
	defer fmt.Println("generateKey module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("resethw", ResetUserHandler))
}

