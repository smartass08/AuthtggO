package test

import (
	"AuthtggO/helper/db"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func testHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserOwner(ctx.EffectiveUser.Id) {
		return nil
	}
	fmt.Println()
	/*_, err := b.SendMessage(ctx.EffectiveChat.Id,authGG.GetInfo(firstArgument, secondArgument),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyToMessageId: ctx.EffectiveMessage.MessageId,
			DisableWebPagePreview: true},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}*/
	db.AddLicenses(db.GetLicenses())
	return nil
}

func LoadTestHandler(updater *ext.Updater)  {
	defer fmt.Println("getInfo module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("test", testHandler))
}

