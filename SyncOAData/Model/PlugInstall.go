package Model

type PlugInstall struct {
	Id          int32  `json:"id,omitempty" gorm:"primary_key"`
	TaskName    string `json:"task_name,omitempty"`
	TaskGroup   string `json:"task_group,omitempty"`
	PlugName    string `json:"plug_name,omitempty"`
	PlugVersion string `json:"plug_version,omitempty"`
	PlugType    string `json:"plug_type,omitempty"`
	Types       string `json:"types,omitempty"`
	Name        string `json:"name,omitempty"`
	Parameter   string `json:"parameter,omitempty"`
	Notes       string `json:"notes,omitempty"`
	CreatedUser string `json:"created_user,omitempty"`
	UpdatedUser string `json:"updated_user,omitempty"`

	PlugInstallVersionList []PlugInstallVersion `gorm:"ForeignKey:PlugInstallId;AssociationForeignKey:Id"`
}

func (PlugInstall) TableName() string {
	return "plug_install"
}
