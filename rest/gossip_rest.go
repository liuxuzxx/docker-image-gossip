/*
 * @Description:专门用来处理gossip广播image的操作
 */
package rest

import (
	"fmt"
	"image/gossip/gossip"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//
// 定义一个专门处理gossip协议上传镜像的接口
//

const (
	fileParam = "file"
	hostName  = "docker-image-gossip-headless-service.share-components.svc.cluster.local"
	dataDir   = "./data"
)

var (
	gossiper = gossip.NewGossiper(hostName)
)

func GossipImage(c *gin.Context) {
	file, err := c.FormFile(fileParam)
	if err != nil {
		logrus.Warnln("No file parameter!")
		c.String(http.StatusBadRequest, fmt.Sprintf("%s parameter not exists!", fileParam))
		return
	}

	targetFilePath := strings.Join([]string{dataDir, file.Filename}, "/")
	c.SaveUploadedFile(file, targetFilePath)
	c.String(http.StatusOK, fmt.Sprintf("%s uploaded success!", file.Filename))

	go dockerClient.LoadImage(targetFilePath)
	gossiper.SpreadImages(targetFilePath)
}

func init() {
	os.MkdirAll(dataDir, os.ModePerm)
	logrus.Infoln("创建数据目录：", dataDir)
}
