package db

import (
	"AuthtggO/logHelper"
	"AuthtggO/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

var (
	alltgUsers      []*TGUser
	allTgUsersMutex sync.Mutex
)
type TGUser struct {
	TgID   int64 `json:"tgID"`
	KeyUser []string `json:"keyUser"`
	BindedKeys []string `json:"bindedKeys"`
}

func InitTgUser() error {
	tgUsers, err := FetchAllTgUsers()
	if err != nil {
		return err
	}
	for _, v := range tgUsers{
		AddtgUserToLocal(*v)
	}
	return nil
}


func FetchAllTgUsers() (users []*TGUser, err error) {
	logger := logHelper.GetLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("TgUsers")
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
			tempLicenseSlice := make([]string, 0)
			tempKeyUserSlice := make([]string, 0)
			for _, vv := range result["data"].(primitive.M)["bindedkeys"].(primitive.A){
				tempLicenseSlice = append(tempLicenseSlice, vv.(string))
			}
			for _, vv := range result["data"].(primitive.M)["keyuser"].(primitive.A){
				tempKeyUserSlice = append(tempKeyUserSlice, vv.(string))
			}
			users = append(users, &TGUser{
				TgID:          result["data"].(primitive.M)["tgid"].(int64),
				KeyUser:       tempKeyUserSlice,
				BindedKeys:    tempLicenseSlice,
			})
		}
	}
	return users, nil
}

func AddTgUserToDB(user TGUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("TgUsers")
	doc, err := toDoc(user)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, bson.M{"User": user.TgID, "data": doc})
	if err != nil {
		return err
	}
	return nil
}
func AddtgUserToLocal(trigger TGUser) {
	allTgUsersMutex.Lock()
	defer allTgUsersMutex.Unlock()
	alltgUsers = append(alltgUsers, &trigger)
}

func UpdateTgUserDB(user TGUser) error{
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("TgUsers")
	doc, err := toDoc(user)
	if err != nil {
		return err
	}
	filter := bson.M{
		"License": user.TgID,
	}
	_, err = collection.UpdateOne(ctx, filter, bson.D{{
		"$set", bson.M{"License": user.TgID, "data": doc}}})
	if err != nil {
		return err
	}
	return nil
}

func UpdateTgUserLocal(user TGUser) bool {
	allTgUsersMutex.Lock()
	defer allTgUsersMutex.Unlock()
	for i,v := range alltgUsers{
		if v.TgID == user.TgID {
			alltgUsers[i] = &user
			return true
		}
	}
	return false
}

func UpdateTgUser(user TGUser) bool {
	err := UpdateTgUserDB(user); if err != nil {
		return false
	}
	errr := UpdateTgUserLocal(user); if errr != true {
		return false
	}
	return true
}

func AddtgUser(user TGUser) error {
	err := AddTgUserToDB(user); if err != nil {
		return err
	}
	AddtgUserToLocal(user)
	return nil
}

func CheckTgUserExists(userId int64) bool{
	for _,v := range alltgUsers{
		if v.TgID == userId{
			return true
		}
	}
	return false
}
func GetOneTgUser(userId int64) (*TGUser, bool) {
	for _,v := range alltgUsers{
		if v.TgID == userId{
			return v, true
		}
	}
	return &TGUser{}, false
}

func IstgUserScammer(license string, username string) bool{
	for _,v := range alltgUsers{
		for _, vv := range v.KeyUser{
			if vv == username {
				return true
			}
		}
		for _, vv := range v.BindedKeys{
			if vv == license{
				return true
			}
		}
	}
	return false
}