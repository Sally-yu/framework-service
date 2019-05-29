package model

import (
	"fmt"
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Notif struct {
	Key   string    `json:"key" bson:"key"`
	Title string    `json:"title" bson:"title"`
	Time  time.Time `json:"time" bson:"time"`
	Type  string    `json:"type" bson:"type"`
	Level int       `json:"level" bson:"level"`
	New   bool      `json:"new" bson:"new"`
}

const (
	DBNAME = "notifDB"
	CONAME = "notifcollection"
)

//生成新通知
func (notif *Notif) NewNotif() (error, *Notif) {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	fmt.Println(notif)
	notif.New = true
	notif.Time = time.Now().UTC()
	err := db.Collection.Insert(&notif)
	if err != nil {
		return err, nil
	}
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	return nil, notif
}

//更新通知状态等
func (notif *Notif) Update() (error, *Notif) {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	fmt.Println(notif)
	err := db.Collection.Update(bson.M{"key":notif.Key},notif)
	if err != nil {
		return err, nil
	}
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	return nil, notif
}

//移除通知
func (notif *Notif)Remove() error  {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	fmt.Println(notif)
	err := db.Collection.Remove(bson.M{"key":notif.Key})
	if err != nil {
		return err
	}
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	return nil
}