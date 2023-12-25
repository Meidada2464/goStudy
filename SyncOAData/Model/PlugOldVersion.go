package Model

type PlugOldVersion struct {
	Id             int32  `json:"id,omitempty" gorm:"primary_key"`
	PlugRegisterId int32  `json:"plug_register_id,omitempty"`
	PlugName       string `json:"plug_name,omitempty"`
	PlugVersion    string `json:"plug_version,omitempty"`
	PlugParameter  string `json:"plug_parameter,omitempty"`
	Introduce      string `json:"introduce,omitempty"`
	CreatedUser    string `json:"created_user,omitempty"`
	UpdatedUser    string `json:"updated_user,omitempty"`
}

func (PlugOldVersion) TableName() string {
	return "plug_old_version"
}
