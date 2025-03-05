package podcasts

import (
	"net/http"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	podcastsDto "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/services"
	userService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/services"
	"github.com/labstack/echo/v4"
)

type PodcastHandler struct {
	PodcastService *podcasts.PodcastService
	UserService    *userService.UserService
}

func NewPodcastHandler(
	podcastService *podcasts.PodcastService,
	userService *userService.UserService,
) *PodcastHandler {
	return &PodcastHandler{
		PodcastService: podcastService,
		UserService:    userService,
	}
}

func (h *PodcastHandler) GetAllPodcasts(c echo.Context) error {
	var getAllPodcastsRequestDto podcastsDto.GetAllPodcastsRequestDto
	res, ok := base.BindAndValidate(c, &getAllPodcastsRequestDto)
	if !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	getAllPodcastsRequestDto.BindPaginationParams(c)

	podcasts, err := h.PodcastService.GetAllPodcasts(c, getAllPodcastsRequestDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, podcasts)
}

func (h *PodcastHandler) GetRecommendedPodcasts(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}

	userCategoriesIds, err := h.UserService.GetUserCategoriesIds(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	recommendedPodcasts, err := h.PodcastService.GetRecommendedPodcasts(userID, userCategoriesIds)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, recommendedPodcasts)
}
