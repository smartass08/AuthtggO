package InfoPackage

import (
	"AuthtggO/db"
	"AuthtggO/helper/authGG"
	"AuthtggO/utils"
	"fmt"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

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
			return message
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
		if basicInfo.Created == ""{
			result, _ := authGG.FetchOneLicenseInfo(input2)
			db.FetchAndUpdateOneLicense(input2, result)

			message +=  fmt.Sprintf("\n<b>Creation Date</b>: <code>%v</code>", result["created"])
		} else {
			message += fmt.Sprintf("\n<b>Creation Date</b>: <code>%v</code>", basicInfo.Created)
		}
		return message
	}
	return "Human mismatch"
}

func CreateKey(amount int, duration int, prefix string, channel chan string)  {
	result, err :=  authGG.CreateKey(amount, duration, prefix)
	if err != nil   {
		channel <- err.Error()
		close(channel)
		return
	}
	if result["status"] == "failed" {
		channel <- fmt.Sprintf("%v", result)
		close(channel)
		return
	}
	time.Sleep(time.Second*2)
	for _, v := range result{
		channel <- v.(string)
		wg.Add(1)
		v := v
		fmt.Println(v)
		go func() {
			defer wg.Done()
			fmt.Printf("Adding license to DB %v\n", v.(string))
			err := db.FetchAndAddOneLicense(v.(string))
			if err != nil {
				return
			}
		}()
		time.Sleep(time.Second*2)
	}
	if err != nil {
		channel <- err.Error()
		close(channel)
		return
	}
	close(channel)
}