package resethwid

import (
	"AuthtggO/db"
	"AuthtggO/helper/authGG"
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"strings"
	"sync"
	"time"
)
var CheckMNutex sync.Mutex

func mutexCheckForced(){
	defer CheckMNutex.Unlock()
	CheckMNutex.Lock()
	time.Sleep(time.Second*3)
}

func ResetUserHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	var (
		admin bool = false
		userKey string
		message string = "Default text"
		valid bool = false
		loaderUsername string
		resetCount int64 = 0
		msg *gotgbot.Message
	)
	if utils.IsUserSudo(ctx.EffectiveUser.Id) || utils.IsUserOwner(ctx.EffectiveUser.Id) {
		admin = true
	}
	if !db.CheckTgUserExists(ctx.EffectiveUser.Id) && !admin{
		b.SendMessage(ctx.EffectiveChat.Id,  "You have not registered with me yet, Press /register for more information",
			&gotgbot.SendMessageOpts{
				ParseMode:             "HTML",
				ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
				DisableWebPagePreview: true},
		)
		return nil
	}

	if !(len(strings.Split(ctx.EffectiveMessage.Text, " ")) == 2)  {
		b.SendMessage(ctx.EffectiveChat.Id,  "Please send me 1 loader username to whom you wanna reset",
			&gotgbot.SendMessageOpts{
				ParseMode:             "HTML",
				ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
				DisableWebPagePreview: true},
		)
		return nil
	}
	loaderUsername =  strings.Split(ctx.EffectiveMessage.Text, " ")[1]
	confirmChannel := make(chan bool)
	if db.CheckTgUserExists(ctx.EffectiveUser.Id) || admin == true{
		chack, err := db.GetOneUser(loaderUsername)
		if !err || chack.UserName == ""{
			valid = false
			message = "User not found"
		} else {
			valid = true
		}
		if valid{
			someKeys := db.GetallKeysOfUser(loaderUsername)
			tempLicense, _ := db.GetOneLicense(someKeys[0])
			userKey = tempLicense.Key
			if admin{
				valid = true
				go addCounttoDb(confirmChannel, tempLicense.Key)
			} else {
				tempTgUSer, _ := db.GetOneTgUser(ctx.EffectiveUser.Id)
				for _, v := range tempTgUSer.KeyUser{
					if strings.ToLower(v) == strings.ToLower(loaderUsername){
						count, err := db.FetchResetCount(tempLicense.Key)
						if err != nil {
							message = err.Error()
							b.SendMessage(ctx.EffectiveChat.Id,  err.Error(),
								&gotgbot.SendMessageOpts{
									ParseMode:             "HTML",
									ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
									DisableWebPagePreview: true},
							)
							return nil
						}
						if count >= 8{
							message = "You have done more than 8 resets, Please contact admin"
							valid = false
						} else {
							go addCounttoDb(confirmChannel, tempLicense.Key)
						}

					}
				}
			}

		}

	} else {
		message = "User not found"
	}
	if valid {
		someChan := make(chan string)
		go func() {
			mutexCheckForced()
			result, err := authGG.ResetHwid(loaderUsername)
			if err != nil {
				time.Sleep(time.Second *2)
				mutexCheckForced()
				result, err = authGG.ResetHwid(loaderUsername)
				if err != nil {
					someChan <- "Error has occured, Please contact admin for manual reset."
				}

			}
			someChan <- result
		}()
		msg, _ = b.SendMessage(ctx.EffectiveChat.Id,  "Hwid Reset Initated, Please wait...",
			&gotgbot.SendMessageOpts{
				ParseMode:             "HTML",
				ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
				DisableWebPagePreview: true},
		)
		for {
			message, ok := <- someChan
			if ok{
				confirmChannel <- true
				if !admin{
					resetCount, _ = db.FetchResetCount(userKey)
					message += fmt.Sprintf("\nTotal resets done this month := <code>%v</code>", resetCount + 1)
				}
				_, err := b.EditMessageText(message, &gotgbot.EditMessageTextOpts{MessageId: msg.MessageId, ParseMode: "HTML", ChatId: msg.Chat.Id})
				if err != nil {
					fmt.Println(err)
					return err
				}
				break
			}
		}
		return nil

	}

	_, err := b.SendMessage(ctx.EffectiveChat.Id, message,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			ReplyToMessageId:      ctx.EffectiveMessage.MessageId,
			DisableWebPagePreview: true},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	confirmChannel <- false
	return nil
}

func LoadResetUserHandler(updater *ext.Updater)  {
	defer fmt.Println("Reset User module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("resetme", ResetUserHandler))
}

func addCounttoDb(confirm chan bool, key string){
	for {
		kek, ok := <- confirm
		if ok{
			if kek {
				db.AddOneToResetCount(key, time.Now())
			}
			return

		}
	}
}