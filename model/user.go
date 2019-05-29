package model

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
	"time"
)

type User struct {
	Key   string `json:"key" bson:"key" form:"key"`
	Phone string `json:"phone" bson:"phone" form:"phone"`
	Email string `json:"email" bson:"email" form:"email"`
	Sex   string `json:"sex" bson:"sex" form:"sex"`
	Group string `json:"group" bson:"group" form:"group"`

	Uname     string    `json:"username" bson:"username" form:"username"`
	Pwd       string    `json:"password" bson:"password" form:"password"`
	Role      string    `json:"role" bson:"role" form:"role"`
	Signtime  time.Time `json:"signtime" bson:"signtime" form:"signtime"`
	Logintime time.Time `json:"logintime" bson:"logintime" formL:"logintime"`
	Status    string    `json:"status" bson:"status" form:"status"` //0未激活未认证，1正常使用，2临时禁用或小黑屋
}

const (
	USERDBNAME = "userdb"
	USERCONAME = "usercol"
)

const (
	SUname = "admin"
	SUpwd  = "admin"
)

//新增用户
func (user *User) Insert() (bool, string) {
	if user.Uname == SUname || user.Exist() {
		return false, "该用户名已存在"
	}
	if len(user.Pwd) < 6 || len(user.Pwd) > 16 {
		return false, "密码长度6至16位"
	}
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	user.Signtime = time.Now().UTC()
	user.Status = "1"
	id, _ := uuid.NewRandom()
	user.Key = id.String()

	err := db.Collection.Insert(&user)
	defer db.CloseDB()
	if err != nil {
		fmt.Println(err)
		return false, "创建用户信息出错"
	}
	return true, "创建用户信息成功"
}

//查找用户，使用用户名
func (user *User) FindByName() error {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"username": user.Uname}).One(&user)
	if err != nil {
		return err
	}
	return nil
}

//手机查找用户
func (user *User) FindByPhone() error {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"phone": user.Phone}).One(&user)
	if err != nil {
		return err
	}
	return nil
}

//邮箱查找用户
func (user *User) FindByEmail() error {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"email": user.Email}).One(&user)
	if err != nil {
		return err
	}
	return nil
}

//用户验证，使用密码
func (user *User) Auth() (bool, string) {
	if SUname == user.Uname { //内置用户
		if SUpwd == user.Pwd {
			return true, "验证成功"
		} else {
			return false, "无效的用户名或密码"
		}
	} else { //普通用户
		var pwd = user.Pwd
		if err := user.FindByName(); err != nil {
			return false, "无效的用户名或密码"
		}
		fmt.Println(pwd)
		fmt.Println(user.Pwd)
		if string(user.Pwd)!= string(pwd) {
			return false, "无效的用户名或密码"
		}

		switch user.Status {
		case "0":
			return false, "用户未激活或未认证"
			break
		case "1":
			return true, ""
			break
		case "2":
			return false, "用户暂不可用"
			break
		default:
			return false, "用户信息不存在"
			break
		}

		return false, ""
	}
}

//登录更新时间
func (user *User) Login() error {
	if SUname == user.Uname && SUpwd == user.Pwd {
		return nil
	}
	user.Logintime = time.Now().UTC()
	err := user.Update()
	if err != nil {
		return err
	}
	return nil
}

//更新
func (user *User) Update() error {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := db.Collection.Update(bson.M{"key": user.Key}, user); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//移除用户
func (user *User) Remove() error {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := user.FindByName(); err != nil {
		return err
	}
	if err := db.Collection.Remove(bson.M{"key": user.Key}); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//用户名，手机，邮箱已使用
func (user *User) Exist() bool {
	name := user.FindByName() == nil
	phone := user.FindByPhone() == nil
	email := user.FindByEmail() == nil
	return name || phone || email
}
