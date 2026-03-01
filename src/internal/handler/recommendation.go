package handler

import (
	"errors"
	"fmt"
	"net/http"
	"recommendation-system/src/internal/model/interfaces"
	_appErr "recommendation-system/src/utils/error"
	"strconv"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type recommendationHandler struct {
	recommendationService interfaces.RecommendationService
}

func NewRecommendationHandler(recommendationService interfaces.RecommendationService) *recommendationHandler {
	return &recommendationHandler{
		recommendationService,
	}
}

func (h *recommendationHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/users/:userID/recommendations", h.GetUserRecommendations)
	e.GET("/recommendations/batch", h.GetBatchRecommendation)

}

func (h *recommendationHandler) GetUserRecommendations(c *echo.Context) error {
	ctx := c.Request().Context()
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		errObj := _appErr.HttpResponseError{
			Error:   "invalid_parameter",
			Message: "Invalid userID parameter",
		}
		return c.JSON(http.StatusBadRequest, errObj)
	}

	limit := 10
	rawLimit := c.QueryParam("limit")
	fmt.Println("rawLimit", rawLimit)
	if rawLimit != "" {
		limit, err = strconv.Atoi(rawLimit)
		if err != nil {
			errObj := _appErr.HttpResponseError{
				Error:   "invalid_parameter",
				Message: "Invalid limit parameter",
			}
			return c.JSON(http.StatusBadRequest, errObj)
		}
		if limit < 0 || limit > 50 {
			errObj := _appErr.HttpResponseError{
				Error:   "invalid_parameter",
				Message: "Invalid limit parameter",
			}
			return c.JSON(http.StatusBadRequest, errObj)
		}
	}

	response, err := h.recommendationService.GetUserRecommendations(ctx, userID, limit)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errObj := _appErr.HttpResponseError{
			Error:   "internal_error",
			Message: err.Error(),
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user not found" {
			statusCode = http.StatusNotFound
			errObj.Error = "user_not_found"
			errObj.Message = fmt.Sprintf("User with ID %v does not exist", userID)
		}

		return c.JSON(statusCode, errObj)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *recommendationHandler) GetBatchRecommendation(c *echo.Context) error {
	type Params struct {
		Page  int `query:"page"`
		Limit int `query:"limit"`
	}

	params := Params{
		Page:  1,
		Limit: 20,
	}
	if err := c.Bind(&params); err != nil {
		errObj := _appErr.HttpResponseError{
			Error:   "invalid_parameter",
			Message: "Invalid query parameters",
		}
		return c.JSON(http.StatusBadRequest, errObj)
	}

	if params.Page <= 0 {
		errObj := _appErr.HttpResponseError{
			Error:   "invalid_parameter",
			Message: "Invalid page parameter",
		}
		return c.JSON(http.StatusBadRequest, errObj)
	}

	if params.Limit < 0 || params.Limit > 100 {
		errObj := _appErr.HttpResponseError{
			Error:   "invalid_parameter",
			Message: "Invalid limit parameter",
		}
		return c.JSON(http.StatusBadRequest, errObj)
	}

	response, err := h.recommendationService.GetBatchRecommendation(c.Request().Context(), params.Page, params.Limit)
	if err != nil {
		errObj := _appErr.HttpResponseError{
			Error:   "internal_error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errObj)
	}

	return c.JSON(http.StatusOK, response)
}
