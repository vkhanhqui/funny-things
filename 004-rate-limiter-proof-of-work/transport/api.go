package transport

import (
	"net/http"
	"proof-of-work/app"

	"github.com/gin-gonic/gin"
)

type PowAPI interface {
	GetIssue(c *gin.Context)
	VerifyIssue(c *gin.Context)
}

type powAPI struct {
	svc app.PowService
}

func NewPowAPI(svc app.PowService) PowAPI {
	return &powAPI{svc}
}

func (r *powAPI) GetIssue(c *gin.Context) {
	issue := r.svc.GetIssue()
	c.IndentedJSON(http.StatusOK, issue)
}

func (r *powAPI) VerifyIssue(c *gin.Context) {
	var req app.VerifyIssueReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.svc.VerifyIssue(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}
