package model

type Attribute struct {
	Key         string `json:"key" form:"key" bson:"key"`
	Name        string `json:"name" form:"name" bson:"name"`
	Code        string `json:"code" form:"code" bson:"code"`
	ValueType   string `json:"valuetype" form:"valuetype" bson:"valuetype" `
	Unit        string `json:"unit" form:"unit" bson:"unit"`
	Description string `json:"description" form:"description" bson:"description"`
	Sum         bool   `json:"sum" form:"sum" bson:"sum"`
}
