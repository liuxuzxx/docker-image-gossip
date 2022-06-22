/*
 * @Description:放置docker相关的rest的操作
 */
package rest

import (
	"fmt"
	"image/gossip/docker"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// 放置RESTFul接口部分的代码逻辑部分
//

var (
	dockerClient = docker.NewDockerClient()
)

func LoadDockerImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	targetFilePath := "./data/" + file.Filename
	c.SaveUploadedFile(file, targetFilePath)
	c.String(http.StatusOK, fmt.Sprintf("%s uploaded success!", file.Filename))

	dockerClient.LoadImage(targetFilePath)
}
