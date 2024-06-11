package register

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

// 得到本机IP
func GetLocalIP() string {
	// 获取本机IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	// 获取第一个非回环地址
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// bytesToGB 将字节转换为GB
func BytesToGB(b uint64) float64 {
	return float64(b) / (1024 * 1024 * 1024)
}

// formatUptime 将秒转换为可读的时间格式
func FormatUptime(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
}

// getUptime 获取系统运行时间（Linux）
func GetUptime() (uint64, error) {
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var uptimeSeconds float64
	_, err = fmt.Fscanf(file, "%f", &uptimeSeconds)
	if err != nil {
		return 0, err
	}

	return uint64(uptimeSeconds), nil
}

// getHostInfo 获取主机信息
func GetHostInfo() (map[string]interface{}, error) {
	// 获取CPU架构和数量
	architecture := runtime.GOARCH
	cpus := runtime.NumCPU()
	ip := GetLocalIP()

	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	// 获取操作系统平台和版本
	platform := runtime.GOOS
	release := runtime.Version()

	return map[string]interface{}{
		"hostname":     hostname,
		"ip":           ip,
		"architecture": architecture,
		"cpus":         cpus,
		"platform":     platform,
		"release":      release,
	}, nil
}
