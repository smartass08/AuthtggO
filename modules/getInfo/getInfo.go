package getInfo

import (
	"AuthtggO/helper/InfoPackage"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/smartass08/AuthGG-Go/Admin"
	"strings"
)

func getInfoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserOwner(ctx.EffectiveUser.Id) {
		return nil
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
	msg, err := b.SendMessage(ctx.EffectiveChat.Id,InfoPackage.GetInfo(firstArgument, secondArgument),
		&gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyToMessageId: ctx.EffectiveMessage.MessageId,
		DisableWebPagePreview: true},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if strings.ToLower(firstArgument) == "lic" || strings.ToLower(firstArgument) == "license"{
		apiObject := Admin.Administration{}
		err := apiObject.Init(utils.GetApiHash())
		if err != nil {
			return err
		}
		CreationDateRaw, err := apiObject.FetchLicenseInfo(secondArgument)
		if err != nil{
			return err
		}
		newEditedText := fmt.Sprintf("%v\n<b>Creation Date</b>: <code>%v</code>",msg.Text,CreationDateRaw["created"])
		_, err = b.EditMessageText(newEditedText, &gotgbot.EditMessageTextOpts{
			ChatId: msg.Chat.Id,
			MessageId: msg.MessageId,
			ParseMode: "HTML",
			DisableWebPagePreview: true,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func LoadgetInfoHandler(updater *ext.Updater)  {
	defer fmt.Println("getInfo module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("getinfo", getInfoHandler))
}
