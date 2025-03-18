package container

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

type Docker struct {
	Client  *client.Client
	AuthStr string `json:"authStr"`
}

type PullResponse struct {
	Status      string         `json:"status"`
	ErrorDetail map[string]any `json:"errorDetail"`
}

func newDocker() *Docker {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed connect docker")
	}
	authConfig := registry.AuthConfig{
		Username:      "admin",          // 替换为你的用户名
		Password:      "Ctx1ytxA@3zdj",  // 替换为你的密码
		ServerAddress: "reg.safedog.cn", // 替换为你的镜像仓库地址，包括端口号，例如 "registry.example.com:80" 或 "registry.example.com:443"
	}
	_, err = cli.RegistryLogin(context.TODO(), authConfig)
	if err != nil {
		log.Fatalf("Failed loggin docker")
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Fatalf("Error encoding auth config: %v", err)
	}
	return &Docker{
		Client:  cli,
		AuthStr: base64.URLEncoding.EncodeToString(encodedJSON),
	}
}

// Pull pulls an image from a specified repository using docker command.

func (d *Docker) Pull(imageName string) (err error) {
	out, err := d.Client.ImagePull(context.Background(), imageName, types.ImagePullOptions{
		RegistryAuth: d.AuthStr,
	})
	if err != nil {
		return err
	}
	defer out.Close()
	decoder := json.NewDecoder(out)
	for {
		var response PullResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				// 读取结束
				break
			}
			log.Printf("读取输出时出错: %v", err)
			return err
		}

		// 打印拉取镜像的名称（status 字段包含了镜像名称）
		if response.Status != "" {
			fmt.Println(response.Status)
		}
	}

	return err
}

// Tag tags an existing image with a new tag using docker command.
func (d *Docker) Tag(imageName string, newTag string) {
	if err := d.Client.ImageTag(context.Background(), imageName, newTag); err != nil {
		log.Fatalf("err: %v", err)
	}
}

// Push pushes a tagged image to a specified repository using docker command.
func (d *Docker) Push(imageName string) (err error) {
	out, err := d.Client.ImagePush(context.Background(), imageName, types.ImagePushOptions{
		RegistryAuth: d.AuthStr,
	})
	if err != nil {
		log.Fatalf("Failed push image %v,err: %v", imageName, err)
	}
	defer out.Close()
	decoder := json.NewDecoder(out)
	for {
		var response PullResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				// 读取结束
				break
			}
			log.Printf("读取输出时出错: %v", err)
			return err
		}

		// 打印拉取镜像的名称（status 字段包含了镜像名称）
		if response.ErrorDetail != nil {
			log.Printf("push输出时出错: %v", response.ErrorDetail)
		}
	}
	return err
}

// Push pushes a tagged image to a specified repository using docker command.
func (d *Docker) LoadAndPush(filePath string, repo string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.Background()
	response, err := d.Client.ImageLoad(ctx, f, true)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 获取所有镜像ID
	loadedImages := []string{}
	decoder := json.NewDecoder(response.Body)
	re := regexp.MustCompile(`Loaded image: ([^:]+):(.+)`)
	for {
		var jm jsonmessage.JSONMessage
		if err := decoder.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if jm.Stream != "" {
			matches := re.FindStringSubmatch(strings.TrimSpace(jm.Stream))
			if len(matches) > 2 {
				imageName := matches[1]
				tag := matches[2]
				repoTag := imageName + ":" + tag
				loadedImages = append(loadedImages, repoTag)
			}
		}
	}
	for _, v := range loadedImages {
		newStringList := strings.Split(v, "/")
		newString := fmt.Sprintf("%v/%v", repo, newStringList[len(newStringList)-1])
		d.Tag(v, newString)
		d.Push(newString)
	}

	return nil
}
