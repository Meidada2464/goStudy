/**
 * Package httpC
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/3 14:22
 */

package httpC

import "time"

const (
	HTTP_FAILED              = -1
	HTTP_CODE_502            = 502
	HTTP_CODE_503            = 503
	MAX_TIMEOUT              = 10
	MAX_CONNECT_TIMEOUT      = 10 * time.Second
	DEFAULT_DETECT_FREQUENCY = time.Minute
	DEFAULT_HTTP_VERSION     = "2.0"
)

type (
	httpTaskParam struct {
		followRedirects   bool              // 是否跟随重定向，类似curl -L参数
		skipVerify        bool              // 是否跳过证书验证，类似curl -k参数
		enableHttp2       bool              // 是否使用http2,目前只支持http/1.1、http2
		maxRedirects      *int              // 最大重定向次数
		maxTimeout        time.Duration     // 最大超时时间，类似curl --max-time参数
		connectTimeout    time.Duration     // 连接超时时间，类似curl --connect-timeout参数
		HttpMethod        string            // 当前只支持post/get
		url               string            // 请求url
		srcIp             string            // 探测绑定ip
		basicAuth         *string           // 基础认证user:pass，类似curl -u参数
		header            map[string]string // 附加请求头，类似curl -H参数
		payload           []byte            // 请求体，类似curl -d参数
		expectedHttpCodes []int             // 期望的http状态码，下载速度之类在状态码不一致时不计算，可为空
		httpVersion       HttpVersion       // http协议版本
	}

	httpDlInfo struct {
		status            int
		fileSize          int64
		dstIp             string
		connTime          float64
		dnsTime           float64
		tlsTime           float64
		firstByteTime     float64
		downloadTime      float64
		totalDownloadTime float64
		downloadRate      float64
	}

	HttpVersion string
)

var (
	HttpVersion1dot0 HttpVersion = "1.0"
	HttpVersion1dot1 HttpVersion = "1.1"
	HttpVersion2dot0 HttpVersion = "2.0"
	HttpVersion3dot0 HttpVersion = "3.0"
)
