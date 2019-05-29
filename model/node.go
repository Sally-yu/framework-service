package model

import (
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
)

//树节点，控制为二级树
type Node struct {
	Title    string      `json:"title" bson:"key"`
	Key      string      `json:"key" bson:"key"`
	Author   string      `json:"author" bson:"author"`
	Expanded bool        `json:"expanded" bson:"expanded"`
	Children []ChildNode `json:"children" bson:"children"`
}

func (node *Node) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&node)
	if err != nil {
		return err
	}
	return nil
}

//返回查重结果
func (node *Node) FindTitle(db database.DbConnection, title string) string {
	db.ConnDB()
	db.Collection.Find(bson.M{"Title": title}).One(&node)
	if len(node.Title) > 0 {
		return "0" //重复
	}
	return "1"
}

//以key字段查找
func (node *Node) Find(db database.DbConnection) (error, *Node) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"key": node.Key}).One(&node)
	if err != nil {
		return err, nil
	}
	return nil, node
}

func (node *Node) FindAll(db database.DbConnection) (error, []Node) {
	db.ConnDB()
	res := []Node{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	return nil, res
}

//以key字段删除
func (node *Node) Remove(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Remove(bson.M{"key": node.Key})
	if err != nil {
		return err
	}
	return nil
}

//以key字段更新，先删后插
func (node *Node) Update(db database.DbConnection) error {
	err := node.Remove(db)
	if err != nil {
		return err
	}
	err = node.Save(db)
	if err != nil {
		return err
	}
	return nil
}
