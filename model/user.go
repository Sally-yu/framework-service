package model

import (
	"fmt"
	"framework-service/crypt"
	"github.com/google/uuid"
	"github.com/wenzhenxi/gorsa"
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
	Logintime time.Time `json:"logintime" bson:"logintime" form:"logintime"`
	Status    string    `json:"status" bson:"status" form:"status"` //0未激活未认证，1正常使用，2临时禁用或小黑屋

	Img string `json:"img" form:"img" bson:"img"` //存储头像用
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
	if err:=user.FindByName();err==nil{
		return false,"用户名已存在！"
	}
	user.Signtime = time.Now().Local() //本地时区时间。mongo存档时转为UTC，从数据库取出会自动附上时区时差
	user.Status = "1"
	id, _ := uuid.NewRandom()
	user.Key = id.String()

	user.Encrypt()//密文密码存数据库
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
func (user *User) Auth() (bool, string,string) {
	if b:=user.Encrypt();!b{
		fmt.Println("验证出错")
		return false,"验证出错",""
	}
	var Pwd=user.Pwd
	if err := user.FindByName(); err != nil {
		return false, "无效的用户名或密码",""
	}
	fmt.Println(Pwd)
	fmt.Println(user.Pwd)
	if !user.ComparePwd(Pwd) {
		return false, "无效的用户名或密码",""
	}

	switch user.Status {
	case "0":
		return false, "用户未激活或未认证",""
		break
	case "1":
		return true, "",user.Key
		break
	case "2":
		return false, "用户暂不可用",user.Key
		break
	default:
		return false, "用户信息不存在",user.Key
		break
	}

	return false, "",""
}

//用户验证，使用密码
func (user *User) AuthKey() (bool, string) {
	if b:=user.Encrypt();!b{
		fmt.Println("验证出错")
		return false,"验证出错"
	}
	var Pwd=user.Pwd
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
	user.Logintime = time.Now().Local()
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
	//user.Encrypt()  //仅更新普通信息时，不修改密码，user中不包含密码信息，不需加密。
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


//加密 加密密码
func (user *User) Encrypt() bool  {
	var pubkey = crypt.PublicKey
	pwd,err:=gorsa.PublicEncrypt(user.Pwd,pubkey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(pwd)
	user.Pwd=pwd
	return true
}

//加密后的密文 解密比较。数据库存储加密字段
func (user *User)ComparePwd(pwd string) bool{
	var prvkey= crypt.Pirvatekey
	userPwd,_:=gorsa.PriKeyDecrypt(user.Pwd,prvkey)
	enPwd,_:=gorsa.PriKeyDecrypt(pwd,prvkey)
	if string(userPwd)==string(enPwd){
		return true
	}else {
		return false
	}
}