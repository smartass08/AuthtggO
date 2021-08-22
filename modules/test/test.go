package test

import (
	"AuthtggO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/smartass08/AuthGG-Go/Admin"
)

func testHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if !utils.IsUserOwner(ctx.EffectiveUser.Id) {
		return nil
	}
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return err
	}
	fmt.Println(apiObject.FetchLicenseInfo("Nra-MI73N-SH0JN-72DI2-NA2PA"))
	return nil
}

func LoadTestHandler(updater *ext.Updater)  {
	defer fmt.Println("getInfo module loaded!")
	updater.Dispatcher.AddHandler(handlers.NewCommand("test", testHandler))
}

