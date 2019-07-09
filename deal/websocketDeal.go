package deal

import (
	"encoding/json"
	"framework-service/model"
	"time"
)


type Data struct {
	Keys []string      `json:"keys"`
	Time time.Duration `json:"time"`
}

func DeviceValueWS(conn *Connection,d []byte) { //实际是int64
	var (
		err error
	)
	data := Data{}
	json.Unmarshal(d, &data)
	for {
		res:=[]model.Res{}
		for k := range data.Keys {
			r := GetAttValue(data.Keys[k])
			res = append(res, r)
		}
		flow,_:=json.Marshal(res)
		if err = conn.WriteMessage(flow); err != nil {
			return
		}
		time.Sleep(data.Time * time.Second)

	}
}

func NotifyListWS(conn *Connection){
	noti := model.Notif{}
	for {
		_, res := noti.AllNotif()
		flow, _ := json.Marshal(res)
		if err := conn.WriteMessage(flow); err != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
}