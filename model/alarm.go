package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Alarm struct {
	Key      string   `json:"key" form:"key" bson:"key"`
	Name     string   `json:"name" form:"name" bson:"name"`
	Strategy AlarmStg `json:"strategy" form:"strategy" bson:"strategy"`
	Time     string   `json:"time" form:"time" bson:"time"` //开始时间
	EndTime  string   `json:"endtime" form:"endtime" bson:"endtime"` //结束时间
	Duration int64    `json:"duration" form:"duration" bson:"duration"` //持续时长
	Active   bool     `json:"active" form:"active" bson:"active"` //是否在报警状态
}

const (
	alarmDBNAME  = "alarmDB"
	alarmCOLNAME = "alarmCol"
)

//新增报警
func (alarm *Alarm) Insert() error {
	db := database.DbConnection{alarmDBNAME, alarmCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	alarm.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	alarm.Active = true
	id, _ := uuid.NewRandom()
	alarm.Key = id.String()
	err := db.Collection.Insert(&alarm)
	if err != nil {
		fmt.Println(err.Error())
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

//结束报警，更新结束时间和时长
func (alarm *Alarm) End() error {
	alarm.EndTime = time.Now().Local().Format("2006-01-02 15:04:05") //添加结束时间
	if alarm.Active {
		times, _ := time.Parse("2006-01-02 15:04:05", alarm.Time)
		timee, _ := time.Parse("2006-01-02 15:04:05", alarm.EndTime)
		long := (timee.Unix() - times.Unix()) / 1000
		alarm.Duration = long
		alarm.Active = false
		err := alarm.Update()
		if err != nil {
			return err
		}
	}
	return nil
}
