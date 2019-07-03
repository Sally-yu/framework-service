package model


//grafana dashboard
type Dashboard struct {
	Id int
	Version int
	Slug string
	Title string
	Data  string
	Org_id int
	Created string
	Updated string
	Updated_by int
	Created_by int
	Gnet_id int 
	Plugin_id string
	Folder_id int
	Is_folder int
	Has_acl int
	Uid string
}

func (dash *Dashboard)Find()  {
	
}