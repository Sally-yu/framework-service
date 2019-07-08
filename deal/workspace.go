package deal

import (
	"encoding/json"
	"fmt"
	"framework-service/database"
	"framework-service/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)


type ResponseW struct {
	Status string
}

const (
	path         = "/usr/local/nginx/html/assets/img"
	uploadPath   = "/usr/local/nginx/html/assets/upload"
	DbnameWork   = "imgdb"
	CnameWork    = "imgcollection"
	DocDnameWork = "docdb"
)

//const path="c:/img"

//跨域头
func CorsHeaderWork(w http.ResponseWriter) http.ResponseWriter {
	println("CorsHeaderWork")
	w.Header().Set("Access-Control-Allow-Origin", "*")                                                      //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,X-Requested-With,Authorization,text/html") //header的类型	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json") //返回数据格式是json
	println("Cors Done")
	return w
}




//上传保存图片
func PathServer(w http.ResponseWriter, r *http.Request, p string) {
	CorsHeaderWork(w) //优先处理跨域，否则后续函数不会执行
	if "POST" == r.Method {
		fmt.Println(r)
		file, handler, err := r.FormFile("file") //antd上传控件formdata中文件key为file
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}

		filename := handler.Filename //取文件名
		defer file.Close()
		if _, err := os.Stat(p); err == nil {
			fmt.Println("path exists 1", p)
		} else {
			fmt.Println("path not exists ", p)
			err := os.MkdirAll(p, 0711)
			// check again
			if err != nil {
				io.WriteString(w, "Error creating directory")
				return
			}
			if _, err := os.Stat(p); err == nil {
				fmt.Println("path exists 2", p)
			}
		}

		if _, err := os.Stat(p + "/" + filename); err == nil {
			os.Remove(p + "/" + filename) //删除文件
		}
		f, err := os.Create(p + "/" + filename) //创建文件
		if err != nil {
			return
		}
		defer f.Close()
		var data ResponseW
		data.Status = "success"
		fmt.Println(filename)
		json.NewEncoder(w).Encode(data)
		io.Copy(f, file) //复制文件内容
	}
}

//保存到数据库
func SaveLink(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var jsono map[string]string
		json.Unmarshal(body, &jsono) //json解析
		db := database.DbConnection{DbnameWork, CnameWork, nil, nil, nil}
		img := model.Img{jsono["deviceid"], jsono["imgurl"]}
		err := img.Save(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data := ResponseW{}
		data.Status = "success"
		json.NewEncoder(w).Encode(data)
	}
}

