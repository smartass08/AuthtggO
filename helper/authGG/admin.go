package authGG

import (
	"AuthtggO/utils"
	"fmt"
	"github.com/smartass08/AuthGG-Go/Admin"
	"time"
)

func CreateKey(amount int, duration int, prefix string) (map[string]interface{}, error) {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return nil, err
	}
	key, err := apiObject.GenerateLicense(amount, duration, prefix)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func DeleteKey(keyRaw []string, channel chan string) error {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return err
	}
	for _, v := range keyRaw{
		fmt.Println(v)
		check, err := apiObject.DeleteKey(v)
		if err != nil {
			channel <- fmt.Sprintf("`%v` : %v\n", v, err.Error())
			continue
		}
		channel <- fmt.Sprintf("`%v` : %v\n", v, check["info"])
		time.Sleep(time.Second*2)
	}
	close(channel)
	return nil
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

func ResetHwid(username string) (string, error) {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return "", nil
	}
	result, err := apiObject.ResetHwid(username)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("<code>%v</code>\nStatus:  <code>%v</code>", result["info"].(string),  result["status"].(string)), nil

}

func FetchOneLicenseInfo(license string) (map[string]interface{}, error){
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return nil, err
	}
	result, err := apiObject.FetchLicenseInfo(license)
	if err != nil {
		return result, err
	}
	return result, nil
}