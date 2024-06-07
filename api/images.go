package api

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		if len(image.RepoTags) > 0 {
			imageStrs = append(imageStrs, image.RepoTags[0])
		}
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

	// 如果是GET请求，获取一个docker image 的详细信息
	if method == http.MethodGet {
		ViewImage(w, r)
	}

	// 如果是DELETE请求，删除一个docker image
	if method == http.MethodDelete {
		ImageDeleter(w, r)
	}
}

// 查看一个docker image 的详细信息
func ViewImage(w http.ResponseWriter, r *http.Request) {
	Cors(w)

	// ImageDetails 定义了我们想要返回的镜像详细信息的结构
	type ImageDetails struct {
		WorkingDir     string   `json:"WorkingDir"`
		Entrypoint     []string `json:"Entrypoint"`
		Cmd            []string `json:"Cmd"`
		Id             string   `json:"Id"`
		Created        string   `json:"Created"`
		Size           string   `json:"Size"`
		Architecture   string   `json:"Architecture"`
		RepositoryTags []string `json:"RepositoryTags"`
		Os             string   `json:"Os"`
		DockerVersion  string   `json:"DockerVersion"`
	}

	// 获取参数
	params := strings.Split(r.URL.Path, "/")
	if len(params) < 4 {
		log.Fatalf("Error: no image name")
	}

	// 获取docker image 的名称
	var imageName string
	if params[len(params)-2] == "images" {
		imageName = params[len(params)-1]
	} else {
		imageName = params[len(params)-2] + "/" + params[len(params)-1]
	}

	// fmt.Println("Inspect Docker image: ", imageName)

	// 获取一个docker image 的详细信息
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

	// 格式化详细信息
	formattedDetails := ImageDetails{
		WorkingDir:     image.Config.WorkingDir,
		Entrypoint:     image.Config.Entrypoint,
		Cmd:            image.Config.Cmd,
		Id:             image.ID,
		Created:        image.Created,
		Size:           fmt.Sprintf("%.2f MB", float64(image.Size)/1024/1024),
		Architecture:   image.Architecture,
		RepositoryTags: image.RepoTags,
		Os:             image.Os,
		DockerVersion:  image.DockerVersion,
	}

	// fmt.Println("Inspect Docker image: ", formattedDetails)

	// 返回docker image 的详细信息
	json.NewEncoder(w).Encode([]ImageDetails{formattedDetails}) // 返回一个数组，方便前端处理

}

// 删除一个docker image
func ImageDeleter(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	params := strings.Split(r.URL.Path, "/")
	if len(params) < 4 {
		log.Fatalf("Error: no image name")
	}

	// 获取docker image 的名称
	var imageName string
	if params[len(params)-2] == "images" {
		imageName = params[len(params)-1]
	} else {
		imageName = params[len(params)-2] + "/" + params[len(params)-1]
	}

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
	fmt.Println("Delete success: ", imageName)
	json.NewEncoder(w).Encode("Delete success: " + imageName)
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

	fmt.Println("Pulling Docker image: ", imageName)

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Docker client: %v", err), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	// 拉取Docker镜像，打印输出流
	out, err := cli.ImagePull(r.Context(), imageName, image.PullOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error pulling Docker image: %v", err), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 实时读取输出流并打印到控制台
	var output strings.Builder
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // 实时打印到控制台
		output.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error reading Docker image pull response: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回成功
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Pull success: " + imageName)
}
