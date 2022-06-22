/*
 * @Description:放置docker相关的rest的操作
 */
package rest

import (
	"fmt"
	"image/gossip/docker"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//
// 放置RESTFul接口部分的代码逻辑部分
//

var (
	dockerClient = docker.NewDockerClient()
	imageDataDir = "./data/images"
)

func LoadDockerImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	targetFilePath := path.Join(imageDataDir, file.Filename)
	c.SaveUploadedFile(file, targetFilePath)
	dockerClient.LoadImage(targetFilePath)
	c.String(http.StatusOK, fmt.Sprintf("%s uploaded success!", file.Filename))
}

func init() {
	err := os.MkdirAll(imageDataDir, os.ModePerm)
	if err != nil {
		logrus.Fatalln("Create dir error")
	}
}
