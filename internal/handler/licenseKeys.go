package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"RecurroControl/models"
)

func generateUniqueKey() string {
	const charset = "ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnopqrstuvwxyz123456789"

	seed := int64(time.Now().UnixNano())
	src := rand.NewSource(seed)
	r := rand.New(src)

	key := make([]byte, 36)
	for i := 0; i < 36; i++ {
		key[i] = charset[r.Intn(len(charset))]
	}

	return string(key)
}

func (h *Handler) createLicenseKeys(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	var inputKeys models.InputCreate
	if err := c.BindJSON(&inputKeys); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.Users.GetUser(userID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if inputKeys.CountKeys > 50 {
		newErrorResponse(c, http.StatusBadRequest, "Превышено максимальное количество ключей")
		return
	}

	var licenseKeys []models.LicenseKeys

	for x := 0; x < inputKeys.CountKeys; x++ {
		var tempLicenseKey models.LicenseKeys

		tempLicenseKey.Creator = user.Login
		tempLicenseKey.Holder = inputKeys.Holder
		tempLicenseKey.TTLCheat = inputKeys.TTLCheat
		tempLicenseKey.Cheat = inputKeys.Cheat
		tempLicenseKey.LicenseKeys = generateUniqueKey()

		licenseKeys = append(licenseKeys, tempLicenseKey)
	}

	err = h.services.LicenseKeys.CreateLicenseKeys(licenseKeys)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"keys": licenseKeys})
}

func (h *Handler) getLicenseKeys(c *gin.Context) {

	userID, err := getUserId(c)
	if err != nil {
		return
	}

	pageStr, ok := c.GetQuery("page")
	if !ok {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат страницы"})
		return
	}

	offset := (page - 1) * 100

	keys, err := h.services.LicenseKeys.GetLicenseKeys(userID, 100, offset)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(keys) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "bad page")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"keys": keys})
}
