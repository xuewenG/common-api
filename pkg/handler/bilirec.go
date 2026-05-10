package handler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/xuewenG/common-api/pkg/bark"
)

// EventData 表示事件数据结构
type EventData struct {
	RoomId int64  `json:"RoomId"`
	Name   string `json:"Name"`
	Title  string `json:"Title"`
}

// Event 表示完整的事件结构
type Event struct {
	EventType string     `json:"EventType"`
	EventId   string     `json:"EventId"`
	EventData *EventData `json:"EventData"`
}

// RoomEventProcessor 表示房间事件处理器
type RoomEventProcessor struct {
	roomId     int64
	isLive     bool
	logger     *zerolog.Logger
	mutex      sync.Mutex
	cancelFunc context.CancelFunc
}

// BiliRecHandler 表示事件处理
type BiliRecHandler struct {
	*Handler
	barkClient      *bark.BarkClient
	roomProcessors  map[int64]*RoomEventProcessor
	processorsMutex sync.RWMutex
}

func NewBiliRecHandler(i do.Injector) (*BiliRecHandler, error) {
	h, err := newHandler(i)
	if err != nil {
		return nil, err
	}

	return &BiliRecHandler{
		Handler:        h,
		barkClient:     do.MustInvoke[*bark.BarkClient](i),
		roomProcessors: make(map[int64]*RoomEventProcessor),
	}, nil
}

func (h *BiliRecHandler) OnEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.Info().Msgf("解析事件数据失败: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析事件数据失败: " + err.Error()})
		return
	}

	if event.EventType == "" || event.EventData.RoomId == 0 {
		h.logger.Info().Msgf("缺少必要的事件字段")
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要的事件字段"})
		return
	}

	if event.EventType != "StreamStarted" && event.EventType != "StreamEnded" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	h.logger.Info().Msgf("开始处理: %+v", event)

	h.getOrCreateRoomProcessor(
		event.EventData.RoomId,
	).processEvent(
		&event, h.barkClient,
	)

	h.logger.Info().Msgf("处理完成: %+v", event)

	c.JSON(http.StatusOK, gin.H{"message": "事件处理成功"})
}

// getOrCreateRoomProcessor 获取或创建房间事件处理器
func (h *BiliRecHandler) getOrCreateRoomProcessor(roomId int64) *RoomEventProcessor {
	h.processorsMutex.Lock()
	defer h.processorsMutex.Unlock()

	processor, exists := h.roomProcessors[roomId]
	if !exists {
		processor = &RoomEventProcessor{
			roomId: roomId,
			logger: h.logger,
		}
		h.roomProcessors[roomId] = processor
	}

	return processor
}

// processEvent 处理单个事件
func (rp *RoomEventProcessor) processEvent(event *Event, barkClient *bark.BarkClient) {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()

	switch event.EventType {
	case "StreamEnded":
		if rp.cancelFunc != nil {
			rp.cancelFunc()
		}

		ctx, cancel := context.WithCancel(context.Background())
		rp.cancelFunc = cancel

		name := event.EventData.Name

		go func() {
			select {
			case <-time.After(60 * time.Second):
				rp.mutex.Lock()
				rp.isLive = false
				rp.mutex.Unlock()

				err := barkClient.SendMessage(&bark.BarkMessage{
					Title: "停止推流通知",
					Body:  fmt.Sprintf("%s - 已停止推流", name),
				})
				if err != nil {
					rp.logger.Error().Err(err).Msg("发送 Bark 消息失败")
				}
			case <-ctx.Done():
				return
			}
		}()

	case "StreamStarted":
		// 先取消之前的停止通知定时器（如果有）
		if rp.cancelFunc != nil {
			rp.cancelFunc()
			rp.cancelFunc = nil
		}

		if rp.isLive {
			return
		}

		rp.isLive = true

		err := barkClient.SendMessage(&bark.BarkMessage{
			Title: "开始推流通知",
			Body:  fmt.Sprintf("%s - 已开始推流", event.EventData.Name),
			Url:   fmt.Sprintf("bilibili://live/%d", event.EventData.RoomId),
		})
		if err != nil {
			rp.logger.Error().Err(err).Msg("发送 Bark 消息失败")
		}
	}
}
