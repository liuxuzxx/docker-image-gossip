/*
 * @Description: 压测Promtail的日志搜集缓慢的问题，收集到的时间和实际的时间差别很大，怀疑是CPU的问题，现在需要测试下
 */
package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoadLog(c *gin.Context) {
	countStr := c.Param("count")
	logrus.Info("获取到的数据是:", countStr)
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logrus.Warn("转换数字错误，请重新传递过来:", countStr)
		c.String(http.StatusBadRequest, fmt.Sprintf("错误的count参数:%s", countStr))
		return
	}
	for start := int64(0); start < count; start = start + 1 {
		logrus.Trace("trace msg")
		logrus.Debug("debug msg")
		logrus.Info("info msg")
		logrus.Warn("warn msg")
		logrus.Error("error msg")
	}
	c.String(http.StatusOK, fmt.Sprintf("%v uploaded success!", time.Now()))
}
