package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sdslabs/status/pkg/api/app/oauth"
	"github.com/sdslabs/status/pkg/api/response"
	"github.com/sdslabs/status/pkg/database"
)

// UpdateUser contains the user's information to be updated.
type UpdateUser struct {
	Name string
}

// GetUser fetches the user by ID or Email.
func GetUser(ctx *gin.Context) {
	parameter := ctx.Param("id")
	isEmail := false
	id, err := strconv.Atoi(parameter)
	if err != nil {
		isEmail = validateEmail(parameter)
		if !isEmail {
			ctx.JSON(http.StatusBadRequest, response.HTTPError{
				Error: err.Error(),
			})
			return
		}
		u, err := database.GetUserByEmail(parameter)
		getUserResponse(ctx, u, err)
		return
	}

	u, err := database.GetUserByID(uint(id))
	getUserResponse(ctx, u, err)
}

// GetCurrentUser fetches the current user.
func GetCurrentUser(ctx *gin.Context) {
	user, ok := oauth.CurrentUserFromCtx(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.HTTPError{
			Error: "user not logged in",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.HTTPUserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// DeleteCurrentUser deletes the current user.
func DeleteCurrentUser(ctx *gin.Context) {
	user, ok := oauth.CurrentUserFromCtx(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.HTTPError{
			Error: "user not logged in",
		})
		return
	}

	if err := database.DeleteUserByID(user.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.HTTPEmpty{})
}

// UpdateCurrentUser updates the user.
func UpdateCurrentUser(ctx *gin.Context) {
	// validates input
	var updateUser UpdateUser
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	u, ok := oauth.CurrentUserFromCtx(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.HTTPError{
			Error: "user not logged in",
		})
		return
	}

	v := database.User{}
	v.Name = updateUser.Name

	user, err := database.UpdateUserByID(u.ID, &v)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.HTTPError{
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

// validateEmail validates given string as an email address.
func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

// getUserResponse is a helper function for GetUser.
func getUserResponse(ctx *gin.Context, user *database.User, err error) {
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
