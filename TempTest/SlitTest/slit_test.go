package SlitTest

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"goStudy/imroc/req"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Router struct {
	middlewareChain []string
}

type (
	ProvRes struct {
		Code      int          `json:"code,omitempty"`
		Data      []ProResData `json:"data"`
		Msg       string       `json:"msg,omitempty"`
		RequestId string       `json:"request_id,omitempty"`
	}

	ProResData struct {
		ID        int64  `json:"id"`
		CountryID int    `json:"country_id"`
		RegionID  int    `json:"region_id"`
		Name      string `json:"name"`
		Ename     string `json:"ename"`
		Code      string `json:"code"`
		Status    int    `json:"status"`
		Remark    string `json:"remark"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

func TestUncode(t *testing.T) {
	//str := `{"code": 0,"request_id": "1701748710.5442_782933","data": [{"id": 1,"country_id": 1,"region_id": 4,"name": "北京","ename": "beijing","code": "beijing","status": 1,"remark": "beijing","created_at": "2020-02-17 14:57:32","updated_at": "2020-02-17 14:57:32"}]}`
	//
	//var result ProvRes
	//
	//err := json.Unmarshal([]byte(str), &result)
	//if err != nil {
	//	fmt.Println("Error decoding JSON:", err)
	//	return
	//}

	var res ProvRes
	provinceIdMap := make(map[int64]string, 10)
	resp, err := req.Get("https://service-taishan.bs58i.baishancloud.com/api/location/province")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	statusCode := resp.Response().StatusCode
	if statusCode == 200 {
		err = resp.ToJSON(&res)
		if err != nil {
			fmt.Printf("%+v\n", err)

		}
	}

	if statusCode != 200 {
		fmt.Printf("%+v\n", err)

		return
	}

	for _, d := range res.Data {
		provinceIdMap[d.ID] = d.Name
	}

	fmt.Printf("%+v\n", provinceIdMap)

}

func TestMapNull(t *testing.T) {
	ipInfo := MapT()
	ipInfoStr, _ := json.Marshal(ipInfo)
	fmt.Println("aaa", string(ipInfoStr))
}

func MapT() (aInfo map[string]int) {
	aInfo = map[string]int{}
	return
}

type (
	PingParams struct {
		DetectFrequency int32  `json:"detect_frequency"` // 探测频率 单位 s/次
		PkgSize         int    `json:"package_size"`     // 包大小，byte
		PkgSizeStr      string `json:"-"`                // 包大小，byte
		PkgNum          int    `json:"package_num"`      // 包数量
		PkgNumStr       string `json:"-"`                // 包数量
		IsNeedMtr       bool   `json:"is_need_mtr"`      // 是否需要mtr探测
		Threshold       int32  `json:"threshold"`        // mtr探测触发阈值
		MaxHop          int32  `json:"max_hop"`          //

		RealDetectFrequency time.Duration `json:"-"` // 探测频率
		Wait                time.Duration `json:"-"` // 包间隔等待时间
		Timeout             time.Duration `json:"-"` // 探测超时时间
	}

	Host struct {
		ID               uint            `gorm:"primaryKey" json:"id,omitempty"`
		CreatedAt        time.Time       `gorm:"column:created_at" json:"createdAt,omitempty"`
		UpdatedAt        time.Time       `gorm:"column:updated_at" json:"updatedAt,omitempty"`
		DeletedAt        *gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deletedAt,omitempty"`
		Endpoint         string          `gorm:"column:endpoint" json:"ep"`            // 主机名，可能非标或者ip
		IP               string          `gorm:"column:ip" json:"ip"`                  // ip拼接，不完整列表, TODO: del
		IsInternalAsset  bool            `gorm:"column:is_internal_asset" json:"iia"`  // 是否内部资产，取决于是否录入在ITIL
		IsAgentInstalled bool            `gorm:"column:is_agent_installed" json:"iai"` // 是否安装探测agent
		AgentVersion     string          `gorm:"column:agent_version" json:"av"`       // agent版本
		ServiceType      string          `gorm:"column:service_type" json:"st"`        // 服务类型，多个用,拼接
		Tags             string          `gorm:"column:tags;" json:"tags"`             // 标签
		ServiceTypes     []string        `gorm:"-" json:"-"`                           // 资源类型列表
		UseStatus        int             `gorm:"column:use_status" json:"us"`          // 使用状态
		HangStatus       int             `gorm:"column:hang_status" json:"hs"`         // 挂起状态
		FaultStatus      int             `gorm:"column:fault_status" json:"fs"`        // 故障状态
		ISP              string          `gorm:"column:isp" json:"isp"`                // 运营商，多个用,拼接
		ISPs             []string        `gorm:"-" json:"-"`                           // 运营商列表
		Country          string          `gorm:"column:country" json:"co"`             // 国家
		Province         string          `gorm:"column:province" json:"prov"`          // 省份
		City             string          `gorm:"column:city" json:"city"`              // 城市
		Node             string          `gorm:"column:node" json:"node"`              // 节点
		CacheGroup       string          `gorm:"column:cache_group" json:"cg"`         // cache组
		Region           string          `gorm:"column:region" json:"rg"`              // 大区
		Status           bool            `gorm:"column:status" json:"status"`          // 内部状态，目前没用，后续可以根据最后存活时间或者tcp是否还在建联来做一定的判断
		NetType          int32           `gorm:"column:net_type" json:"nt"`            // 网络类型
		LiveAt           *time.Time      `gorm:"column:live_at" json:"la"`             // 最后存活时间

		// 以下两个字段是用来记录小节点（a,b类）设备账号粒度的带宽，用于雾探测选点时使用
		AccountIpv4Bandwidth float64 `gorm:"account_ipv4_bandwidth" json:"av4bw"` // 账号是ipv4的最大带宽
		AccountIpv6Bandwidth float64 `gorm:"account_ipv6_bandwidth" json:"av6bw"` // 账号是ipv6的最大带宽

		SupportIPv4 bool `gorm:"-" json:"isv4"` // 是否支持ipv4
		SupportIPv6 bool `gorm:"-" json:"isv6"` // 是否支持ipv6
		IsIntranet  bool `gorm:"-" json:"isin"` // 是否是内网设备(全部是内网的IP)
	}
)

