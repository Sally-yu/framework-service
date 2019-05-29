package model

//树子节点
type ChildNode struct {
	Title  string `json:"title" bson:"key"`
	Key    string `json:"key" bson:"key"`
	Author string `json:"author" bson:"author"`
	IsLeaf bool   `json:"isleaf" bson:"isleaf"`
	Fav    bool   `json:"fav" bson:"fav"`
	Share  bool   `json:"share" bson:"share"`
}
