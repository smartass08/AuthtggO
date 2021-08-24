package generateKey

import (
	"AuthtggO/helper/InfoPackage"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func generateKeyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserSudo(ctx.EffectiveUser.Id){
		if !utils.IsUserOwner(ctx.EffectiveUser.Id){
			return nil
		}
	}
	var message string
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
	wg.Add(1)
	channelKeys := make(chan string, Quantity)
 	go InfoPackage.CreateKey(Quantity, days, prefix, channelKeys)
	if err != nil {
		message = err.Error()
	}
	msg, err := b.SendMessage(ctx.EffectiveChat.Id, "Starting process to create keys\nThis may take some time as API is slow as your grandma",
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
			DisableWebPagePreview: true},
	)
	if err != nil {
		return err
	}

	for i := range channelKeys{
		message +=  fmt.Sprintf("<code>%v</code>\n", i)
		_, err := b.EditMessageText(message, &gotgbot.EditMessageTextOpts{MessageId: msg.MessageId, ChatId: msg.Chat.Id, ParseMode: "HTML"})
		if err != nil {
			return err
		}
	}
	_, err = b.EditMessageText(message+"\n\nAll keys processed", &gotgbot.EditMessageTextOpts{MessageId: msg.MessageId, ChatId: msg.Chat.Id, ParseMode: "HTML"})
	if err != nil {
		return err
	}
	return nil
}

func LoadgenerateKeyHandler(updater *ext.Updater)  {
	defer fmt.Println("generateKey module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("generate", generateKeyHandler))
}

