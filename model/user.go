package model

import (
	"fmt"
	"framework-service/crypt"
	"github.com/google/uuid"
	"github.com/wenzhenxi/gorsa"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"goserver/database"
	"time"
)

type User struct {
	Key       string `json:"key" bson:"key" form:"key"`
	Phone     string `json:"phone" bson:"phone" form:"phone"`
	Email     string `json:"email" bson:"email" form:"email"`
	Sex       string `json:"sex" bson:"sex" form:"sex"`
	Group     string `json:"group" bson:"group" form:"group"`
	Uname     string `json:"username" bson:"username" form:"username"`
	Pwd       string `json:"password" bson:"password" form:"password"`
	Role      string `json:"role" bson:"role" form:"role"` // admin user
	Signtime  string `json:"signtime" bson:"signtime" form:"signtime"`
	Logintime string `json:"logintime" bson:"logintime" form:"logintime"`
	Status    string `json:"status" bson:"status" form:"status"` //0未激活未认证，1正常使用，2临时禁用或小黑屋
	Img       string `json:"img" form:"img" bson:"img"`          //存储头像用
}

const (
	USERDBNAME = "userdb"
	USERCONAME = "usercol"
)

//新增用户
func (user *User) Insert() (bool, string) {
	if len(user.Pwd) < 6 || len(user.Pwd) > 16 {
		return false, "密码长度6至16位"
	}
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	if err := user.FindByName(); err == nil {
		return false, "用户名已存在！"
	}
	user.Signtime = time.Now().Local().Format("2006-01-02 15:04:05")
	user.Status = "1"
	if user.Role != "admin" {
		user.Role = "user" //默认角色权限是user
	}
	id, _ := uuid.NewRandom()
	user.Key = id.String()

	user.Encrypt() //密文密码存数据库
	err := db.Collection.Insert(&user)
	if err != nil {
		fmt.Println(err)
		return false, "创建用户信息出错"
	}
	return true, "创建用户信息成功"
}

//全部用户列表
func (user *User) All() (error, []User) {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	res := []User{}
	err := db.Collection.Find(nil).All(&res)
	if err != nil {
		println(err.Error())
		return err, nil
	}
	return nil, res
}

//查找用户，使用用户名
func (user *User) FindUser() error {
	db := database.DbConnection{USERDBNAME, USERCONAME, nil, nil, nil}
	db.ConnDB()
	defer db.CloseDB()
	err := db.Collection.Find(bson.M{"key": user.Key}).One(&user)
	if err != nil {
		return err
	}
	return nil
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
func (user *User) Auth() (bool, string, string) {
	if b := user.Encrypt(); !b {
		fmt.Println("验证出错")
		return false, "验证出错", ""
	}
	var Pwd = user.Pwd
	if err := user.FindByName(); err != nil {
		return false, "无效的用户名或密码", ""
	}
	fmt.Println(Pwd)
	fmt.Println(user.Pwd)
	if !user.ComparePwd(Pwd) {
		return false, "无效的用户名或密码", ""
	}

	switch user.Status {
	case "0":
		return false, "用户未激活或未认证", ""
		break
	case "1":
		return true, "", user.Key
		break
	case "2":
		return false, "用户暂不可用", user.Key
		break
	default:
		return false, "用户信息不存在", user.Key
		break
	}

	return false, "", ""
}

//用户验证，使用密码
func (user *User) AuthKey() (bool, string) {
	if b := user.Encrypt(); !b {
		fmt.Println("验证出错")
		return false, "验证出错"
	}
	var Pwd = user.Pwd
	if err := user.FindByName(); err != nil {
		return false, "无效的用户名或密码"
	}
	fmt.Println(Pwd)
	fmt.Println(user.Pwd)
	if !user.ComparePwd(Pwd) {
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

//登录更新时间
func (user *User) Login() error {
	user.Logintime = time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Println(user)
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
	if user.Pwd == "" { //密码置空不修改，原密码 对应用户列表修改信息时不输入密码的情况
		u := user
		u.FindUser()
		user.Pwd = u.Pwd
	} else {
		if len(user.Pwd) <= 16 {
			user.Encrypt() //不超过16位，用户信息整个加密传输的情况，对应用户列表修改用户信息时
		}
	}
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
	if err := db.Collection.Remove(bson.M{"key": user.Key}); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}


//
func (user *User) AsAdmin()error  {
	user.Role="admin"
	err:=user.Update()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

//加密 加密密码 bcrypt
func (user *User) Encrypt() bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	pwd := string(hash)  // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(pwd)
	user.Pwd = pwd
	return true
}

//加密后的密文 解密比较。数据库存储加密字段
func (user *User) ComparePwd(pwd string) bool {
	var prvkey = crypt.Pirvatekey
	inPwd, _ := gorsa.PriKeyDecrypt(pwd, prvkey) //解密前台传输的密文
	err:= bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(inPwd)) //bcrypt比较数据库密码
	if err==nil{
		return true
	} else {
		return false
	}
}
