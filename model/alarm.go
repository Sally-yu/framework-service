package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Alarm struct {
	Key string `json:"key" form:"key" bson:"key"`
	Name string `json:"name" form:"name" bson:"name"`
	Strategy AlarmStg `json:"strategy" form:"strategy" bson:"strategy"`
	Time time.Time `json:"time" form:"time" bson:"time"`
}

const (
	alarmDBNAME="alarmDB"
	alarmCOLNAME="alarmCol"
)

//新增报警
func (alarm *Alarm) Insert() error{
	db := database.DbConnection{alarmDBNAME, alarmCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	alarm.Time=time.Now().Local()
	id, _ := uuid.NewRandom()
	alarm.Key = id.String()
	err := db.Collection.Insert(&alarm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//全部报警
func (alarm *Alarm) All() (error, []Alarm) {
	db := database.DbConnection{alarmDBNAME, alarmCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	res := []Alarm{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

//key查找报警
func (alarm *Alarm) Find() error {
	db := database.DbConnection{alarmDBNAME, alarmCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"key": alarm.Key}).One(&alarm)
	if err != nil {
		return err
	}
	return nil
}


//更新
func (alarm *Alarm) Update() error {
	db := database.DbConnection{alarmDBNAME, alarmCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Update(bson.M{"key": alarm.Key}, alarm); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}


//删除报警
func (alarm *Alarm) Remove() error {
	db := database.DbConnection{alarmDBNAME, alarmCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Remove(bson.M{"key": alarm.Key}); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
