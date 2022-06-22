/*
 * @Description:专门用来处理gossip广播image的操作
 */
package rest

import (
	"fmt"
	"image/gossip/gossip"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//
// 定义一个专门处理gossip协议上传镜像的接口
//

const (
	fileParam     = "file"
	hostName      = "docker-image-gossip-headless-service.share-components.svc.cluster.local"
	gossipDataDir = "./data/gossip"
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

	targetFilePath := path.Join(gossipDataDir, file.Filename)
	c.SaveUploadedFile(file, targetFilePath)
	c.String(http.StatusOK, fmt.Sprintf("%s uploaded success!", file.Filename))

	go gossiper.SpreadImages(targetFilePath)
}

func init() {
	err := os.MkdirAll(gossipDataDir, os.ModePerm)
	if err != nil {
		logrus.Fatalln("Create dir error:", gossipDataDir)
	} else {
		logrus.Infoln("Create dir success:", gossipDataDir)
	}
}
