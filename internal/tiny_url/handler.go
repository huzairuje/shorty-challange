package tiny_url

import (
	"net/http"

	"shorty-challenge/pkg/response"
	"shorty-challenge/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler() *Handler {
	service := NewService()
	return &Handler{
		Service: service,
	}
}

func (p Handler) CreateTinyUrl(ctx *gin.Context) {
	var req Data
	if err := ctx.Bind(&req); err != nil {
		response.BadRequest(ctx, utils.BadRequest, err.Error())
		return
	}

	if req.Url == "" {
		response.BadRequest(ctx, utils.BadRequest, utils.UrlIsNotSet)
		return
	}

	if req.ShortCode != "" {
		if !utils.IsValidShortCode(req.ShortCode) {
			response.BadRequest(ctx, utils.BadRequest, utils.ShortCodeFailedRegexPattern)
			return
		}

		shortIsExist, err := p.Service.GetSingleData(req.ShortCode)
		if shortIsExist != nil || err != nil {
			response.BadRequest(ctx, utils.BadRequest, utils.ShortCodeExist)
			return
		}
	}

	shortCode, err := p.Service.CreateData(req)
	if err != nil {
		response.InternalServerError(ctx, utils.SomethingWentWrong, err)
		return
	}

	type resp struct {
		ShortCode string `json:"shortCode"`
	}

	data := &resp{
		ShortCode: shortCode,
	}
	response.DataWithoutMeta(ctx, data)
	return
}

func (p Handler) SingleTinyUrl(ctx *gin.Context) {
	shortCode := ctx.Param("shortcode")
	data, err := p.Service.GetSingleData(shortCode)
	if data == nil || err != nil {
		response.CustomResponse(ctx, http.StatusNotFound, utils.NotFound, utils.ShortCodeIsNotExist)
		return
	}
	if data != nil {
		p.Service.UpdateStat(data.ShortCode)
		ctx.Writer.WriteHeader(http.StatusFound)
		ctx.Header("Location", data.Url)
		return
	}
}

func (p Handler) StatSingleTinyUrl(ctx *gin.Context) {
	shortCode := ctx.Param("shortcode")
	data, err := p.Service.GetSingleData(shortCode)
	if data == nil || err != nil {
		response.CustomResponse(ctx, http.StatusNotFound, utils.NotFound, utils.ShortCodeIsNotExist)
		return
	}
	type responseStat struct {
		StartDate     string `json:"startDate"`
		LastSeenDate  string `json:"lastSeenDate"`
		RedirectCount int    `json:"redirectCount"`
	}
	resp := &responseStat{
		StartDate:     data.StartDate,
		LastSeenDate:  data.LastSeenDate,
		RedirectCount: data.RedirectCount,
	}
	response.DataWithoutMeta(ctx, resp)
}
