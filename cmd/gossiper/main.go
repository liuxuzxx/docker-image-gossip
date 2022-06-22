/*
 * @Description:
 */
/*
 * @Description: 执行的主类
 */
package main

import (
	"bytes"
	"image/gossip/rest"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	initLog()
	go initPrometheus()
	logrus.Info("Start docker image gossip ...")
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/docker/load-image", rest.LoadDockerImage)
		v1.POST("/docker/gossip-image", rest.GossipImage)
		v1.GET("/promtail/:count/load-logs", rest.LoadLog)
		v1.POST("/log/statistics", rest.StatisticsCount)
	}
	router.Run(":8080")
}

func initPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9599", nil)
}

func initLog() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	writerBuffer := &bytes.Buffer{}
	stdoutWriter := os.Stdout
	os.Mkdir("logs", os.ModePerm)
	fileWriter, err := os.OpenFile("logs/docker-image-gossip.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(writerBuffer, stdoutWriter, fileWriter))
}
