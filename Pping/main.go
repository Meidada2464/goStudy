/**
 * Package Pping
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/10 16:06
 */

package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var (
	liveInp   = pflag.StringP("interface", "i", "", "interface name")
	fName     = pflag.StringP("read", "r", "", "pcap captured file")
	filterOpt = pflag.StringP("filter", "f", "", "pcap filter applied to packets")
	not_tcp   = 0
	no_TS     = 0
	offTm     int64
)

func main() {
	var (
		sumInt          time.Duration
		filtLocal       bool
		timeToRun       time.Duration
		maxPackets      int
		machineReadable bool
		tsvalMaxAge     time.Duration
		flowMaxIdle     time.Duration
	)
	pflag.DurationVarP(&sumInt, "sumInt", "q", 10*time.Second, "interval to print summary reports to stderr")
	pflag.BoolVarP(&filtLocal, "showLocal", "l", false, "show RTTs through local host applications")
	pflag.DurationVarP(&timeToRun, "seconds", "s", 0*time.Second, "stop after capturing for <num> seconds")
	pflag.IntVarP(&maxPackets, "count", "c", 0, "stop after capturing <num> packets")
	pflag.BoolVarP(&machineReadable, "machine", "m", false, "machine readable output")
	pflag.DurationVarP(&tsvalMaxAge, "tsvalMaxAge", "M", 10*time.Second, "max age of an unmatched tsval")
	pflag.DurationVarP(&flowMaxIdle, "flowMaxIdle", "F", 300*time.Second, "flows idle longer than <num> are deleted")

	snif := getSnifFromPcap("eth0", 65535)
	err := snif.SetBPFFilter(*filterOpt)

	if err != nil {
		fmt.Println("SetBPFFilter: ", err)
		return
	}

	// 将原始数据转化成以太网包格式的数据
	src := gopacket.NewPacketSource(snif, layers.LayerTypeEthernet)

	packets := src.Packets()

	for packet := range packets {
		// 处理包
		processPacket(packet)

	}

}

// processPacket 处理获取到的原始的包数据，获取到src和dst，以及包中的时间戳
func processPacket(pkt gopacket.Packet) {
	// 将pkg解析成TCP包结构
	tcpLayer := pkt.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		not_tcp++
		return
	}
	tcp, _ := tcpLayer.(*layers.TCP)
	tsval, tsecr := getTSFromTCPOpts(tcp)
	if tsval == 0 || (tsecr == 0 && !tcp.SYN) {
		no_TS++
		return
	}

	netlayer := pkt.Layer(layers.LayerTypeIPv6)
	if netlayer == nil {
		fmt.Println("no ipv6")
		netlayer = pkt.Layer(layers.LayerTypeIPv4)
		if netlayer == nil {
			fmt.Println("no ipv4")
			return
		}
	}

	// 获取srcIp和dstIp

	var ipsStr, ipdStr string
	if ip, ok := netlayer.(*layers.IPv4); ok {
		ipsStr = ip.SrcIP.String()
		ipdStr = ip.DstIP.String()
		fmt.Println("ipv4 src", ipsStr, "dst", ipdStr)
	} else {
		ip := netlayer.(*layers.IPv6)
		ipsStr = ip.SrcIP.String()
		ipdStr = ip.DstIP.String()
		fmt.Println("ipv6 src", ipsStr, "dst", ipdStr)
	}

	srcStr := ipsStr + ":" + fmt.Sprint(tcp.SrcPort)
	dstStr := ipdStr + ":" + fmt.Sprint(tcp.DstPort)
	fmt.Println("tcp src", srcStr, "dst", dstStr)
}

// getTSFromTCPOpts 用于从 TCP 选项中获取时间戳信息
func getTSFromTCPOpts(tcp *layers.TCP) (uint32, uint32) {
	var tsval, tsecr uint32
	options := tcp.Options
	for _, option := range options {
		if option.OptionType == layers.TCPOptionKindTimestamps && option.OptionLength == 10 { // Timestamp 选项长度为 10 字节
			tsval = binary.BigEndian.Uint32(option.OptionData[:4])
			tsecr = binary.BigEndian.Uint32(option.OptionData[4:8])
			break
		}
	}
	// 从TCPOptions中获取时间戳
	return tsval, tsecr
}

func getSnifFromPcap(nicName string, snapLen int) *pcap.Handle {
	// 创建一个新的非活动 pcap 句柄
	inactive, err := pcap.NewInactiveHandle(nicName)
	defer inactive.CleanUp()
	if err != nil {
		fmt.Println("pcap.NewInactiveHandle: ", err)
		return nil
	}

	err = inactive.SetSnapLen(snapLen)
	if err != nil {
		fmt.Println("SetSnapLen: ", err)
		return nil
	}

	snif, err := inactive.Activate()
	if err != nil {
		fmt.Println("Activate: ", err)
		return nil
	}

	return snif
}

func getSnifFromFile(fileName string) *pcap.Handle {
	openFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("os.Open: ", err)
		return nil
	}

	snif, err := pcap.OpenOfflineFile(openFile)
	if err != nil {
		fmt.Println("pcap.OpenOfflineFile: ", err)
		return nil
	}

	return snif
}
