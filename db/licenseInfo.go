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
	allLicenses      []*Licence
	allLicensesMutex sync.Mutex
)

type Licence struct {
	Key    string `json:"token,omitempty"`
	Rank   string `json:"rank,omitempty"`
	Used   string `json:"used,omitempty"`
	UsedBy string `json:"used_by,omitempty"`
	Days   string `json:"days,omitempty"`
}

func AddLicensesFirstTIme(licenses []Licence) error {
	docs := make([]interface{}, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	for _, v := range licenses {
		doc, _ := toDoc(v)
		docs = append(docs, bson.M{"License": v.Key, "data": *doc})
	}
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Licenses")
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func AddLicense(licence Licence) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Licenses")
	doc, err := toDoc(licence)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, bson.M{"License": licence.Key, "data": doc})
	if err != nil {
		return err
	}
	return nil
}

func FetchAllLicences() (licences []*Licence, err error) {
	logger := logHelper.GetLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Licenses")
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
			licences = append(licences, &Licence{
				Key:    result["data"].(primitive.M)["key"].(string),
				Rank:   result["data"].(primitive.M)["rank"].(string),
				Used:   result["data"].(primitive.M)["used"].(string),
				UsedBy: result["data"].(primitive.M)["usedby"].(string),
				Days:   result["data"].(primitive.M)["days"].(string),
			})
		}
	}
	return licences, nil
}

func GetLicenses() []Licence {
	var licenses []Licence
	for _, v := range authGG.FetchAllLicenses() {
		licenses = append(licenses, Licence{
			Key:    v.(map[string]interface{})["token"].(string),
			Rank:   v.(map[string]interface{})["rank"].(string),
			Used:   v.(map[string]interface{})["used"].(string),
			UsedBy: v.(map[string]interface{})["used_by"].(string),
			Days:   v.(map[string]interface{})["days"].(string),
		})
	}
	return licenses
}

func InitLicense() error {
	licences, err := FetchAllLicences()
	if err != nil {
		return err
	}
	for _, v := range licences {
		AddLicenseToLocal(v)
	}
	return nil
}

func AddLicenseToLocal(trigger *Licence) {
	allLicensesMutex.Lock()
	defer allLicensesMutex.Unlock()
	allLicenses = append(allLicenses, trigger)
}

func UpdateAlllicenses() error {
	var operations []mongo.WriteModel
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Licenses")
	for _, v := range authGG.FetchAllLicenses() {
		newLicense := true
		for ii, vv := range allLicenses {
			if v.(map[string]interface{})["token"].(string) == vv.Key {
				newLicense = false
				if !CompareMapStructLicense(v.(map[string]interface{}), *vv) {
					newLicense := &Licence{
						Key:    vv.Key,
						Rank:   v.(map[string]interface{})["rank"].(string),
						Used:   v.(map[string]interface{})["used"].(string),
						UsedBy: v.(map[string]interface{})["used_by"].(string),
						Days:   v.(map[string]interface{})["days"].(string),
					}
					allLicenses[ii] = newLicense
					doc, err := toDoc(newLicense)
					if err != nil {
						return err
					}
					operation := mongo.NewUpdateOneModel()
					operation.SetFilter(bson.M{"License": newLicense.Key})
					operation.SetUpdate(bson.M{"$set": bson.M{"License": newLicense.Key, "data": doc}})
					operations = append(operations, operation)
				}
			}
		}
		if newLicense {
			newLicense := &Licence{
				Key:    v.(map[string]interface{})["token"].(string),
				Rank:   v.(map[string]interface{})["rank"].(string),
				Used:   v.(map[string]interface{})["used"].(string),
				UsedBy: v.(map[string]interface{})["used_by"].(string),
				Days:   v.(map[string]interface{})["days"].(string),
			}
			doc, err := toDoc(newLicense)
			if err != nil {
				return err
			}
			operation := mongo.NewInsertOneModel()
			operation.SetDocument(bson.M{"License": newLicense.Key, "data": doc})
			operations = append(operations, operation)
			AddLicenseToLocal(newLicense)
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

func CompareMapStructLicense(someMap map[string]interface{}, someStruct Licence) (same bool) {
	if someMap["rank"] != someStruct.Rank || someMap["used_by"] != someStruct.UsedBy || someMap["days"] != someStruct.Days {
		same = false
	}
	same = true
	return
}

func toDoc(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func GetOneLicense(license string) (Licence, bool) {
	for _, v := range allLicenses {
		if v.Key == license{
			return *v, true
		}
	}
	return Licence{}, false
}