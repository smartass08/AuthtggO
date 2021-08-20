package db

import (
	"AuthtggO/helper/authGG"
	"AuthtggO/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var all_licenses []Licence

type Licence struct {
	Key    string `json:"token,omitempty"`
	Rank   string `json:"rank,omitempty"`
	Used   string `json:"used,omitempty"`
	UsedBy string `json:"used_by,omitempty"`
	Days   string `json:"days,omitempty"`
}

func AddLicenses(licenses []Licence) error{
	docs := make([]interface{}, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	for _, v := range licenses{
		doc, _ := toDoc(v)
		docs = append(docs, bson.M{"License": v.Key, "data": *doc})
	}
	fmt.Println(docs)
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Licenses-test")
	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return nil
}

func GetLicenses() []Licence{
	var licenses []Licence
	for _, v := range authGG.FetchAll(){
		licenses = append(licenses, Licence{
			Key:    v.(map[string]interface{})["token"].(string),
			Rank:    v.(map[string]interface{})["rank"].(string),
			Used:    v.(map[string]interface{})["used"].(string),
			UsedBy:  v.(map[string]interface{})["used_by"].(string),
			Days:    v.(map[string]interface{})["days"].(string),
		})
	}
	return licenses
}

func InitLicense()  {
	
}

func toDoc(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}