package transport

import (
	"proof-of-work/app"

	"github.com/gin-gonic/gin"
)

func NewPowRouter(r *gin.RouterGroup, svc app.PowService) {
	p := NewPowAPI(svc)

	r.GET("/", p.GetIssue)
	r.POST("/", p.VerifyIssue)
}
