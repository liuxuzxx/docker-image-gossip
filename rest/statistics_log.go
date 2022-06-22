/*
 * @Description:统计线上的文件的个数信息
 */
package rest

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StatisticsCount(c *gin.Context) {
	param := &StatisticsRequest{}
	c.BindJSON(&param)
	logrus.Printf("查看接收到的参数信息:%v\n!", param)
	go doCount(param)
}

func doCount(request *StatisticsRequest) {
	f, err := os.Create(request.TaragetFilePath)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	count := 0
	fi, err := os.Open(request.FilePath)
	if err != nil {
		logrus.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), request.KeyWord) {
			count = count + 1
			fmt.Fprintln(w, string(a))
		}
	}
	w.Flush()
	logrus.Printf("文件:%s  含有关键字:%s 个数是:%d\n!", request.FilePath, request.KeyWord, count)

}

type StatisticsRequest struct {
	FilePath        string `json:"filePath"`
	KeyWord         string `json:"keyWord"`
	TaragetFilePath string `json:"targetFilePath"`
}
