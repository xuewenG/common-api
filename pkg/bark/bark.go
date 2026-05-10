package bark

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/samber/do/v2"
	"github.com/xuewenG/common-api/pkg/config"
)

// BarkMessage 表示发送到 Bark 的消息结构
type BarkMessage struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Badge     int    `json:"badge,omitempty"`
	AutoCopy  int    `json:"auto_copy,omitempty"`
	Copy      string `json:"copy,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	IsArchive bool   `json:"is_archive,omitempty"`
	Url       string `json:"url,omitempty"`
}

// BarkClient 表示 Bark 客户端
type BarkClient struct {
	config *config.BarkConfig
	client *http.Client
}

// NewBarkClient 创建新的 Bark 客户端
func NewBarkClient(i do.Injector) (*BarkClient, error) {
	cfg := do.MustInvoke[*config.Config](i)

	return &BarkClient{
		config: &cfg.Bark,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// SendMessage 发送消息到 Bark
func (b *BarkClient) SendMessage(message *BarkMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	baseURL := b.config.ServerURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	requestURL := baseURL + b.config.DeviceKey
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("content-type", "application/json")

	if b.config.Username != "" && b.config.Password != "" {
		auth := b.config.Username + ":" + b.config.Password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("authorization", fmt.Sprintf("Basic %s", encodedAuth))
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送 Bark 消息失败: %w", err)
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bark 服务返回错误状态码: %d", resp.StatusCode)
	}

	return nil
}
