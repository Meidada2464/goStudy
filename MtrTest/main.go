package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	// 创建一个日志记录器
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// 示例日志输出
	logger.Println("程序启动")

	// 获取命令行参数
	args := os.Args
	if len(args) < 6 {
		return
	}
	ip := args[1]

	// 获取第二个参数
	maxHops := args[2]
	maxHopsInt, err := strconv.Atoi(maxHops)
	if err != nil {
		logger.Println("error:", err)
		return
	}

	timeOutMs := args[3]
	timeOutMsInt, err := strconv.Atoi(timeOutMs)
	if err != nil {
		logger.Println("error:", err)
		return
	}

	packetSize := args[4]
	packetSizeInt, err := strconv.Atoi(packetSize)
	if err != nil {
		logger.Println("error:", err)
		return
	}

	sntSize := args[5]
	sntSizeInt, err := strconv.Atoi(sntSize)
	if err != nil {
		logger.Println("error:", err)
		return
	}

	sntLimit := args[6]
	sntLimitInt, err := strconv.Atoi(sntLimit)
	if err != nil {
		logger.Println("error:", err)
		return
	}

	params := &Params{
		MaxHops:    maxHopsInt,
		TimeoutMs:  timeOutMsInt,
		PacketSize: packetSizeInt,
		SntSize:    sntSizeInt,
		SntLimit:   sntLimitInt,
	}

	// 获取配置参数
	log.Println("开始执行命令,需要等待约", timeOutMsInt, "秒。。。")
	cmdRespContent, err := execCmd(ip, params, logger)
	if err != nil {
		logger.Println("execCmd error:", err)
	}

	// 解析参数
	log.Println("开始解析参数。。。")
	content := parseContent(cmdRespContent, logger)

	log.Println("开始写入数据。。。")
	wirteMessageTofile(content)
}

func execCmd(ip string, params *Params, log *log.Logger) ([]string, error) {
	ctx := context.Background()
	sntLimit := params.SntSize
	//最大跳数
	maxHop := params.MaxHops
	//重试次数
	sntLimit = params.SntLimit
	if sntLimit == 0 {
		sntLimit = 10
	}
	respContent := make([]string, 0)
	log.Println("获取到数据:", ip, "params:", *params)

	log.Println("执行命令:", "mtr", ip, "-nz", "-c", strconv.Itoa(sntLimit), "-m", strconv.Itoa(maxHop))
	//最大跳数量
	cmd := exec.CommandContext(ctx, "mtr", ip, "-nz", "-c", strconv.Itoa(sntLimit), "-m", strconv.Itoa(maxHop))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("exec stdout error:", err)
		return respContent, err
	}

	if err := cmd.Start(); err != nil {
		log.Println("exec start error:", err)
		return respContent, err
	}

	scanner := bufio.NewScanner(stdout)
	log.Println("exec bufio NewScanner:", "scanner.Scan():", scanner.Scan(), "scanner.Err():", scanner.Err())
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("exec scanner Scan:", "line:", line)
		respContent = append(respContent, line)
	}
	if err := cmd.Wait(); err != nil {
		log.Println("exec wait error:", err, "data:", respContent)
		return respContent, err
	}
	//执行超时
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("exec timeout error:", ctx.Err())
		return []string{}, ctx.Err()
	}
	return respContent, nil
}

// 解析参数
func parseContent(content []string, log *log.Logger) []*DetectMtrHopDetail {
	//无用数据过滤
	log.Println("开始过滤无用数据。。。", "content:", content)
	content = filterUselessContent(content)

	log.Println("数据过滤完毕。", "content:", content)
	hopDetail := make([]*DetectMtrHopDetail, 0)
	for i, line := range content {
		//空格过滤
		log.Println("开始过滤空数据", "第", i, "行", "过滤前数据:", line)
		lineParts := filterWhiteSpace(strings.Fields(line))
		log.Println("数据过滤完毕", "第", i, "行", "过滤后数据:", lineParts)
		//数据量完整性过滤
		if len(lineParts) < lineDataMinCnt {
			log.Println("数据不完整，跳过", "第", i, "行", "过滤后数据:", len(lineParts))
			continue
		}
		//【特殊处理】2024-01-22 mtr探测记录存在AS和IP段粘连的情况，需要配合做兜底解析处理
		if len(lineParts) == lineDataMinCnt {
			log.Println("开始处理S和IP段粘连的情况", "数据:", lineParts)
			lineParts = handleAsIpSplit(lineParts)
			log.Println("S和IP段粘连的情况处理完毕", "数据:", lineParts)
		}
		log.Println("开始解析各个字段数据：AVG、LAST")
		hopDetail = append(hopDetail, parseHopDetail(lineParts))
	}
	return hopDetail
}

