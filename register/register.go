package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 服务信息结构
type ServiceInfo struct {
	ID        string                 `json:"_id"`
	URL       string                 `json:"url"`
	PublicURL string                 `json:"publicUrl"`
	HostInfo  map[string]interface{} `json:"hostInfo"`
}

// 全局变量
var (
	URL            = "http://192.168.0.103:4000" // 替换为你的服务地址
	CENTRAL_SERVER = "http://localhost:8000"     // 替换为你的中央服务器地址
	id             = "GO Server: " + URL         // 替换为你的服务ID
)

func RegisterService() bool {
	hostInfo, err := GetHostInfo()
	if err != nil {
		log.Fatal(err)
	}

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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/register-service", CENTRAL_SERVER), bytes.NewBuffer(jsonData))
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

// 心跳请求结构
type HeartbeatRequest struct {
	ID       string                 `json:"_id"`
	HostInfo map[string]interface{} `json:"hostInfo"`
}

// 发送心跳功能
func SendHeartbeat() bool {
	hostInfo, _ := GetHostInfo()
	heartbeatReq := HeartbeatRequest{ID: id, HostInfo: hostInfo}

	// 将心跳请求转换为json格式
	jsonData, err := json.Marshal(heartbeatReq)
	if err != nil {
		fmt.Printf("Failed to marshal heartbeat request: %s\n", err.Error())
		return false
	}

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/service-heartbeat", CENTRAL_SERVER), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send heartbeat: %s\n", err.Error())
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码
	return resp.StatusCode == http.StatusOK
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
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/unregister-service", CENTRAL_SERVER), bytes.NewBuffer(jsonData))
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
