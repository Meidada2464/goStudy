package Model

type Host struct {
	Id            int64     `json:"id"`
	Hostname      string    `json:"hostname"`
	Ip            string    `json:"ip"`
	AgentVersion  string    `json:"agent_version"`
	PluginVersion string    `json:"plugin_version"`
	MaintainBegin int64     `json:"maintain_begin"`
	MaintainEnd   int64     `json:"maintain_end"`
	LiveAt        int32     `json:"live_at"`
	RemoteIp      string    `json:"remote_ip"`
	Cachegroup    string    `json:"cachegroup"`
	Sertypes      string    `json:"sertypes"`
	GroupHost     GroupHost `gorm:"ForeignKey:HostId;AssociationForeignKey:Id"`
}

func (Host) TableName() string {
	return "host"
}
