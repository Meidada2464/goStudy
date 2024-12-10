/**
 * Package Server
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/9 13:18
 */

package main

import (
	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

var (
	addr = flag.String("s", "localhost", "server address")
	port = flag.Int("p", 8972, "port")
)
var (
	stat         = make(map[string]int)
	lastStatTime = int64(0)
)

func main() {
	flag.Parse()
	conn, err := net.ListenPacket("ip4:udp", "127.0.0.1")
	if err != nil {
		panic(err)
	}
	cc := conn.(*net.IPConn)
	cc.SetReadBuffer(20 * 1024 * 1024)
	cc.SetWriteBuffer(20 * 1024 * 1024)
	handleConn(conn)
}
func handleConn(conn net.PacketConn) {
	for {
		buffer := make([]byte, 1024)
		n, remoteaddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("recv data : ", buffer, " from ", remoteaddr)
		buffer = buffer[:n]
		packet := gopacket.NewPacket(buffer, layers.LayerTypeUDP, gopacket.NoCopy)
		// Get the UDP layer from this packet
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			if app := packet.ApplicationLayer(); app != nil {
				data, err := EncodeUDPPacket(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), uint16(udp.SrcPort), uint16(udp.DstPort), app.Payload())
				if err != nil {
					log.Printf("failed to EncodePacket: %v", err)
					return
				}
				if _, err := conn.WriteTo(data, remoteaddr); err != nil {
					log.Printf("failed to write packet: %v", err)
					conn.Close()
					return
				}
			}
		}
	}
}

func EncodeUDPPacket(localIP, remoteIP net.IP, localPort, remotePort uint16, payload []byte) ([]byte, error) {
	ip := &layers.IPv4{
		Version:  4,
		TTL:      128,
		SrcIP:    localIP,
		DstIP:    remoteIP,
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(localPort),
		DstPort: layers.UDPPort(remotePort),
	}
	udp.SetNetworkLayerForChecksum(ip)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	err := gopacket.SerializeLayers(buf, opts, udp, gopacket.Payload(payload))
	return buf.Bytes(), err
}
