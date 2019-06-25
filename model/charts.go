package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Chart struct {
	Key       string      `json:"key" form:"key" bson:"key"`                   //id
	Name      string      `json:"name" form:"name" bson:"name"`                //name
	StartTime string      `json:"starttime" form:"starttime" bson:"starttime"` //起始时间
	EndTime   string      `json:"endtime" form:"endtime" bson:"endtime"`       //结束时间
	Interval  int64       `json:"interval" form:"interval" bson:"interval"`    //采样间隔 ms
	Refresh   bool        `json:"refresh" form:"refresh" bson:"refresh"`       //自动刷新
	Time      string      `json:"time" form:"time" bson:"time"`                //更新时间
	User      string      `json:"user" form:"user" bson:"user"`                //操作人
	Type      string      `json:"type" form:"type" bson:"type"`                //图表类型
	Device    string      `json:"device" form:"device" bson:"device"`          //设备id
	Attrs     []Attribute `json:"attrs" form:"attr" bson:"attr"`               //设备属性
	Option    interface{} `json:"option" form:"option" bson:"option"`          //chart设置
}

const (
	TYPELINE      string = "line"      //折线图
	TYPEPIE       string = "pie"       //饼图
	TYPEHEAT      string = "heat"      //热力图
	TYPETABLE     string = "table"     //表格
	TYPEBAR       string = "bar"       //柱状图
	TYPEMAP       string = "map"       //地图
	TYPEINDICATOR string = "indicator" //指示图
	TYPEDASHBOARD string = "dashboard" //仪表盘
	TYPELIQUID    string = "liquid"    //液体图
	TYPERADAR     string = "radar"     //雷达图
	TYPEPOLAR     string = "polar"     //极坐标

	chartDB  = "chartdb"
	chartCol = "chartcol"
)

func (c *Chart) Upsert() error {
	db := database.DbConnection{chartDB, chartCol, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if c.Key == "" { //新增时
		id, _ := uuid.NewRandom()
		c.Key = id.String()
		c.Interval = 0
		c.Refresh = false
	}
	c.Time = time.Now().Local().Format("2006-01-02 15:04:05") //更新时间
	_, err := db.Collection.Upsert(bson.M{"key": c.Key}, &c)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
