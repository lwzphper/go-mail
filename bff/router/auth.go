package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lwzphper/go-mall/bff/api/auth"
)

func InitAuth(r *gin.RouterGroup) {
	authGroup := r.Group("auth") //.Use(middleware.Trace())
	{
		authGroup.POST("login", auth.Login)
		authGroup.POST("register", auth.Register)
	}
}
