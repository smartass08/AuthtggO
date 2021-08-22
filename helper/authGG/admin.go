package authGG

import (
	"AuthtggO/utils"
	"fmt"
	"github.com/smartass08/AuthGG-Go/Admin"
	"strings"
)

func CreateKey(amount int, duration int, prefix string) string {
	apiObject := Admin.Administration{}
	var keys []string
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return err.Error()
	}
	key, err := apiObject.GenerateLicense(amount, duration, prefix)
	if err != nil {
		return err.Error()
	}
	for _, v := range key{
		keys = append(keys, v.(string))
	}
	return fmt.Sprintf(" %v", strings.TrimLeft(fmt.Sprintf("%s", utils.ParseSliceKeys(keys)), "[]"))
}

func DeleteKey(keyRaw []string) string {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return err.Error()
	}
	message := ""
	for _, v := range keyRaw{
		check, err := apiObject.DeleteKey(v)
		if err != nil {
			message += fmt.Sprintf("`%v` : %v\n", v, err.Error())
			continue
		}
		message += fmt.Sprintf("`%v` : %v\n", v, check["info"])
	}
	return message
}

func FetchAllLicenses() map[string]interface{} {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return nil
	}
	kek,_ := apiObject.FetchAllLicenseInfo()
	return kek
}

func FetchAllUsers() map[string]interface{} {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return nil
	}
	kek,_ := apiObject.FetchAllUserInfo()
	return kek
}