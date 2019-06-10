package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Condition struct {
	Key   string  `json:"key" form:"key" bson:"key"`
	Value float32 `json:"value" form:"value" bson:"value"`
}

type AlarmStg struct {
	Key        string    `json:"key" form:"key" bson:"key"`
	Name       string    `json:"name" form:"name" bson:"name"`
	Code       string    `json:"code" form:"code" bson:"code"`
	Time       string    `json:"time" form:"time" bson:"time"`
	Device     Device    `json:"device" form:"device" bson:"device"`             //设备，内含属性
	Attribute  Attribute `json:"attribute" form:"attribute" bson:"attribute"`    //设备的属性
	ConditionA Condition `json:"conditiona" form:"conditiona" bson:"conditiona"` //报警条件A
	ConditionB Condition `json:"conditionb" form:"conditionb" bson:"conditionb"` //报警条件B
	Level      string    `json:"level" form:"level" bson:"level"`                //报警级别 0提醒 1警报 2严重
	Interval   int64     `json:"interval" form:"interval" bson:"interval"`
	Note       string    `json:"note" form:"note" bson:"note"`
	Status     bool      `json:"status" form:"status" bosn:"status"`
}

const stgColNAME = "alarmStgCol"

//新增报警
func (alarm *AlarmStg) Insert() error {
	db := database.DbConnection{alarmDBNAME, stgColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	alarm.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	id, _ := uuid.NewRandom()
	alarm.Key = id.String()
	err := db.Collection.Insert(&alarm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//全部策略
func (alarm *AlarmStg) All() (error, []AlarmStg) {
	db := database.DbConnection{alarmDBNAME, stgColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	res := []AlarmStg{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

//key查找策略
func (alarm *AlarmStg) Find() error {
	db := database.DbConnection{alarmDBNAME, stgColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"key": alarm.Key}).One(&alarm)
	if err != nil {
		return err
	}
	return nil
}

//更新
func (alarm *AlarmStg) Update() error {
	db := database.DbConnection{alarmDBNAME, stgColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Update(bson.M{"key": alarm.Key}, alarm); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除策略
func (alarm *AlarmStg) Remove() error {
	db := database.DbConnection{alarmDBNAME, stgColNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Remove(bson.M{"key": alarm.Key}); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
