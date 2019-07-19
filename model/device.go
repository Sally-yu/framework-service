package model

import (
	"fmt"
	"framework-service/database"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Device struct {
	Key           string         `json:"key" form:"key" bson:"key"`                               //id
	Code          string         `json:"code" form:"code" bson:"code"`                            //唯一编号
	Type          string         `json:"type" form:"type" bson:"type"`                            //类型
	Group         string         `json:"group" form:"group" bson:"group"`                         //分组
	Name          string         `json:"name" form:"name" bson:"name"`                            //名称
	Template      DeviceTemplate `json:"template" form:"template" bson:"template"`                //设备模板
	Connect       string         `json:"connect" form:"connect" bson:"connect"`                   //设备连接
	Interval      int64          `json:"interval" form:"interval" bson:"interval"`                //采样间隔
	Model         string         `json:"model" form:"model" form:"model"`                         //规格型号
	GPS           string         `json:"gps" form:"gps" bson:"gps"`                               //gps
	Phone         string         `json:"phone" form:"phone" bson:"phone"`                         //手机号
	Manufacturer  string         `json:"manufacturer" form:"manufacturer" bson:"manufacturer"`    //厂商
	Status        bool           `json:"status" form:"status" bson:"status"`                      //运行状态
	Note          string         `json:"note" form:"note" bson:"note"`                            //备注
	Attrs         []Attribute    `json:"attrs" form:"attrs" bson:"attrs"`                         //设备模板字段
	Time          string         `json:"time" form:"time" bson:"time"`                            //注册时间
	DeviceSetting DeviceSetting  `json:"devicesetting" form:"devicesetting" bson:"devicesetting"` //设备设置参数
	Display       bool           `json:"display" form:"display" bson:"display"`                   //显示卡片
	ServerName    string         `json:"servername" form:"servername" bson:"servername"`          //服务器名称
	ServerAddress string         `json:"serveraddress" form:"serveraddress" bson:"serveraddress"` //服务器地址
}

type DeviceSetting struct {
	CardColor string `json:"cardcolor" form:"cardcolor" bson:"cardcolor"`
}

const (
	deviceDBNAME  = "deviceDB"
	deviceCOLNAME = "deviceCol"
)

//新增设备
func (device *Device) Insert() (bool, string) {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	device.Status = true
	id, _ := uuid.NewRandom()
	device.Key = id.String()
	device.Time = time.Now().Local().Format("2006-01-02 15:04:05")
	err := db.Collection.Insert(&device)
	if err != nil {
		fmt.Println(err)
		return false, "添加设备出错"
	}
	return true, "设备已添加"
}

//全部设备
func (device *Device) All() (error, []Device) {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	res := []Device{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

//key查找设备
func (device *Device) Find() error {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"key": device.Key}).One(&device)
	if err != nil {
		return err
	}
	return nil
}

//name查找设备
func (device *Device) FindByName() error {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"name": device.Name}).One(&device)
	if err != nil {
		return err
	}
	return nil
}

//code查找设备
func (device *Device) FindByCode() error {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"code": device.Code}).One(&device)
	if err != nil {
		return err
	}
	return nil
}

//删除设备
func (device *Device) Remove() error {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Remove(bson.M{"key": device.Key}); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//更新
func (device *Device) Update() error {
	db := database.DbConnection{deviceDBNAME, deviceCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Update(bson.M{"key": device.Key}, device); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

type ResData struct {
	Code  string      `json:"attcode" bson:"code"`
	Type  string      `json:"type" bson:"type"`
	Value interface{} `json:"value" bson:"value"`
}

type Res struct {
	Device string     `json:"device" bson:"device"`
	Data   [] ResData `json:"data" bson:"data"`
}

func (d *Device) GetValue() Res {
	r := Res{}
	for i := range d.Attrs {
		data := ResData{}
		code := d.Attrs[i].Code
		data.Code = code
		err := data.Find()
		if err != nil {
			data.Value = rand.Intn(100)
			data.Type = ""
		}
		fmt.Println(data)
		r.Data = append(r.Data, data)
	}
	r.Device = d.Key
	return r
}

func (r *ResData) Find() error {
	db := database.DbConnection{deviceDBNAME, "deviceValCol", nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"code": r.Code}).One(&r)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	switch v := r.Value.(type) { // 局部变量v是类型转换后的结果
	case nil:
		break
	case int:
		break
	case float32:
		break
	case float64:
		break
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		n10 := math.Pow10(2)
		f = math.Trunc((f+0.5/n10)*n10) / n10
		r.Value = f
		break
	default:
		println("unknown")
		break
	}
	return nil
}
