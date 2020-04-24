package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/api/app/oauth"
	"github.com/sdslabs/status/pkg/api/request"
	"github.com/sdslabs/status/pkg/api/response"
	"github.com/sdslabs/status/pkg/database"
)

// GetUser fetches the user by ID or Email.
func GetUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	isEmail, param, err := paramIsEmail(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	var u *database.User
	if isEmail {
		email, ok := param.(string)
		if !ok {
			logrus.Errorln("email param returned by `paramIsEmail` invalid")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		u, err = database.GetUserByEmail(email)
	} else {
		userid, ok := param.(uint)
		if !ok {
			logrus.Errorln("id param returned by `paramIsEmail` invalid")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		u, err = database.GetUserByID(userid)
	}

	if err != nil {
		if err != database.ErrRecordNotFound {
			logrus.WithError(err).Errorln("cannot get user")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		ctx.JSON(http.StatusNotFound, response.HTTPError{
			Error: fmt.Sprintf("user with id '%s' not found", idParam),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPUserInfo{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	})
}

// GetCurrentUser fetches the current user.
func GetCurrentUser(ctx *gin.Context) {
	user, ok := oauth.CurrentUserFromCtx(ctx)
	if !ok {
		logrus.Errorln("cannot get user from context bucket")
		ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
		return
	}

	u, err := database.GetUserByID(user.ID)
	if err != nil {
		if err != database.ErrRecordNotFound {
			logrus.WithError(err).Errorln("cannot get user")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		ctx.JSON(http.StatusNotFound, response.HTTPError{
			Error: fmt.Sprintf("user with id '%d' not found", user.ID),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPUserInfo{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	})
}

// DeleteCurrentUser deletes the current user.
func DeleteCurrentUser(ctx *gin.Context) {
	user, ok := oauth.CurrentUserFromCtx(ctx)
	if !ok {
		logrus.Errorln("cannot get user from context bucket")
		ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
		return
	}

	if err := database.DeleteUserByID(user.ID); err != nil {
		if err != database.ErrRecordNotFound {
			logrus.WithError(err).Errorln("cannot delete user")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		ctx.JSON(http.StatusNotFound, response.HTTPError{
			Error: fmt.Sprintf("user with id '%d' not found", user.ID),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPEmpty{})
}

// UpdateCurrentUser updates the user.
func UpdateCurrentUser(ctx *gin.Context) {
	var updateUser request.HTTPUserUpdate
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, response.HTTPError{
			Error: err.Error(),
		})
		return
	}

	user, ok := oauth.CurrentUserFromCtx(ctx)
	if !ok {
		logrus.Errorln("cannot get user from context bucket")
		ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
		return
	}

	update := database.User{Name: updateUser.Name}

	u, err := database.UpdateUserByID(user.ID, &update)
	if err != nil {
		if err != database.ErrRecordNotFound {
			logrus.WithError(err).Errorln("cannot update user")
			ctx.JSON(http.StatusInternalServerError, response.HTTPInternalServerError)
			return
		}

		ctx.JSON(http.StatusNotFound, response.HTTPError{
			Error: fmt.Sprintf("user with id '%d' not found", user.ID),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.HTTPUserInfo{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	})
}
