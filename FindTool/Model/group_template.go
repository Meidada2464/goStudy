package Model

type GroupTemplate struct {
	Id         int32  `json:"id" gorm:"primary_key"`
	GroupId    int32  `json:"group_id"`
	TemplateId int32  `json:"template_id"`
	BindUser   string `json:"bind_user"`

	GH *GroupHost `gorm:"ForeignKey:GroupId;AssociationForeignKey:GroupId"`
}

func (GroupTemplate) TableName() string {
	return "group_template"
}
