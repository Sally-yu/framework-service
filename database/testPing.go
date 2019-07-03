package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type DbConnectionTest struct {
	Dbname     string
	Cname      string
	Session    *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
	Ip			string
	Port		string
	User		string
	Password	string
}
//类方法，连接MongoDB数据库
func (test *DbConnectionTest) ConnMongoDB() error {
	var err error
	//var url1 = "mongodb://admin:123456@10.24.20.71:28081"
	var url = "mongodb://"
	if test.User != ""  {
		url = url + test.User
	}
	if test.Password != "" {
		url = url + ":" + test.Password + "@"
	}
	if test.Ip != "" {
		url += test.Ip
	}
	if test.Port != ""{
		url = url + ":" + test.Port
	}
	//fmt.Println("url1=",url1)
	//fmt.Println("url=",url)
	test.Session, err = mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	test.Session.SetMode(mgo.Eventual, true) //不缓存连接模式
	test.Database = test.Session.DB(test.Dbname)
	test.Collection = test.Database.C(test.Cname)
	fmt.Println("connect to database:",test.Database)
	fmt.Println("connect to collection: ",test.Collection)
	return nil
}

