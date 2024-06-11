package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// 服务信息结构
type ServiceInfo struct {
	ID        string                 `json:"_id"`
	URL       string                 `json:"url"`
	PublicURL string                 `json:"publicUrl"`
	HostInfo  map[string]interface{} `json:"hostInfo"`
}

// 获取环境变量中的注册服务URL
func GetCentralServer() string {
	centralServer := os.Getenv("CENTRAL_SERVER")
	if centralServer == "" {
		centralServer = "http://localhost:8000" // 默认值
	}
	return centralServer
}

// 获取环境变量中的本后端服务URL
func GetGoAPIURL() string {
	goAPIURL := os.Getenv("API_URL")
	if goAPIURL == "" {
		goAPIURL = "http://localhost:4000" // 默认值
	}
	return goAPIURL
}

// removePort 函数移除URL中的端口号, 例如 http://localhost:4000 -> http://localhost
func removePort(rawURL string) string {
	if strings.Contains(rawURL, ":") {
		return rawURL[:strings.LastIndex(rawURL, ":")]
	}
	return rawURL
}

// 全局变量
var (
	URL            = GetGoAPIURL()
	CENTRAL_SERVER = GetCentralServer()
	id             = "GO Server: " // 替换为你的服务ID
	hostInfo, _    = GetHostInfo()
)

func RegisterService(port string) bool {
	// 去掉端口号，然后添加新的端口号
	URL = removePort(URL)

	// 通过指针修改全局变量
	URL = URL + ":" + port
	id = id + URL

	// 创建服务信息
	serviceInfo := ServiceInfo{
		ID:        id,
		URL:       URL,
		PublicURL: URL,
		HostInfo:  hostInfo,
	}

	// 将服务信息转换为json格式
	jsonData, err := json.Marshal(serviceInfo) // 将服务信息转换为json格式
	if err != nil {
		log.Fatal(err)
	}

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/backend/register-service", CENTRAL_SERVER), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to register service: %s\n", err.Error())
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	} else {
		fmt.Printf("Failed to register service: %s\n", resp.Status)
		return false
	}
}

// 发送心跳功能
func SendHeartbeat() bool {
	// 创建服务信息
	serviceInfo := ServiceInfo{
		ID:        id,
		URL:       URL,
		PublicURL: URL,
		HostInfo:  hostInfo,
	}

	// 将服务信息转换为json格式
	jsonData, err := json.Marshal(serviceInfo) // 将服务信息转换为json格式
	if err != nil {
		log.Fatal(err)
	}

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/backend/register-service", CENTRAL_SERVER), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to register service: %s\n", err.Error())
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	} else {
		fmt.Printf("Failed to register service: %s\n", resp.Status)
		return false
	}
}

// 注销请求结构
type UnregisterRequest struct {
	ID string `json:"_id"`
}

// 注销服务功能
func UnregisterService() {
	unregisterReq := UnregisterRequest{ID: id}

	// 将注销请求转换为json格式
	jsonData, err := json.Marshal(unregisterReq)
	if err != nil {
		fmt.Printf("Failed to marshal unregister request: %s\n", err.Error())
		return
	}

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/backend/unregister-service", CENTRAL_SERVER), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to unregister service: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Service unregistered")
	} else {
		fmt.Printf("Failed to unregister service: %s\n", resp.Status)
	}
}
