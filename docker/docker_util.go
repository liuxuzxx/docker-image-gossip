/*
 * @Description:Docker 操作相关的动作
 */
package docker

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

//
// 放置docker操作相关的命令
//

type DockerClient struct {
	cli       *client.Client
	imagePath string
}

func (d *DockerClient) init() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panic(err)
	}
	d.cli = cli
}

func (d *DockerClient) ListImage() {
	ctx := context.Background()
	images, err := d.cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	for _, image := range images {
		log.Printf("ImageName: %s, RepoTags: %v\n", image.ID, image.RepoTags)
	}
}

func (d *DockerClient) LoadImage(tarPath string) {
	file, err := os.Open(tarPath)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	response, err := d.cli.ImageLoad(ctx, file, true)
	if err != nil {
		log.Fatal(err)
	}
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(response.Body)
	defer response.Body.Close()
	logrus.Infoln("查看载入镜像的结果信息:", buffer.String())
	err = os.Remove(tarPath)
	if err != nil {
		logrus.Warnln("Remote file:", tarPath, "fail")
	} else {
		logrus.Infoln("Clear file:", tarPath, "success")
	}
}

func (d *DockerClient) SaveImage(name, tag string) {
	imageName := strings.Join([]string{name, tag}, ":")
	ctx := context.Background()
	reader, err := d.cli.ImageSave(ctx, []string{imageName})
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	tarFilePath := filepath.Join(d.imagePath, strings.Join([]string{tag, "tar"}, "."))
	file, fileErr := os.OpenFile(tarFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	defer file.Close()

	bufferWriter := bufio.NewWriterSize(file, 4096)

	bufferBytes := make([]byte, 4096)

	for count, _ := reader.Read(bufferBytes); count > 0; count, _ = reader.Read(bufferBytes) {
		count, writeErr := bufferWriter.Write(bufferBytes)
		if writeErr != nil {
			log.Fatal(writeErr)
		}
		flushErr := bufferWriter.Flush()
		if flushErr != nil {
			log.Fatal(flushErr)
		}
		log.Printf("Write byte count is:%d!", count)
	}
}

func NewDockerClient() *DockerClient {
	dockerClient := &DockerClient{
		imagePath: "/media/liuxu/data/component/bitbucket/docker-image-gossip/images",
	}
	dockerClient.init()
	return dockerClient
}
