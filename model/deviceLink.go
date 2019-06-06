package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type DeviceLink struct {
	Key       string                 `json:"key" form:"key" bson:"key"`
	DeviceKey string                 `json:"devicekey" form:"devicekey" bson:"devicekey"` //设备的key
	Template  string                 `json:"template" form:"template" bson:"template"`    //关联的模板
	Values    map[string]interface{} `json:"values" form:"values" bson:"values"`          //关联模板的字段和值 map 类型随意
}

const LinkColNAME = "linkCol"

//生成新的link
func (link *DeviceLink) NewLink() error {
	db := database.DbConnection{deviceDBNAME, LinkColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	id, _ := uuid.NewRandom()
	link.Key = id.String()
	err := db.Collection.Insert(&link)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除link
func (link *DeviceLink) Remove() error {
	db := database.DbConnection{deviceDBNAME, LinkColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Remove(bson.M{"key": link.Key})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//更新link
func (link *DeviceLink) Update() error {
	db := database.DbConnection{deviceDBNAME, LinkColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Update(bson.M{"key": link.Key},link)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//全部link
func (link *DeviceLink)All() (error,[]DeviceLink)  {
	var list []DeviceLink
	db := database.DbConnection{deviceDBNAME, LinkColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err:=db.Collection.Find(nil).All(&list)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	return nil, list
}

//查找link
func (link *DeviceLink)Find()  error {
	db := database.DbConnection{deviceDBNAME, LinkColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err:=db.Collection.Find(bson.M{"key":link.Key}).One(&link)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//按设备key查找link
func (link *DeviceLink)FindByDevice()  error {
	db := database.DbConnection{deviceDBNAME, LinkColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err:=db.Collection.Find(bson.M{"devicekey":link.DeviceKey}).One(&link)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}