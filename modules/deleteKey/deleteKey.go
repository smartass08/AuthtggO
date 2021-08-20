package deleteKey

import (
	"AuthtggO/helper/authGG"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strings"
)

func deleteKeyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) >= 1 ) {
		_, err := ctx.EffectiveMessage.Reply(b, "I need you to send me atleast 1 license key", nil)
		if err != nil {
			return err
		}
		return nil
	}
	keys := strings.Split(ctx.EffectiveMessage.Text, " ")[1:]

	_, err := b.SendMessage(ctx.EffectiveChat.Id, authGG.DeleteKey(keys),
		&gotgbot.SendMessageOpts{
			ParseMode:             "MARKDOWN",
			ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
			DisableWebPagePreview: true},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func LoaddeleteKeyHandler(updater *ext.Updater)  {
	defer fmt.Println("deleteKey module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("deletekey", deleteKeyHandler))
}

