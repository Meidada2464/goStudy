package Model

type GroupHost struct {
	Id      int32 `json:"id" gorm:"primary_key"`
	HostId  int64 `json:"host_id"`
	GroupId int32 `json:"group_id"`

	Host *Host `gorm:"ForeignKey:Id;AssociationForeignKey:HostId"`
}

func (GroupHost) TableName() string {
	return "group_host"
}
