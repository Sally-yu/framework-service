package model

import (
	"fmt"
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
)

type Chart struct {
	X      int          `json:"x" bson:"x"`
	Y      int          `json:"y" bson:"y"`
	Width  int          `json:"width" bson:"width"`
	Height int          `json:"height" bson:"height"`
	Key    string       `json:"key" bson:"key"`
	Name   string       `json:"name" bson:"name"`
	Option EchartOption `json:"option" bson:"option"`
	Data   interface{}  `json:"data" bson:"data"`
	Type   string       `json:"type" bson:"type"`
}

type EchartOption struct {

}

const (
	TYPELINE="line"
	TYPEBAR="bar"
	TYPEPIE="pie"
	TYPERING="ring"
	TYPERADAR="radar"
	TYPEDASHBOARD="dashboard"
	TYPESINGLE="single"

	CHARTDB="chartdb"
	CHARTCOL="chartcol"
)


func (c *Chart)Upsert() error  {
	db:=database.DbConnection{CHARTDB,CHARTCOL,nil,nil,nil}
	db.ConnDB()
	defer db.CloseDB()
	if _,err := db.Collection.Upsert(bson.M{"key": c.Key}, &c); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
