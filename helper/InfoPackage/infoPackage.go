package InfoPackage

import (
	"AuthtggO/db"
	"AuthtggO/utils"
	"fmt"
	"strings"
)

func GetInfo(input string, input2 string) string {
	var message string
	switch strings.ToLower(input){
	case "user", "u":
		basicInfo, err := db.GetOneUser(input2)
		if !err {
			message = "No user exists"
			return message
		}
		keys := db.GetallKeysOfUser(basicInfo.UserName)
		message =  fmt.Sprintf(
			"<b>Username</b> : <code>%v</code>\n\n" +
				"<b>Last IP</b> : <a href='https://ifconfig.co/?ip=%v'>%v</a>\n" +
				"<b>Expiry Date</b> : <code>%v</code>\n" +
				"<b>Binded Keys</b> : [\n %v\n]",basicInfo.UserName, basicInfo.LastIP,basicInfo.LastIP, basicInfo.ExpiryDate,strings.Trim(fmt.Sprintf("%s", utils.ParseSliceKeys(keys)), "[]"))
		return message
	case "license", "lic":
		basicInfo, err := db.GetOneLicense(input2)
		if !err {
			message = "No such license exists"
		}
		if basicInfo.Used == "1"{
			basicInfo.Used = "True"
		} else {
			basicInfo.Used = "False"
		}
		message := fmt.Sprintf(
			"<b>License</b> : <code>%v</code>\n\n" +
					"<b>In Use</b> : %v" ,basicInfo.Key, basicInfo.Used )

		if basicInfo.Used == "True"{
			message += fmt.Sprintf("\nUsed By : %v", basicInfo.UsedBy)
		}
		return message
	}
	return "Human mismatch"
}
