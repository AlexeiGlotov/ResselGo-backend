package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"RecurroControl/models"
)

func (h *Handler) createAccessKeys(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	var access_key models.AccessKey
	if err := c.BindJSON(&access_key); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err, ErrNotFields)
		return
	}

	users, err := h.services.Users.GetUser(userID)

	if users.Role == models.Salesman {
		newErrorResponse(c, http.StatusForbidden, nil, ErrAccessDenied)
		return
	}

	if users.Role == models.Distributors {
		if access_key.Role == models.Distributors || access_key.Role == models.Admin {
			newErrorResponse(c, http.StatusForbidden, nil, ErrAccessDenied)
			return
		}
	}

	if users.Role == models.Reseller {
		if access_key.Role == models.Distributors || access_key.Role == models.Admin || access_key.Role == models.Reseller {
			newErrorResponse(c, http.StatusForbidden, nil, ErrAccessDenied)
			return
		}
	}

	if users.Role != models.Admin && users.Role != models.Distributors && users.Role != models.Reseller && users.Role != models.Salesman {
		newErrorResponse(c, http.StatusForbidden, nil, ErrAccessDenied)
		return
	}

	key, err := h.services.AccessKeys.CreateAccessKey(userID, access_key.Role)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err, ErrServerError)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"key": key,
	})

}

func (h *Handler) getAccessKey(c *gin.Context) {

	userID, err := getUserId(c)
	if err != nil {
		return
	}

	user, err := h.services.Users.GetUser(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err, ErrServerError)
		return
	}

	if user.Role == models.Salesman {
		newErrorResponse(c, http.StatusUnauthorized, err, ErrAccessDenied)
		return
	}

	key, err := h.services.AccessKeys.GetAccessKey(user.Login, user.Role)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err, ErrServerError)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"key": key,
	})
}
