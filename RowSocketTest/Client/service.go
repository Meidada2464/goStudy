/**
 * Package Client
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/9 13:21
 */

package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

func main() {
	conn, err := net.ListenPacket("ip4:udp", "127.0.0.1")
	if err != nil {
		panic(err)
	}
	data, err := EncodeUDPPacket(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), 8972, 0, []byte("hello"))
	if err != nil {
		log.Printf("failed to EncodePacket: %v", err)
		return
	}
	remoteAddr := &net.IPAddr{IP: net.ParseIP("127.0.0.1")}
	if _, err := conn.WriteTo(data, remoteAddr); err != nil {
		log.Printf("failed to write packet: %v", err)
		conn.Close()
		return
	}

	log.Println("send ok,data : ", data)
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		log.Fatal(err)
	}
	buffer = buffer[:n]
	packet := gopacket.NewPacket(buffer, layers.LayerTypeUDP, gopacket.NoCopy)
	// Get the UDP layer from this packet
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		if app := packet.ApplicationLayer(); app != nil {
			fmt.Printf("reply: %s\n", app.Payload())
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
