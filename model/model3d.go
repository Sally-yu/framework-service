package model

import (
	"fmt"
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
)

type Model3d struct {
	Key      string `json:"key" bson:"key"`
	Code     string `json:"code" bson:"code"`
	Name     string `json:"name" bson:"name"`
	Url      string `json:"url" bson:"url"`
	Cover    string `json:"cover" bson:"cover"`
	Released bool   `json:"released" bson:"released"`
}

const (
	ModelDB  = "modeldb"
	ModelCol = "modelcol"
)

func (m *Model3d) Insert() (bool, string) {
	db := database.DbConnection{ModelDB, ModelCol, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	_, err := db.Collection.Upsert(bson.M{"key": m.Key}, &m)
	if err != nil {
		fmt.Println(err)
		return false, "添加出错"
	}
	return true, "已添加"
}

func (m *Model3d) FindAll() (error, []Model3d) {
	db := database.DbConnection{ModelDB, ModelCol, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	res := []Model3d{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	return nil, res
}

func (m *Model3d) Update() error {
	db := database.DbConnection{ModelDB, ModelCol, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Update(bson.M{"key": m.Key}, m); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
