package db

import (
	"AuthtggO/helper/authGG"
	"AuthtggO/logHelper"
	"AuthtggO/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	allUsers      []*User
	allUsersMutex sync.Mutex
)

type User struct {
	UserName   string `json:"username,omitempty"`
	Rank       string `json:"rank,omitempty"`
	ExpiryDate string `json:"expiry_date,omitempty"`
	LastIP     string `json:"Lastip,omitempty"`
	LastLogin  string `json:"lastlogin"`
	Email      string `json:"email"`
	Hwid       string `json:"hwid"`
}

func AddUsersFirstTIme(users []User) error {
	docs := make([]interface{}, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	for _, v := range users {
		doc, _ := toDoc(v)
		docs = append(docs, bson.M{"User": v.UserName, "data": *doc})
	}
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Users")
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Users")
	doc, err := toDoc(user)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, bson.M{"User": user.UserName, "data": doc})
	if err != nil {
		return err
	}
	return nil
}
func FetchAllUsers() (users []*User, err error) {
	logger := logHelper.GetLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Users")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			logger.Error(err)
		} else {
			users = append(users, &User{
				Email:      result["data"].(primitive.M)["email"].(string),
				ExpiryDate: result["data"].(primitive.M)["expirydate"].(string),
				Hwid:       result["data"].(primitive.M)["hwid"].(string),
				Rank:       result["data"].(primitive.M)["rank"].(string),
				LastLogin:  result["data"].(primitive.M)["lastlogin"].(string),
				LastIP:     result["data"].(primitive.M)["lastip"].(string),
				UserName:   result["data"].(primitive.M)["username"].(string),
			})
		}
	}
	return users, nil
}

func GetUsers() (users []User) {
	for _, v := range authGG.FetchAllUsers() {
		users = append(users, User{
			UserName:   v.(map[string]interface{})["username"].(string),
			Rank:       v.(map[string]interface{})["rank"].(string),
			ExpiryDate: v.(map[string]interface{})["expiry_date"].(string),
			LastIP:     v.(map[string]interface{})["lastip"].(string),
			LastLogin:  v.(map[string]interface{})["lastlogin"].(string),
			Email:      v.(map[string]interface{})["email"].(string),
			Hwid:       v.(map[string]interface{})["hwid"].(string),
		})
	}
	return
}
func InitUsers() error {
	users, err := FetchAllUsers()
	if err != nil {
		return err
	}
	for _, v := range users {
		AddUserToLocal(v)
	}
	return nil
}

func AddUserToLocal(trigger *User) {
	allUsersMutex.Lock()
	defer allUsersMutex.Unlock()
	allUsers = append(allUsers, trigger)
}

func UpdateAllUsers() error {
	var operations []mongo.WriteModel
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Users")
	for _, v := range authGG.FetchAllUsers() {
		newUserCheck := true
		for ii, vv := range allUsers {
			if v.(map[string]interface{})["username"].(string) == vv.UserName {
				newUserCheck = false
				if !CompareMapStructUser(v.(map[string]interface{}), *vv) {
					newUser := &User{
						UserName:   v.(map[string]interface{})["username"].(string),
						Rank:       v.(map[string]interface{})["rank"].(string),
						ExpiryDate: v.(map[string]interface{})["expiry_date"].(string),
						LastIP:     v.(map[string]interface{})["lastip"].(string),
						LastLogin:  v.(map[string]interface{})["lastlogin"].(string),
						Email:      v.(map[string]interface{})["email"].(string),
						Hwid:       v.(map[string]interface{})["hwid"].(string),
					}
					allUsers[ii] = newUser
					doc, err := toDoc(newUser)
					if err != nil {
						return err
					}
					operation := mongo.NewUpdateOneModel()
					operation.SetFilter(bson.M{"User": newUser.UserName})
					operation.SetUpdate(bson.M{"$set": bson.M{"User": newUser.UserName, "data": doc}})
					operations = append(operations, operation)
				}
			}
		}
		if newUserCheck {
			newUser := &User{
				UserName:   v.(map[string]interface{})["username"].(string),
				Rank:       v.(map[string]interface{})["rank"].(string),
				ExpiryDate: v.(map[string]interface{})["expiry_date"].(string),
				LastIP:     v.(map[string]interface{})["lastip"].(string),
				LastLogin:  v.(map[string]interface{})["lastlogin"].(string),
				Email:      v.(map[string]interface{})["email"].(string),
				Hwid:       v.(map[string]interface{})["hwid"].(string),
			}

			doc, err := toDoc(newUser)
			if err != nil {
				return err
			}
			operation := mongo.NewInsertOneModel()
			operation.SetDocument(bson.M{"User": newUser.UserName, "data": doc})
			operations = append(operations, operation)
			AddUserToLocal(newUser)
		}
	}
	if len(operations) > 0 {
		bulkOption := options.BulkWriteOptions{}
		bulkOption.SetOrdered(false)
		_, err := collection.BulkWrite(ctx, operations, &bulkOption)
		if err != nil {
			return err
		}
	}
	return nil
}

func CompareMapStructUser(someMap map[string]interface{}, someStruct User) (same bool) {
	if someMap["rank"] != someStruct.Rank || someMap["email"] != someStruct.Email || someMap["username"] != someStruct.UserName || someMap["hwid"] != someStruct.Hwid || someMap["lasptip"] != someStruct.LastIP || someMap["lastlogin"] != someStruct.LastLogin || someMap["expiry_dateA"] != someStruct.ExpiryDate {
		same = false
	}
	same = true
	return
}
func GetOneUser(username string) (User, bool) {
	for _, v := range allUsers {
		if v.UserName == username {
			return *v, true
		}
	}
	return User{}, false
}

func GetallKeysOfUser(username string) []string {
	allkeys := make([]string, 0)
	for _, v := range allLicenses {
		if v.UsedBy == username {
			allkeys = append(allkeys, v.Key)
		}
	}
	return allkeys
}