//返回imgurl
func FindImg(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var device map[string]string
		json.Unmarshal(body, &device) //json解析
		img := model.Img{device["deviceid"], ""}
		err, _ := img.Find(database.DbConnection{DbnameWork, CnameWork, nil, nil, nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data := ResponseW{}
		data.Status = img.Imgurl
		json.NewEncoder(w).Encode(data)
	}
}

func BackImg(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var post map[string]string
		json.Unmarshal(body, &post) //json解析
		back := model.Back{}
		err, _ := back.Find(database.DbConnection{DbnameWork, CnameWork, nil, nil, nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		json.NewEncoder(w).Encode(back)
	}
}

func WorkSpace(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		var post = struct {
			Opt       string
			Workspace model.WorkSpace
		}{}
		worklist := []model.WorkSpace{}
		json.Unmarshal(body, &post) //json解
		var err error
		switch post.Opt {
		case "save": //保存
			err = post.Workspace.Save(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(post.Workspace) //response一个workspace
			break
		case "find":
			err, _ = post.Workspace.Find(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(post.Workspace) //response一个workspace
			break
		case "all":
			err, worklist = post.Workspace.FindAll(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(worklist) //response一个workspace
			break
		case "released":
			err, worklist = post.Workspace.Release(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode(worklist) //response一个workspace
			break
		case "delete":
			err = post.Workspace.Remove(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil})
			json.NewEncoder(w).Encode("delete success")
			break
		default:
			break
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	PathServer(w, r, uploadPath)
}

func SaveSvg(w http.ResponseWriter, r *http.Request) {
	PathServer(w, r, path)
}

//更新指定的自定义分组
func UpdateCus(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		cus := model.Cus{}
		json.Unmarshal(body, &cus) //json解析
		c := model.Cus{Divid: cus.Divid}
		err, _ := c.Find(database.DbConnection{DocDnameWork, "cus", nil, nil, nil})
		fmt.Println("c:", c.Svg)
		fmt.Println("cus:", cus.Svg)
		if err == nil {
			a := append(cus.Svg, c.Svg...) //合并数组，追加svg
			cus.Svg = a
		}
		cus.RemoveRepeat() //svg去重
		for i := 0; i < len(cus.Svg); i++ {
			cus.Svg[i].Svg = strings.Replace(cus.Svg[i].Svg, ".svg", "", 1) //去.svg后缀
		}
		err = cus.Update(database.DbConnection{DocDnameWork, "cus", nil, nil, nil})
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		data := ResponseW{}
		data.Status = "update"
		json.NewEncoder(w).Encode(data)
	}
}

//获取所有的自定义分组
func CusSvg(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "GET" == r.Method {
		svgList := []model.Cus{}
		svg := model.Cus{}
		var err error
		err, svgList = svg.FindAll(database.DbConnection{DocDnameWork, "cus", nil, nil, nil})
		json.NewEncoder(w).Encode(svgList) //response一个workspace
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
	}
}

func FindName(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	if "POST" == r.Method {
		body, _ := ioutil.ReadAll(r.Body) //获取post的数据
		work := model.WorkSpace{}
		var name = ""
		json.Unmarshal(body, &name) //json解析
		res := ResponseW{}
		res.Status = work.FindName(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil}, name)
		json.NewEncoder(w).Encode(res) //response一个workspace
	}
}

func AutoCode(w http.ResponseWriter, r *http.Request) {
	CorsHeaderWork(w)
	//if "GET" == r.Method {
		work := model.WorkSpace{}
		_, list := work.FindAll(database.DbConnection{DocDnameWork, "workspace", nil, nil, nil})
		var codes [] int
		if len(list)<1{
			json.NewEncoder(w).Encode(1)
			return
		}
		for _, value := range list {
			if strings.HasPrefix(value.Code, "TOPO-") {
				i, _ := strconv.Atoi(value.Code[5:])
				codes = append(codes, i)
			}
		}
		sort.Ints(codes)
		if len(codes)<1{
			json.NewEncoder(w).Encode(1)
			return
		}
		json.NewEncoder(w).Encode(codes[len(codes)-1]+1)
	//}
}

//topob部分的后台文件服务
func FileServer(){
	//预制图标存取
	http.HandleFunc("/assets/img", SaveSvg)          //保存svg 自带get文件服务器
	http.HandleFunc("/assets/img/save", SaveLink)    //保存设备和svg的联系
	http.HandleFunc("/assets/upload", Upload)        //上传自定义svg ，自带get文件服务器
	http.HandleFunc("/assets/img/deviceid", FindImg) //加载图片
	http.HandleFunc("/assets/img/back", BackImg)

	//布局、自定义图标存取
	http.HandleFunc("/assets/img/cussvg", CusSvg)    //get上传的自定义svg
	http.HandleFunc("/assets/updateCus", UpdateCus)  //更新自定义信息
	http.HandleFunc("/workspace", WorkSpace)         //保存工作区
	http.HandleFunc("/workspace/findname", FindName) //保存工作区
	http.HandleFunc("/code", AutoCode)               //最大编号

	fs := http.FileServer(http.Dir("/usr/local/nginx/html/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs)) //开启assets文件夹服务
	println("goserver start… listen 9098")
	err := http.ListenAndServe(":9098", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}