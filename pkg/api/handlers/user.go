package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sdslabs/status/pkg/api/response"
	"github.com/sdslabs/status/pkg/database"
)

// GetUser fetches the user by ID or Email.
func GetUser(ctx *gin.Context) {
	parameter := ctx.Param("id")
	isEmail := false
	id, err := strconv.Atoi(parameter)
	if err != nil {
		isEmail = ValidateEmail(parameter)
		if !isEmail {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		u, err := database.GetUserByEmail(parameter)
		GetUserHelperFunc(ctx, u, err)
	}

	u, err := database.GetUserByID(uint(id))
	GetUserHelperFunc(ctx, u, err)
}

// DeleteUser deletes the current user.
func DeleteUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := database.DeleteUserByID(uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	ctx.Status(http.StatusOK)
}

// UpdateUser updates the user.
func UpdateUser(ctx *gin.Context) {
	// validates input
	var updateUser database.User
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	u, err := database.GetUserByID(uint(id))
	user, err := database.UpdateUserNameByEmail(u.Email, &updateUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.HTTPUserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// ValidateEmail validates given string as an email address.
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

// GetUserHelperFunc is a helper function for GetUser.
func GetUserHelperFunc(ctx *gin.Context, user *database.User, err error) {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPUserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}
