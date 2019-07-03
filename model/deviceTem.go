package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DeviceTemplate struct {
	Key         string      `json:"key" form:"key" bson:"key"`
	Code        string      `json:"code" form:"code" bson:"code"`
	Name        string      `json:"name" form:"name" bson:"name"`
	Description string      `json:"description" form:"description" bson:"description"`
	Attrs       []Attribute `json:"attrs" form:"attrs" bson:"attrs"`
	Time        string      `json:"time" form:"time" bson:"time"`
}

const TemColNAME = "deviceTemCol"

//生成新的temp
func (temp *DeviceTemplate) New() error {
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	id, _ := uuid.NewRandom()
	temp.Key = id.String()
	temp.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	err := db.Collection.Insert(&temp)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除temp
func (temp *DeviceTemplate) Remove() error {
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Remove(bson.M{"key": temp.Key})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//name查找设备
func (temp *DeviceTemplate) FindByName() error {
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"name": temp.Name}).One(&temp)
	if err != nil {
		return err
	}
	return nil
}

//code查找设备
func (temp *DeviceTemplate) FindByCode() error {
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"code": temp.Code}).One(&temp)
	if err != nil {
		return err
	}
	return nil
}

//更新temp
func (temp *DeviceTemplate) Update() error {
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Update(bson.M{"key": temp.Key}, temp)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//全部temp
func (temp *DeviceTemplate) All() (error, []DeviceTemplate) {
	var list []DeviceTemplate
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(nil).All(&list)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	return nil, list
}

//查找temp
func (temp *DeviceTemplate) Find() error {
	db := database.DbConnection{deviceDBNAME, TemColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"key": temp.Key}).One(&temp)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
