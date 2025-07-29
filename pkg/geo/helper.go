package geo

import (
	"net"
)

func GetCityName(ip string, langOptional ...string) string {
	netIP := net.ParseIP(ip)

	if netIP.IsLoopback() {
		return "localhost"
	}

	if IsPrivateIP(netIP) {
		return "内网"
	}

	if record, err := G.City(netIP); err == nil {
		lang := "zh-CN"
		if len(langOptional) > 0 {
			lang = langOptional[0]
		}
		if location, exist := record.City.Names[lang]; exist {
			return location
		}
	}

	return ""
}

// 判断是否为私有 IPv4 地址
func IsPrivateIPv4(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return ip4[0] == 10 ||
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
		(ip4[0] == 192 && ip4[1] == 168)
}

// 判断是否为私有 IPv6 地址
func IsPrivateIPv6(ip net.IP) bool {
	ip6 := ip.To16()
	if ip6 == nil {
		return false
	}
	// 检查唯一本地地址（ULA）
	if ip6[0] == 0xfc && ip6[1]&0xfe == 0x00 {
		return true
	}
	// 检查链路本地地址
	if ip6[0] == 0xfe && ip6[1]&0xc0 == 0x80 {
		return true
	}
	return false
}

// 判断是否为私有 IP 地址（包括 IPv4 和 IPv6）
func IsPrivateIP(ip net.IP) bool {
	return ip.IsLoopback() || IsPrivateIPv4(ip) || IsPrivateIPv6(ip)
}
