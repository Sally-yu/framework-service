package model

import (
	"fmt"
	"framework-service/database"
	"gopkg.in/mgo.v2/bson"
)

type DbMgr struct {
	//Key      string   `json:"key" form:"key" bson:"key"`
	//Name     string   `json:"name" form:"name" bson:"name"`
	//Strategy AlarmStg `json:"strategy" form:"strategy" bson:"strategy"`
	//Time     string   `json:"time" form:"time" bson:"time"` //开始时间
	//EndTime  string   `json:"endtime" form:"endtime" bson:"endtime"` //结束时间
	//Duration int64    `json:"duration" form:"duration" bson:"duration"` //持续时长
	//Active   bool     `json:"active" form:"active" bson:"active"` //是否在报警状态

	Servername   string `json:"servername" form:"servername" bson:"servername"`       //服务器名称
	Serverip     string `json:"serverip" form:"serverip" bson:"serverip"`             //服务器IP
	Database     string `json:"database" form:"database" bson:"database"`             //数据库名称
	Serverport   string `json:"serverport" form:"serverport" bson:"serverport"`       //数据库端口
	Databasetype string `json:"databasetype" from:"databasetype" bson:"databasetype"` //数据库类型
	Username     string `json:"username" from:"username" bson:"username"`             //用户名
	Password     string `json:"password" from:"password" bson:"password"`             //密码
	Description  string `json:"description" from:"description" bson:"description"`    //描述
	State        int    `json:"state" bson:"state"`                                   //状态
	//Time     string   `json:"time" form:"time" bson:"time"` //开始时间
}

const (
	dbMgrDBNAME  = "dbMgrDB"
	dbMgrCOLNAME = "dbMgrCol"
)

//添加一条记录
func (mgr *DbMgr) Insert() (bool, string) {
	db := database.DbConnection{dbMgrDBNAME, dbMgrCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Insert(&mgr)
	if err != nil {
		fmt.Println(err)
		return false, "添加数据库出错"
	}
	return true, "数据库已添加"
}

//查询所有的数据库信息
func (mgr DbMgr) All() (error, []DbMgr) {
	db := database.DbConnection{dbMgrDBNAME, dbMgrCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	res := []DbMgr{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	defer db.CloseDB()
	return nil, res
}

//更新一条记录
func (mgr *DbMgr) Update() error {
	db := database.DbConnection{dbMgrDBNAME, dbMgrCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Update(bson.M{"serverip": mgr.Serverip}, mgr); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//删除一条记录
func (mgr *DbMgr) Delete() error {
	db := database.DbConnection{dbMgrDBNAME, dbMgrCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Remove(bson.M{"serverip": mgr.Serverip}); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//查找服务器ip是否存在
func (mgr *DbMgr) FindByServerip() error {
	db := database.DbConnection{dbMgrDBNAME, dbMgrCOLNAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"serverip": mgr.Serverip}).One(&mgr)
	if err != nil {
		return err
	}
	return nil
}

//测试数据库是否连接
func (mgr *DbMgr) TestPing() error {
	fmt.Println("mgr=", mgr)
	fmt.Println("mgr.Databasetype=", mgr.Databasetype)
	if mgr.Databasetype == "MongoDB" {
		db := database.DbConnectionTest{dbMgrDBNAME, dbMgrCOLNAME, nil, nil, nil, mgr.Serverip, mgr.Serverport, mgr.Username, mgr.Password}
		var err = db.ConnMongoDB()
		if err != nil {
			fmt.Println("连接失败！")
			//return err
		} else {
			fmt.Println("连接成功！")
			//return nil
		}
		//defer db.CloseDB()
		return err
	}
	return nil
}
