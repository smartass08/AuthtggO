package utils

import (
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"time"
)

const ConfigJsonPath string = "config.json"

const MaxMessageTextLength int = 4000

type ConfigJson struct {
	BotToken  string  `json:"bot_token"`
	SudoUsers   []int64 `json:"sudo_users"`
	MongoUri    string  `json:"mongo_uri"`
	OwnerId     int64 `json:"owner_id"`
	ApiHash		string `json:"apiHashGG"`
	DataBaseName string `json:"dbName"`
	ModChat		int64	`json:"mod_chat"`
}

var Config *ConfigJson = InitConfig()

func InitConfig() *ConfigJson {
	file, err := ioutil.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}
	var Config ConfigJson
	err = json.Unmarshal([]byte(file), &Config)
	if err != nil {
		log.Fatal(err)
	}
	Config.SudoUsers = append(Config.SudoUsers, Config.OwnerId)
	log.Println(Config.SudoUsers)
	return &Config
}

func IsUserSudo(userId int64) bool {
	for _, i := range Config.SudoUsers {
		if i == userId {
			return true
		}
	}
	return false
}

func IsUserOwner(userId int64) bool {
	return Config.OwnerId == userId
}

func GetDbUri() string {
	return Config.MongoUri
}

func GetBotToken() string {
	return Config.BotToken
}

func GetOwnerId() int64 {
	return Config.OwnerId
}

func GetApiHash() string {
	return Config.ApiHash
}

func GetDataBaseName() string {
	return Config.DataBaseName
}

func GetModChatId() int64 {
	return Config.ModChat
}

func GenerateUserLink(userId int64, parseMode string) string {
	if parseMode == "HTML" || parseMode == "html" {
		return fmt.Sprintf("<a href='tg://user?id=%d'>User</a>", userId)
	}
	return fmt.Sprintf("[User](tg://user?id=%d)", userId)
}

func RemoveDuplicates(elements []int64) []int64 {
	encountered := map[int64]bool{}
	result := []int64{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func ParseInterfaceToInt64(i interface{}) int64 {
	if reflect.TypeOf(i).Name() == "int64" {
		return int64(i.(int64))
	}
	return int64(i.(int32))
}

func ParseSliceKeys(keys []string) []string {
	var updatedKeys []string
	for i, v := range keys{
		updatedKeys = append(updatedKeys, fmt.Sprintf("<code>%v</code>", v))
		if len(keys)-1 == i{
			continue
		}
		updatedKeys = append(updatedKeys, " ,\n")
	}
	return updatedKeys
}

// ReturnRestCounts Returns total amount of resets done in int || Returns -1 if error
func ReturnRestCounts(timeSlice []int64, endDate string, validityInString string) (int64, error){
	expiryDate, err := dateparse.ParseAny(endDate)
	if err != nil {
		//Return 10 if error
		return -1, err
	}
	validityInt, err := strconv.Atoi(validityInString)
	if err != nil {
		return -1, err
	}
	startDate := expiryDate.AddDate(0,0, -validityInt)
	StartDateCount := int(startDate.Month())
	EndDateCount := int(expiryDate.Month())
	monthCount := 0
	for i := StartDateCount; i < EndDateCount; i++{
		monthCount += 1
	}
	var resetCount int64
	for _, v := range timeSlice{
		thisDate := time.Unix(0, v)
		tempCount := 0
		for ii := 0; ii <= monthCount; ii++{
			if inTimeSpan(startDate.AddDate(0,tempCount, 0 ), startDate.AddDate(0,tempCount+1, 0 ), time.Now()){
				if inTimeSpan(startDate.AddDate(0,tempCount, 0 ), startDate.AddDate(0,tempCount+1, 0 ), thisDate){
					resetCount += 1
					break
				}
			} else {
				tempCount += 1
			}
		}

	}
	return resetCount, nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
