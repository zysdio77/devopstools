package method

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (info *RobotInfo) Forward(c *gin.Context) {
	group := c.Param("group")
	notifiers, ok := info.routes[group]
	if !ok || len(notifiers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown group: " + group})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.WithError(err).Error("read request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
		return
	}

	for _, n := range notifiers {
		if err := n.Send(string(body)); err != nil {
			logrus.WithError(err).Error("notify failed")
		}
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
