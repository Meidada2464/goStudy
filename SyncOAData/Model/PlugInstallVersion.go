package Model

type PlugInstallVersion struct {
	Id            int32  `json:"id,omitempty" gorm:"primary_key"`
	PlugInstallId int32  `json:"plug_install_id,omitempty"`
	Types         string `json:"types,omitempty"`
	Name          string `json:"name,omitempty"`
	Version       string `json:"version,omitempty"`
	Parameter     string `json:"parameter,omitempty"`
}

func (PlugInstallVersion) TableName() string {
	return "plug_install_version"
}
