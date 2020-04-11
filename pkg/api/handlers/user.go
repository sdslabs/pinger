package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/pkg/api/reponse"
	"github.com/sdslabs/status/pkg/database"
)

userRouter := router.Group("/users")
	userRouter.GET("/:id", func(ctx *gin.Context) {
		u , err := GetUserByID(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				"Error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"name":  u.Name,
			"email": u.Email,
		})
	}	
	userRouter.DELETE("/delete/:id", func(ctx *gin.Context)) {
		if err := DeleteUserByID(ctx.Param("id")); err := nil {
			ctx.JSON(http.StatusBadRequest, reponse.HTTPError{
				"Error": err.Error(),
			})
			return
		} 
		ctx.JSON(http.StatusOK, gin.H{
			"action": true,
		})
	}
	userRouter.PUT("/update/:id", func(ctx *gin.Context)) {
		var updateUser user
		User, err := UpdateUserByID(ctx.Param("id"), &updateUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				"Error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"name":  User.Name,
			"email": User.Email,
		})
	}
