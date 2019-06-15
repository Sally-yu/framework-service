package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

type Notif struct {
	Key     string `json:"key" bson:"key"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	Time    string `json:"time" bson:"time"`
	Type    string `json:"type" bson:"type"`
	Level   int    `json:"level" bson:"level"`
	New     bool   `json:"new" bson:"new"`
}

const (
	DBNAME = "notifDB"
	CONAME = "notifcollection"
)

//生成新通知
func (notif *Notif) Add() (error, *Notif) {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	key, _ := uuid.NewRandom()
	notif.Key = key.String()
	notif.New = true
	notif.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	err := db.Collection.Insert(&notif)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	return nil, notif
}

func (notif *Notif) AllNotif() (error, []Notif) {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	list := []Notif{}
	err := db.Collection.Find(nil).All(&list)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	//时间倒序
	sort.Slice(list, func(i, j int) bool {
		timei,_:=time.Parse("2006-01-02 15:04:05",list[i].Time)
		timej,_:=time.Parse("2006-01-02 15:04:05",list[j].Time)
		return timei.After(timej)
	})
	return nil, list
}

//更新通知状态等
func (notif *Notif) Update() error {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	fmt.Println(notif)
	err := db.Collection.Update(bson.M{"key": notif.Key}, notif)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//移除通知
func (notif *Notif) Remove() error {
	db := database.DbConnection{DBNAME, CONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	fmt.Println(notif)
	err := db.Collection.Remove(bson.M{"key": notif.Key})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
