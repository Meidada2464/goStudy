package Model

type PlugRegister struct {
	Id            int32  `json:"id,omitempty" gorm:"primary_key"`
	PlugName      string `json:"plug_name,omitempty"`
	PlugVersion   string `json:"plug_version,omitempty"`
	PlugType      string `json:"plug_type,omitempty"`
	Introduce     string `json:"introduce,omitempty"`
	Configuration string `json:"configuration,omitempty"`
	CreatedUser   string `json:"created_user,omitempty"`
	UpdatedUser   string `json:"updated_user,omitempty"`

	PlugOldVersion []PlugOldVersion `gorm:"ForeignKey:PlugRegisterId;AssociationForeignKey:Id"`
}

func (PlugRegister) TableName() string {
	return "plug_register"
}
