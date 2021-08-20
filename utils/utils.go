package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
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

func ParseInterfaceToInt(i interface{}) int {
	if reflect.TypeOf(i).Name() == "int32" {
		return int(i.(int32))
	}
	return int(i.(int64))
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