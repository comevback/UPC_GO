package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// get /api/images 是获取所有docker images 的列表
// get /api/images/:imageName 是获取一个docker image 的详细信息
// post  /api/files/:filename 是上传一个zip文件，解压并利用 buildpack 创建一个docker image
// delete /api/images/:imageName 是删除一个docker image

// 列出所有docker images
func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	// 跨域请求
	Cors(w)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()

	// 获取所有docker images
	images, err := cli.ImageList(r.Context(), image.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing Docker images: %v", err)
	}

	// 将数据格式化为字符串
	var imageStrs []string
	for _, image := range images {
		imageStrs = append(imageStrs, image.RepoTags[0])
	}

	// 返回所有docker images
	json.NewEncoder(w).Encode(imageStrs)
}

// 获取一个docker image 的详细信息
func ImageProcessor(w http.ResponseWriter, r *http.Request) {
	// 跨域请求
	Cors(w)
	method := r.Method

	// 获取参数
	params := strings.Split(r.URL.Path, "/")
	if len(params) < 4 {
		log.Fatalf("Error: no image name")
	}

	// 获取docker image 的名称
	imageName := params[len(params)-2] + "/" + params[len(params)-1]

	// 如果是GET请求，获取一个docker image 的详细信息
	if method == http.MethodGet {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatalf("Error creating Docker client: %v", err)
		}
		defer cli.Close()

		// 获取docker image 的详细信息
		image, _, err := cli.ImageInspectWithRaw(r.Context(), imageName)
		if err != nil {
			log.Fatalf("Error inspecting Docker image: %v", err)
		}

		// 返回docker image 的详细信息
		json.NewEncoder(w).Encode(image)
	}

	// 如果是DELETE请求，删除一个docker image
	if method == http.MethodDelete {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatalf("Error creating Docker client: %v", err)
		}
		defer cli.Close()

		// 删除docker image
		_, err = cli.ImageRemove(r.Context(), imageName, image.RemoveOptions{})
		if err != nil {
			log.Fatalf("Error removing Docker image: %v", err)
		}

		// 返回删除成功
		json.NewEncoder(w).Encode("Delete success")
	}
}

// pull 一个docker image
func ImagePuller(w http.ResponseWriter, r *http.Request) {
	Cors(w)
	method := r.Method

	if method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取参数
	params := strings.Split(r.URL.Path, "/")
	if len(params) < 4 {
		http.Error(w, "Error: no image name", http.StatusBadRequest)
		return
	}

	// 获取docker image 的名称
	imageName := params[len(params)-2] + "/" + params[len(params)-1]

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Docker client: %v", err), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	// 拉取Docker镜像
	out, err := cli.ImagePull(r.Context(), imageName, image.PullOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error pulling Docker image: %v", err), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 读取输出流，确认拉取成功
	buf := new(strings.Builder)
	_, err = io.Copy(buf, out)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading Docker image pull response: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回成功
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Pull success: " + buf.String())
}
