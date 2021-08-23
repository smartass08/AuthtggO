package deleteKey

import (
	"AuthtggO/helper/authGG"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func deleteKeyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	var message string
	if !(len(strings.Split(ctx.EffectiveMessage.Text, "\n")) >= 1 ) {
		_, err := ctx.EffectiveMessage.Reply(b, "I need you to send me atleast 1 license key", nil)
		if err != nil {
			return err
		}
		return nil
	}
	keys := strings.Split(ctx.EffectiveMessage.Text, "\n")[1:]
	transmitChannel := make(chan string, len(keys))
	wg.Add(1)
	go func() {
		wg.Done()
		go authGG.DeleteKey(keys, transmitChannel)
	}()
	msg, err := b.SendMessage(ctx.EffectiveChat.Id,"Starting process for deleting the keys",
		&gotgbot.SendMessageOpts{
			ParseMode:             "MARKDOWN",
			ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
			DisableWebPagePreview: true},
	)
	if err != nil {
		return err
	}
	for v := range transmitChannel{
		message += v
		b.EditMessageText(message, &gotgbot.EditMessageTextOpts{MessageId: msg.MessageId, ParseMode: "MARKDOWN", ChatId: msg.Chat.Id})
	}
	b.EditMessageText(message + "\n\nAll keys processed", &gotgbot.EditMessageTextOpts{MessageId: msg.MessageId, ParseMode: "MARKDOWN", ChatId: msg.Chat.Id})
	return nil
}

func LoaddeleteKeyHandler(updater *ext.Updater)  {
	defer fmt.Println("deleteKey module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("deletekey", deleteKeyHandler))
}

