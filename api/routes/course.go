package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/UTDNebula/nebula-api/api/controllers"
)

func CourseRoute(router *gin.Engine) {
	// All routes related to courses come here
	courseGroup := router.Group("/course")

	courseGroup.OPTIONS("", controllers.Preflight)
	courseGroup.GET("", controllers.CourseSearch)
	courseGroup.GET(":id", controllers.CourseById)
	courseGroup.GET("all", controllers.CourseAll)

	// Endpoint to get the list of sections of the queried course, courses
	courseGroup.GET("/sections", controllers.CourseSectionSearch())
	courseGroup.GET("/:id/sections", controllers.CourseSectionById())
	courseGroup.GET("/professors", controllers.CourseProfessorSearch)
	courseGroup.GET("/:id/professors", controllers.CourseProfessorById)
}
