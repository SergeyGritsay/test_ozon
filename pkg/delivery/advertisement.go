package delivery

import (
	"advertisement-service/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createAdvertisement(c *gin.Context) {
	var input models.Advertisement

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"id":     -1,
			"status": http.StatusUnprocessableEntity,
		})

		return
	}

	adId, err := h.services.Advertisement.CreateAdvertisement(input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"id":     -1,
			"status": http.StatusUnprocessableEntity,
		})

		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":     adId,
		"status": http.StatusOK,
	})
}

func (h *Handler) getAdvertisement(c *gin.Context) {
	var input models.AdvertisementGetRequest

	c.BindJSON(&input)

	if !input.ValidateFields() {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"id":     -1,
			"status": http.StatusUnprocessableEntity,
		})

		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{})

		return
	}

	ad, err := h.services.Advertisement.GetAdvertisement(id, input.Fields)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{})

		return
	}

	c.JSON(http.StatusOK, ad)
}

func (h *Handler) getAdvertisementsList(c *gin.Context) {
	var input models.AdvertisementsGetRequest
	const pageSize = 10

	c.BindJSON(&input)

	if !input.ValidateSort() {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status": http.StatusUnprocessableEntity,
		})

		return
	}

	resp, err := h.services.Advertisement.GetAdvertisements(input.Sort)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status": http.StatusUnprocessableEntity,
		})
		return
	}

	paginatedResp := resp

	start := int(input.Page) * pageSize
	end := (int(input.Page) + 1) * pageSize

	if start < len(resp) {
		if end < len(resp) {
			paginatedResp = paginatedResp[start:end]
		} else {
			paginatedResp = paginatedResp[start:]
		}

		c.JSON(http.StatusOK, paginatedResp)

		return
	} else {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"status": http.StatusUnprocessableEntity,
		})
	}
}
