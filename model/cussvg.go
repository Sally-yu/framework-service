package model

import (
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
)

type Svg struct {
	Svg      string `json:"svg" bson:"svg"`
	Deviceid string `json:"deviceid" bson:"deviceid"`
	Status   string `json:"status" bson:"status"`
}

type Cus struct {
	Divid   string `json:"divid" bson:"divid"`
	Name    string `json:"name" bson:"name"`
	Display bool   `json:"display" bson:"display"`
	Svg     []Svg  `json:"svg" bson:"svg"`
}


//单个cus
func (cus *Cus) Find(db database.DbConnection) (error, *Cus) {
	db.ConnDB()
	err := db.Collection.Find(bson.M{"divid": cus.Divid}).One(&cus)
	if err != nil {
		return err, nil
	}
	defer db.CloseDB()
	return nil, cus
}

//返回所有自定义分组信息
func (cus *Cus) FindAll(db database.DbConnection) (error, []Cus) {
	db.ConnDB()
	res := []Cus{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		print(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

func (cus *Cus) Save(db database.DbConnection) error {
	db.ConnDB()
	err := db.Collection.Insert(&cus)
	if err != nil {
		return err
	}
	defer db.CloseDB()
	return nil
}

func (cus *Cus) Remove(db database.DbConnection) error {
	db.ConnDB()
	db.Collection.Remove(bson.M{"divid": cus.Divid})
	defer db.CloseDB()
	return nil
}

func (cus *Cus) Update(db database.DbConnection) error {
	err := cus.Remove(db)
	if err != nil {
		return err
	}
	err = cus.Save(db)
	if err != nil {
		return err
	}
	defer db.CloseDB()
	return nil
}

//去重
func (cus *Cus) RemoveRepeat() {
	arr:=[]Svg{}
	svg:=cus.Svg
	for i := 0; i < len(svg); i++ {
		repeat := false
		for j := i + 1; j < len(svg); j++ {
			if svg[i].Svg == svg[j].Svg {
				repeat = true
				break
			}
		}
		if !repeat {
			arr = append(arr, svg[i])
		}
	}
	cus.Svg=arr
}