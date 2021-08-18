package authGG

import (
	"AuthtggO/utils"
	"fmt"
	"github.com/smartass08/AuthGG-Go/Admin"
	"strings"
)

func GetInfo(input string, input2 string) string {
	apiObject := Admin.Administration{}
	err := apiObject.Init(utils.GetApiHash())
	if err != nil {
		return err.Error()
	}
	switch strings.ToLower(input){
	case "user", "u":
		basicInfo, err := apiObject.FetchUserInfo(input2)
		if err != nil {
			return err.Error()
		}
		keys, err := apiObject.FetchAllUsedLicenses(basicInfo["username"].(string))
		if err != nil{
			return err.Error()
		}
		return fmt.Sprintf(
			"<b>Username</b> : <code>%v</code>\n\n" +
					"<b>Last IP</b> : <a href='https://ifconfig.co/?ip=%v'>%v</a>\n" +
					"<b>Expiry Date</b> : <code>%v</code>\n" +
					"<b>Binded Keys</b> : [\n %v\n]",basicInfo["username"], basicInfo["lastip"],basicInfo["lastip"], basicInfo["expiry"],strings.Trim(fmt.Sprintf("%s", utils.ParseSliceKeys(keys)), "[]"))
	case "license", "lic":
		basicInfo, err := apiObject.FetchLicenseInfo(input2)

		if err != nil {
			fmt.Println(err)
			return err.Error()
		}
		message := fmt.Sprintf(
			"<b>License</b> : <code>%v</code>\n\n" +
				"<b>Creation Date</b> : <code>%v</code>\n" +
				"<b>In Use</b> : %v" ,basicInfo["license"], basicInfo["created"], basicInfo["used"])

		if strings.ToLower(basicInfo["used"].(string)) == "true"{
			message += fmt.Sprintf("\nUsed By : %v", basicInfo["used_by"])
		}
		return message
	}
	return "Human mismatch"
}

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