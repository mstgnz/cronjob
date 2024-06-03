package models

type Group struct {
	ID   int    `json:"id"`
	UID  int    `json:"uid"`
	Name string `json:"name"`
}

func (g *Group) GetLogs(offset, limit int) []*Group {
	Groups := []*Group{}
	return Groups
}

func (g *Group) GetLog(id int) []*Group {
	Groups := []*Group{}
	return Groups
}

func (g Group) CreateLog(id int, is_error bool, log string) Group {
	return g
}
