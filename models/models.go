package models

import "github.com/gin-gonic/gin"

type StandardRequest struct {
	Name string
	Code string
}

type StandardResponse struct {
	ID   string
	Name string
	Code string
}

func Authenticate(g *gin.Context) {

}
