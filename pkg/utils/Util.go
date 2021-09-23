package utils

import (
	"bytes"
	"hash/crc32"
	"net"
	"strconv"
	"strings"
)

func Hash(s string, len uint32) uint {
	v := crc32.ChecksumIEEE([]byte(s))
	var num uint32
	if v >= 0 {
		num = v % len
	}
	if -v >= 0 {
		num = -v % len
	}
	return uint(num)
}

func GetAllNodeAddrs() map[int]string {
	addrs := make(map[int]string)
	return addrs
}

//设置分片规则
func RuleTransAddr(key string) string {
	//rpc启用3台
	num := uint32(3)
	k := int(Hash(key, num))
	v := GetAllNodeAddrs()[k]
	return v
}

func StringIpToInt(ip string) int {
	ipSegs := strings.Split(ip, ".")
	var ipInt int = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func IpIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

func GetLocalIp() string {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var ip net.IP
	for _, addr := range ips {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}
		return ip.String()
	}
	return ""
}

func SliceToMap(arr []string) map[string]bool {
	data := make(map[string]bool)
	for _, v := range arr {
		data[v] = true
	}
	return data
}

func FilterRepeat(arr []string) []string {
	data := make(map[string]bool)
	res := make([]string, 0)
	for _, v := range arr {
		if data[v] {
			continue
		}
		res = append(res, v)
	}
	return res
}