func TestMasual(t *testing.T) {
	str := "{\"detect_frequency\":10,\"package_size\":10,\"package_num\":2,\"is_need_mtr\":true,\"threshold\":20,\"max_hop\":7,\"is_use_tcp_ping\":true,\"tcp_ping_port\":80}"

	var Ping PingParams

	err := json.Unmarshal([]byte(str), &Ping)
	if err != nil {
		return
	}

	a := gjson.Get(str, "detect_frequency").Int()
	fmt.Println("a", a)
}

func TestGorm(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/fftest?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	host := Host{
		ID:           1,
		NetType:      0,
		Endpoint:     "Biophilia.local==aaa",
		AgentVersion: "2.0.4",
	}
	err = db.Table("hosts_tests").Updates(&host).Error
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	err = db.Table("hosts_tests").Model(&host).Select("net_type").Updates(Host{NetType: host.NetType}).Error
	if err != nil {
		fmt.Println("err2:", err)
		return
	}
}

func TestMd5Key(t *testing.T) {
	AddHeaderToken(map[string]string{}, "H5hPC0=lLB7xy,*Kr8=")
	unix := time.Now().Add(-5 * time.Minute).Unix()

	fmt.Println(unix)
}

func AddHeaderToken(headers map[string]string, key string) {
	if headers == nil {
		headers = make(map[string]string)
	}
	nowstamp := time.Now().Unix()
	headers["timestamp"] = strconv.FormatInt(nowstamp, 10)
	token := Md5Encrypt(strconv.FormatInt(nowstamp, 10), key)
	if token != "" {
		bearerToken := "Bearer " + token
		headers["Authorization"] = bearerToken
	}

	fmt.Println(headers)
}

func Md5Encrypt(enc string, key string) string {
	md5 := md5.New()
	md5.Write([]byte(enc + key))
	md5Str := hex.EncodeToString(md5.Sum(nil))
	return md5Str
}

// 获取探测平台的签名
func TestGetSignature(t *testing.T) {
	clientId := "1"
	headers := make(map[string]string)
	unix := time.Now().Unix()
	formatInt := strconv.FormatInt(unix, 10)
	message := strings.Join([]string{"7v0pES5mfq4BnInSZoaylba28D1b6wlY", formatInt}, "\n")
	signature := generateSignature("7v0pES5mfq4BnInSZoaylba28D1b6wlY", message)

	headers["Authorization"] = signature
	headers["ClientId"] = clientId
	headers["Timestamp"] = formatInt

	fmt.Println(headers)
}

func generateSignature(secret, message string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func TestSomething(t *testing.T) {
	fmt.Println("TestSomething")
	assert.New(t)
}
