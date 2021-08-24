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
	Key     string `json:"token,omitempty"`
	Rank    string `json:"rank,omitempty"`
	Used    string `json:"used,omitempty"`
	UsedBy  string `json:"used_by,omitempty"`
	Days    string `json:"days,omitempty"`
	Created string `bson:"created,omitempty"`
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
	AddLicenseToLocal(&licence)
	return nil
}

func UpdateLicense(licence Licence) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := DatabaseClient.Mongo.Database(utils.GetDataBaseName()).Collection("Licenses")
	doc, err := toDoc(licence)
	if err != nil {
		return err
	}
	filter := bson.M{
		"License": licence.Key,
	}
	_, err = collection.UpdateOne(ctx, filter, bson.D{{
		"$set", bson.M{"License": licence.Key, "data": doc}}})
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
			licenseToBeAdded := &Licence{
				Key:    result["data"].(primitive.M)["key"].(string),
				Rank:   result["data"].(primitive.M)["rank"].(string),
				Used:   result["data"].(primitive.M)["used"].(string),
				UsedBy: result["data"].(primitive.M)["usedby"].(string),
				Days:   result["data"].(primitive.M)["days"].(string),
			}
			if result["data"].(primitive.M)["created"] != nil {
				licenseToBeAdded.Created = result["data"].(primitive.M)["created"].(string)
			}
			licences = append(licences, licenseToBeAdded)
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
					time.Sleep(time.Second * 2)
					tempLicenseFetch, err := authGG.FetchOneLicenseInfo(vv.Key)
					if err == nil {
						newLicense.Created = tempLicenseFetch["created"].(string)
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
			time.Sleep(time.Second * 3)
			tempLicenseFetch, err := authGG.FetchOneLicenseInfo(v.(map[string]interface{})["token"].(string))
			if err == nil {
				newLicense.Created = tempLicenseFetch["created"].(string)
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

func FetchAndAddOneLicense(license string) error {
	result, err := authGG.FetchOneLicenseInfo(license)
	var tempVar string
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	for _, v := range authGG.FetchAllLicenses() {
		if v.(map[string]interface{})["token"].(string) == license {
			tempVar = v.(map[string]interface{})["days"].(string)
		}
	}
	err = AddLicense(Licence{
		Key:     result["license"].(string),
		Rank:    result["rank"].(string),
		Used:    result["used"].(string),
		UsedBy:  result["used_by"].(string),
		Days:    tempVar,
		Created: result["created"].(string),
	})
	if err != nil {
		return nil
	}
	return err
}
func FetchAndUpdateOneLicense(license string, info map[string]interface{}) error {
	if info == nil {
		time.Sleep(time.Second * 2)
		info, _ = authGG.FetchOneLicenseInfo(license)
	}
	tempNonUpdated, check := GetOneLicense(license)
	if !check {
		time.Sleep(time.Second * 2)
		err := FetchAndAddOneLicense(license)
		if err != nil {
			return err
		}
	}
	toBeUpdatedLicense := Licence{
		Key:     info["license"].(string),
		Rank:    info["rank"].(string),
		Used:    info["used"].(string),
		UsedBy:  info["used_by"].(string),
		Days:    tempNonUpdated.Days,
		Created: info["created"].(string),
	}
	err := UpdateLicense(toBeUpdatedLicense)
	if err != nil {
		return nil
	}
	UpdateOneLicenseLocal(toBeUpdatedLicense)
	return err
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
		if v.Key == license {
			return *v, true
		}
	}
	return Licence{}, false
}

func UpdateOneLicenseLocal(licence Licence) {
	for i, v := range allLicenses {
		if v.Key == licence.Key {
			allLicenses[i] = &licence
		}
	}
}
