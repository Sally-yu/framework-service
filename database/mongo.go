package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

const (
	//MONGODB_URL = "10.24.20.71:27017"
	//MONGODB_URL = "127.0.0.1:27017"
	MONGODB_URL="mongodb://admin:123456@10.24.20.71:28081" //绑定端口到28081
)

type DbConnection struct {
	Dbname     string
	Cname      string
	Session    *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
}

//类方法，连接数据库
func (db *DbConnection) ConnDB() error {
	var err error
	db.Session, err = mgo.Dial(MONGODB_URL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//defer session.Close()
	//session.SetMode(mgo.Monotonic,true)
	db.Session.SetMode(mgo.Eventual, true) //不缓存连接模式
	db.Database = db.Session.DB(db.Dbname)
	db.Collection = db.Database.C(db.Cname)
	fmt.Println("connect to database:",db.Database)
	fmt.Println("connect to collection: ",db.Collection)
	return nil
}

//类方法,切换collection
func (db *DbConnection) SwitchCollection(cname string) {
	if cname != db.Cname { //切换collection 重新生成session
		db.Cname = cname
		db.CloseDB()
		db.ConnDB()
	}
}

//类方法，关闭数据连接session
func (db *DbConnection) CloseDB() {
	//连接数据库时session模式为Eventual模式，不关session也可以，不会缓存
	defer db.Session.Close()
	fmt.Println("close database connection ")
}
