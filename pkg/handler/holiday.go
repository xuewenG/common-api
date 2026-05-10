package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

// Holiday 假期数据结构
type Holiday struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// NextHolidayResponse 下一个假期响应结构
type NextHolidayResponse struct {
	Code int      `json:"code"`
	Data *Holiday `json:"data"`
}

type HolidayHandler struct {
	*Handler
}

func NewHolidayHandler(i do.Injector) (*HolidayHandler, error) {
	h, err := newHandler(i)
	if err != nil {
		return nil, err
	}

	return &HolidayHandler{
		Handler: h,
	}, nil
}

// GetNextHoliday 获取下一个假期
func (h *HolidayHandler) GetNextHoliday(c *gin.Context) {
	now := time.Now()
	var nextHoliday *Holiday

	for _, holiday := range h.config.Holidays {
		holidayDate, err := time.Parse("2006-01-02", holiday.Date)
		if err != nil {
			h.logger.Error().Err(err).Str("date", holiday.Date).Msg("解析假期日期失败")
			continue
		}

		if holidayDate.After(now) {
			nextHoliday = &Holiday{
				Name: holiday.Name,
				Date: holiday.Date,
			}
			break
		}
	}

	response := NextHolidayResponse{
		Code: 200,
		Data: nextHoliday,
	}

	c.JSON(http.StatusOK, response)
}
