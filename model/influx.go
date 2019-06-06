package model

type Influx struct {
	DbName       string   `json:"dbname" form:"dbname" bson:"dbname"`
	Address      string   `json:"address" form:"address" bson:"address"`
	Username     string   `json:"username" form:"username" bson:"username"`
	Password     string   `json:"password" form:"password" bson:"password"`
	Measurements []string `json:"measurements" form:"measurements" bson:"measurements"`
}
