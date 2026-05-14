package method

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func (info *RobotInfo) Alart(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.WithError(err).Error("read request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
		return
	}

	status := gjson.Get(string(body), "status").Str
	ct := gjson.GetMany(string(body),
		"alerts.0.annotations.summary",
		"alerts.0.labels.desc",
		"alerts.0.annotations.value")
	cTags := []string{"告警：", "单位：", "状态(报警时的触发状态，恢复后的详细状态定请访问：http://prometheus.example.com:9090/"}

	var bb strings.Builder
	if status == "resolved" {
		bb.WriteString("故障已解决：")
	} else if status == "firing" {
		bb.WriteString("故障警告：")
	} else {
		bb.WriteString(status)
	}
	bb.WriteString("\n\n")
	for i, j := range ct {
		bb.WriteString(cTags[i])
		bb.WriteString(j.Str)
		bb.WriteString("\n")
	}

	notifiers, ok := info.routes["alert"]
	if !ok || len(notifiers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no alert route configured"})
		return
	}

	for _, n := range notifiers {
		if err := n.Send(bb.String()); err != nil {
			logrus.WithError(err).Error("notify failed")
		}
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