func wirteMessageTofile(mtrHopData []*DetectMtrHopDetail) {
	if len(mtrHopData) == 0 {
		fmt.Println("wirteMessageTofile: no data")
		return
	}
	// 遍历数据
	var metrics []mtrData
	for _, hotDataDetail := range mtrHopData {
		var metrData mtrData
		tags := make(map[string]string)
		fields := make(map[string]interface{})
		//字段填充
		tags["任务名"] = "test-mtr"
		tags["TTL"] = strconv.Itoa(hotDataDetail.TTL)
		tags["目标地址"] = hotDataDetail.Address
		tags["供应商标识"] = hotDataDetail.AsSign

		fields["丢包率"] = hotDataDetail.Loss
		fields["发包数"] = hotDataDetail.Snt
		fields["最后一个ICMP包往返时间"] = hotDataDetail.Last
		fields["所有ICMP包往返时间平均值"] = hotDataDetail.Avg
		fields["所有ICMP包中最快往返时间"] = hotDataDetail.Best
		fields["所有ICMP包中最慢往返时间"] = hotDataDetail.Wrst
		fields["延迟的标准偏差"] = hotDataDetail.StDev

		metrData.Tag = tags
		metrData.Fild = fields
		metrics = append(metrics, metrData)
	}

	// 将数据写入文件
	// 打开文件以写入数据
	filename := "/tmp/mtrData.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// 将每个 mtrData 对象序列化为 JSON 并逐行写入文件
	for _, item := range metrics {
		jsonData, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("Error marshaling data: %v\n", err)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		_, err = file.WriteString("\n")
		if err != nil {
			fmt.Printf("Error writing newline to file: %v\n", err)
			return
		}
	}

	fmt.Printf("Data successfully written to %s\n", filename)
}

// 数据截断，只保留mtr记录即可，其他信息去除
func filterUselessContent(contents []string) []string {
	filterContents := make([]string, 0)
	if len(contents) > 2 {
		filterContents = contents[2:]
	}
	return filterContents
}

// 空数据清理
func filterWhiteSpace(parts []string) []string {
	filterParts := make([]string, 0)
	for _, part := range parts {
		if part == "" {
			continue
		}
		filterParts = append(filterParts, part)
	}
	return filterParts
}

func handleAsIpSplit(lineParts []string) []string {
	newLinkParts := make([]string, 0)
	for index, content := range lineParts {
		//特殊处理
		if index == specialAsIndex {
			if len(content) < splitIndex {
				continue
			}
			//as信息
			asInfo := content[:splitIndex]
			//ip信息
			ipInfo := content[splitIndex:]
			newLinkParts = append(newLinkParts, asInfo)
			newLinkParts = append(newLinkParts, ipInfo)
			continue
		}
		newLinkParts = append(newLinkParts, content)
	}
	return newLinkParts
}

func parseHopDetail(lineParts []string) *DetectMtrHopDetail {
	hop := &DetectMtrHopDetail{
		TTL:     parseTTL(lineParts),
		AsSign:  parseAsSign(lineParts),
		Address: parseAddress(lineParts),
		Loss:    parseLoss(lineParts),
		Snt:     parseInt(lineParts[4]),
		Last:    parseFloat(lineParts[5]),
		Avg:     parseFloat(lineParts[6]),
		Best:    parseFloat(lineParts[7]),
		Wrst:    parseFloat(lineParts[8]),
		StDev:   parseFloat(lineParts[9]),
	}
	return hop
}

func parseTTL(parts []string) int {
	ttlPart := parts[0]
	ttlPart = strings.ReplaceAll(ttlPart, ".", "")
	ttlPart = strings.ReplaceAll(ttlPart, " ", "")
	ttl, err := strconv.Atoi(ttlPart)
	if err != nil {
		return 0
	}
	return ttl
}

// 解析供应商标识
func parseAsSign(parts []string) string {
	asSign := parts[1]
	return asSign
}

// 解析目标ip
func parseAddress(parts []string) string {
	address := parts[2]
	return address
}

// 解析丢包率
func parseLoss(parts []string) float64 {
	lossPartStr := parts[3]
	lossPartStr = strings.ReplaceAll(lossPartStr, "%", "")
	loss, err := strconv.ParseFloat(lossPartStr, 64)
	if err != nil {
		return 0
	}
	return loss
}

func parseInt(valStr string) int {
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0
	}
	return int(val)
}
func parseFloat(valStr string) float64 {
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0
	}
	return val

}
