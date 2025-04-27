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

func NewPodcastHandler(podcastService *podcasts.PodcastService, userService *userService.UserService) *PodcastHandler {
	return &PodcastHandler{
		PodcastService: podcastService,
		UserService:    userService,
	}
}

func (h *PodcastHandler) GetAllPodcasts(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}

	var getAllPodcastsRequestDto podcastsDto.GetAllPodcastsRequestDto
	res, ok := base.BindAndValidate(c, &getAllPodcastsRequestDto)
	if !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	getAllPodcastsRequestDto.BindPaginationParams(c)

	response := h.PodcastService.GetAllPodcasts(getAllPodcastsRequestDto, userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *PodcastHandler) GetRecommendedPodcasts(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}

	userCategoriesIDs, err := h.UserService.GetUserCategoriesIDs(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := h.PodcastService.GetRecommendedPodcasts(userID, userCategoriesIDs)
	return c.JSON(response.HTTPStatus, response)
}

func (h *PodcastHandler) GetTrendingPodcasts(c echo.Context) error {
	userID, _ := c.Get("user_id").(string)

	response := h.PodcastService.GetTrendingPodcasts(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *PodcastHandler) GetPodcastDetails(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}

	var getPodcastDetailsRequestDto podcastsDto.GetPodcastDetailsRequestDto
	res, ok := base.BindAndValidate(c, &getPodcastDetailsRequestDto)
	if !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.PodcastService.GetPodcastDetails(getPodcastDetailsRequestDto.ID, userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *PodcastHandler) LikePodcast(c echo.Context) error {
	var reqDto podcastsDto.LikePodcastRequestDto
	if res, ok := base.BindAndValidate(c, &reqDto); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	podcastID := c.Param("podcast_id")
	response := h.PodcastService.LikePodcast(podcastID, reqDto.AddLikes)
	return c.JSON(response.HTTPStatus, response)
}

func (h *PodcastHandler) GetPodcastsByCategory(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}

	var getPodcastsByCategoryRequestDto podcastsDto.GetPodcastsByCategoryRequestDto
	res, ok := base.BindAndValidate(c, &getPodcastsByCategoryRequestDto)
	if !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	getPodcastsByCategoryRequestDto.BindPaginationParams(c)

	response := h.PodcastService.GetPodcastsByCategory(getPodcastsByCategoryRequestDto, userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *PodcastHandler) DownloadPodcast(c echo.Context) error {
	podcastID := c.Param("podcast_id")
	userID := c.Get("user_id").(string)

	response := h.PodcastService.ToggleDownloadPodcast(userID, podcastID)
	return c.JSON(response.HTTPStatus, response)
}
