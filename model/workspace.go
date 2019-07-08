package model

import (
	"fmt"
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	Category string `json:"category" bson:"category"` //图形种类
	Text     string `json:"text" bson:"text"`         //备注用的内容
	Svg      string `json:"svg" bson:"svg"`
	Key      int    `json:"key" bson:"key"`
	Loc      string `json:"loc" bson:"loc"`
	Deviceid string `json:"deviceid" bson:"deviceid"` //关联设备的id
	Status   string `json:"status" bson:"status"`     //运行状态指示
}

type Link struct {
	From       int      `json:"from" bson:"from"`
	To         int      `json:"to" bson:"to"`
	FromPortId string   `json:"fromPortId" bson:"fromPortId"` //与前端页面中绑定的字段key重名，保存连线发出口
	ToPortId   string   `json:"toPortId" bson:"toPortId"`     //与前端页面中绑定的字段key重名，保存连线结束口
	Points     []string `json:"points" bson:"points"`
}

type WorkSpace struct {
	Code          string `json:"code" bson:"code"` //编号 默认自增
	Name          string `json:"name" bson:"name"`
	Key           string `json:"key" bson:"key"`
	Class         string `json:"class" bson:"class"`
	Sence         string `json:"sence" bson:"sence"` //关联的场景
	Note          string `json:"note" bson:"note"`
	Released      bool   `json:"released" bson:"released"` //发布标记
	NodeDataArray []Node `json:"nodeDataArray" bson:"nodeDataArray"`
	LinkDataArray []Link `json:"linkDataArray" bson:"linkDataArray"`
	Cover         string `json:"cover" bson:"cover"` //封面图片
}

func (workspc *WorkSpace) Save(db database.DbConnection) error {
	db.ConnDB()
	fmt.Println(workspc)
	err := db.Collection.Insert(&workspc)
	if err != nil {
		return err
	}
	defer db.CloseDB() //关闭数据库连接，不关闭会增加新的数据库连接
	return nil
}

func (workspc *WorkSpace) FindName(db database.DbConnection, name string) string {
	db.ConnDB()
	db.Collection.Find(bson.M{"name": name}).One(&workspc)
	if len(workspc.Name) > 0 {
		return "0" //重复
	}
	defer db.CloseDB()
	return "1"
}

func (workspc *WorkSpace) Find(db database.DbConnection) (error, *WorkSpace) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"key": workspc.Key}).One(&workspc)
	if err != nil {
		fmt.Println(workspc)
		return err, nil
	}
	defer db.CloseDB()
	return nil, workspc
}

func (workspc *WorkSpace) FindAll(db database.DbConnection) (error, []WorkSpace) {
	db.ConnDB()
	res := []WorkSpace{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

func (workspc *WorkSpace) Release(db database.DbConnection) (error, []WorkSpace) {
	db.ConnDB()
	res := []WorkSpace{}
	err := db.Collection.Find(bson.M{"released": true}).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

func (workspc *WorkSpace) Remove(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Remove(bson.M{"key": workspc.Key})
	if err != nil {
		return err
	}
	defer db.CloseDB()
	return nil
}

func (workspc *WorkSpace) Update(db database.DbConnection) error {
	err := workspc.Remove(db)
	if err != nil {
		return err
	}
	err = workspc.Save(db)
	if err != nil {
		return err
	}
	defer db.CloseDB()
	return nil
}

