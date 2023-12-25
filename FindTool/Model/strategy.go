package Model

type Strategy struct {
	Id                int32  `json:"id" gorm:"primary_key"`
	Metric            string `json:"metric"`
	Tags              string `json:"tags"`
	MaxStep           int32  `json:"max_step"`
	Priority          int32  `json:"priority"`
	Func              string `json:"func"`
	Op                string `json:"op"`
	RightValue        string `json:"right_value"`
	Note              string `json:"note"`
	RunBegin          string `json:"run_begin"`
	RunEnd            string `json:"run_end"`
	TemplateId        int32  `json:"template_id"`
	FieldTransform    string `json:"field_transform"`
	Description       string `json:"description"`
	HandleDescription string `json:"handle_description"`
	Owner             string `json:"owner"`
	Notify            string `json:"notify"`
	Stage             int32  `json:"stage"`
	NoData            int32  `json:"no_data"`
	Step              int32  `json:"step"`
	Status            int32  `json:"status"`
	RecoverNotify     int32  `json:"recover_notify"`
	SilencesTime      int32  `json:"silences_time"`
	MarkTags          string `json:"mark_tags"`
	WebLink           string `json:"web_link"`
	CoverDrives       string `json:"cover_drives"`
	MasterType        int32  `json:"master_type"`
	UnRecoveredTime   int32  `json:"UnRecoveredTime" gorm:"column:un_recovered_time"`
	CreatedUser       string `json:"created_user"`
	LastUpdatedUser   string `json:"last_updated_user"`

	GT *GroupTemplate `gorm:"ForeignKey:TemplateId;AssociationForeignKey:TemplateId"`
}

func (Strategy) TableName() string {
	return "strategy"
}
