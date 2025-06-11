package handlers

import (
	"github.com/gin-gonic/gin"
	logs "github.com/j3rryCodes/system-monitor/internal/handlers/log"
	"github.com/j3rryCodes/system-monitor/internal/handlers/system"
)

func Handle(r *gin.Engine) {
	system.RegisterRoute(r)
	logs.RegisterRoute(r)
}
