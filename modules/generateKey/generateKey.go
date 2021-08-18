package generateKey

import (
	"AuthtggO/helper/authGG"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strconv"
	"strings"
)

func generateKeyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserSudo(ctx.EffectiveUser.Id) {
		return nil
	}
	Quantity := 1
	overrideDefaultQuantity := false
	if len(strings.Split(ctx.EffectiveMessage.Text, " ")) == 4{
		overrideDefaultQuantity = true
	}
	if !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) >= 3 ) ||  !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) <= 4 ){
		_, err := ctx.EffectiveMessage.Reply(b, "I need you to send me 2 / 3 arguments", nil)
		if err != nil {
			return err
		}
		return nil
	}
	prefix := strings.Split(ctx.EffectiveMessage.Text, " ")[1]
	days, err :=strconv.Atoi(strings.Split(ctx.EffectiveMessage.Text, " ")[2])
	if err != nil {
		ctx.Message.Reply(b, "Specify correct option for number of days", &gotgbot.SendMessageOpts{})
		return err
	}
	if overrideDefaultQuantity{
		Quantity, err = strconv.Atoi(strings.Split(ctx.EffectiveMessage.Text, " ")[3])
		if err != nil {
			ctx.Message.Reply(b, "Specify correct option for number of quanity", &gotgbot.SendMessageOpts{})
			return err
		}
	}
	_, err = b.SendMessage(ctx.EffectiveChat.Id, authGG.CreateKey(Quantity, days, prefix),
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

func LoadgenerateKeyHandler(updater *ext.Updater)  {
	defer fmt.Println("generateKey module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("generate", generateKeyHandler))
}

