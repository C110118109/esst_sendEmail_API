package linebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"esst_sendEmail/internal/pkg/log"
)

// LineBotService LINE Bot 服務介面
type LineBotService interface {
	SendProjectStep1Notification(data *ProjectStep1Data) error
	SendProjectStep2Notification(data *ProjectStep2Data) error
}

type lineBotService struct {
	channelAccessToken string
	groupID            string
	httpClient         *http.Client
}

// ProjectStep1Data 第一階段專案資料
type ProjectStep1Data struct {
	ProjectID    string
	ProjectName  string
	ContactName  string
	ContactPhone string
	ContactEmail string
	Owner        string
	Remark       string
	Equipments   []Equipment
	CreatedTime  time.Time
}

// ProjectStep2Data 第二階段專案資料
type ProjectStep2Data struct {
	ProjectID              string
	ProjectName            string
	ContactName            string
	ExpectedDeliveryPeriod string
	ExpectedDeliveryDate   string
	ExpectedContractPeriod string
	ContractStartDate      string
	ContractEndDate        string
	DeliveryAddress        string
	SpecialRequirements    string
	Equipments             []Equipment
	UpdatedTime            time.Time
}

// Equipment 設備資料
type Equipment struct {
	PartNumber  string
	Quantity    int64
	Description string
}

// LINE Messaging API 的訊息結構
type lineMessage struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}

type textMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// New 建立新的 LINE Bot 服務
func New() LineBotService {
	return &lineBotService{
		channelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		groupID:            os.Getenv("LINE_GROUP_ID"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendProjectStep1Notification 發送第一階段專案報備通知
func (s *lineBotService) SendProjectStep1Notification(data *ProjectStep1Data) error {
	message := s.buildStep1Message(data)
	return s.sendMessage(message)
}

// SendProjectStep2Notification 發送第二階段專案報備通知
func (s *lineBotService) SendProjectStep2Notification(data *ProjectStep2Data) error {
	message := s.buildStep2Message(data)
	return s.sendMessage(message)
}

// buildStep1Message 建立第一階段訊息
func (s *lineBotService) buildStep1Message(data *ProjectStep1Data) string {
	var msg bytes.Buffer

	msg.WriteString("📋 【專案報備通知 - 第一階段】\n")
	msg.WriteString("━━━━━━━━━━━━━━━━━━━━\n\n")

	// 專案基本資訊
	msg.WriteString("📌 專案基本資訊\n")
	msg.WriteString(fmt.Sprintf("• 專案編號: %s\n", data.ProjectID))
	msg.WriteString(fmt.Sprintf("• 專案名稱: %s\n", data.ProjectName))
	msg.WriteString(fmt.Sprintf("• 建立時間: %s\n", data.CreatedTime.Format("2006-01-02 15:04:05")))
	msg.WriteString(fmt.Sprintf("• 目前狀態: 第一階段\n\n"))

	// 聯絡人資訊
	msg.WriteString("👤 聯絡人資訊\n")
	msg.WriteString(fmt.Sprintf("• 聯絡人: %s\n", data.ContactName))
	if data.ContactPhone != "" {
		msg.WriteString(fmt.Sprintf("• 電話: %s\n", data.ContactPhone))
	}
	if data.ContactEmail != "" {
		msg.WriteString(fmt.Sprintf("• 信箱: %s\n", data.ContactEmail))
	}
	msg.WriteString(fmt.Sprintf("• 雙欣負責人: %s\n\n", data.Owner))

	// 設備清單
	if len(data.Equipments) > 0 {
		msg.WriteString("🔧 設備清單\n")
		for i, eq := range data.Equipments {
			msg.WriteString(fmt.Sprintf("%d. 料號: %s\n", i+1, eq.PartNumber))
			msg.WriteString(fmt.Sprintf("    數量: %d\n", eq.Quantity))
			if eq.Description != "" {
				msg.WriteString(fmt.Sprintf("    說明: %s\n", eq.Description))
			}
		}
		msg.WriteString("\n")
	}

	// 備註
	if data.Remark != "" {
		msg.WriteString("📝 備註\n")
		msg.WriteString(fmt.Sprintf("%s\n\n", data.Remark))
	}

	msg.WriteString("⚠️ 提醒:專案得標後,請記得填寫第二階段交貨資訊")

	return msg.String()
}

// buildStep2Message 建立第二階段訊息
func (s *lineBotService) buildStep2Message(data *ProjectStep2Data) string {
	var msg bytes.Buffer

	msg.WriteString("✅ 【專案報備通知 - 第二階段完成】\n")
	msg.WriteString("━━━━━━━━━━━━━━━━━━━━\n\n")

	// 專案基本資訊
	msg.WriteString("📌 專案基本資訊\n")
	msg.WriteString(fmt.Sprintf("• 專案編號: %s\n", data.ProjectID))
	msg.WriteString(fmt.Sprintf("• 專案名稱: %s\n", data.ProjectName))
	msg.WriteString(fmt.Sprintf("• 聯絡人: %s\n", data.ContactName))
	msg.WriteString(fmt.Sprintf("• 更新時間: %s\n", data.UpdatedTime.Format("2006-01-02 15:04:05")))
	msg.WriteString(fmt.Sprintf("• 目前狀態: 第二階段完成\n\n"))

	// 交貨資訊
	msg.WriteString("📦 交貨資訊\n")
	msg.WriteString(fmt.Sprintf("• 預計交貨期: %s\n", data.ExpectedDeliveryPeriod))
	msg.WriteString(fmt.Sprintf("• 預計交貨日: %s\n", formatDate(data.ExpectedDeliveryDate)))
	msg.WriteString(fmt.Sprintf("• 預計履約期: %s\n", data.ExpectedContractPeriod))

	if data.ContractStartDate != "" && data.ContractStartDate != "-" {
		msg.WriteString(fmt.Sprintf("• 履約開始日: %s\n", formatDate(data.ContractStartDate)))
	}
	if data.ContractEndDate != "" && data.ContractEndDate != "-" {
		msg.WriteString(fmt.Sprintf("• 履約結束日: %s\n", formatDate(data.ContractEndDate)))
	}
	msg.WriteString("\n")

	// 設備清單
	if len(data.Equipments) > 0 {
		msg.WriteString("🔧 設備清單\n")
		for i, eq := range data.Equipments {
			msg.WriteString(fmt.Sprintf("%d. 料號: %s\n", i+1, eq.PartNumber))
			msg.WriteString(fmt.Sprintf("   數量: %d\n", eq.Quantity))
			if eq.Description != "" {
				msg.WriteString(fmt.Sprintf("   說明: %s\n", eq.Description))
			}
		}
		msg.WriteString("\n")
	}

	// 交貨地址
	if data.DeliveryAddress != "" {
		msg.WriteString("📍 交貨地址\n")
		msg.WriteString(fmt.Sprintf("%s\n\n", data.DeliveryAddress))
	}

	// 特殊需求
	if data.SpecialRequirements != "" {
		msg.WriteString("⚡ 特殊需求\n")
		msg.WriteString(fmt.Sprintf("%s\n\n", data.SpecialRequirements))
	}

	msg.WriteString("✨ 專案第二階段交貨資訊已完整填寫")

	return msg.String()
}

// sendMessage 發送訊息到 LINE 群組
func (s *lineBotService) sendMessage(text string) error {
	if s.channelAccessToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}
	if s.groupID == "" {
		return fmt.Errorf("LINE_GROUP_ID is not set")
	}

	// 檢查訊息長度 (LINE 限制 5000 字元)
	if len(text) > 5000 {
		log.Info("Message too long, splitting into multiple messages")
		return s.sendLongMessage(text)
	}

	message := lineMessage{
		To: s.groupID,
		Messages: []interface{}{
			textMessage{
				Type: "text",
				Text: text,
			},
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Error("Failed to marshal message:", err)
		return err
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Failed to create request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.channelAccessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Error("Failed to send request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		log.Error("LINE API error:", errorResponse)
		return fmt.Errorf("LINE API returned status %d: %v", resp.StatusCode, errorResponse)
	}

	log.Info("LINE notification sent successfully to group:", s.groupID)
	return nil
}

// sendLongMessage 發送長訊息(分割成多則)
func (s *lineBotService) sendLongMessage(text string) error {
	const maxLength = 4500 // 留一些緩衝空間

	for len(text) > 0 {
		end := len(text)
		if end > maxLength {
			end = maxLength
			// 嘗試在換行處分割
			if idx := bytes.LastIndexByte([]byte(text[:end]), '\n'); idx > 0 {
				end = idx
			}
		}

		if err := s.sendMessage(text[:end]); err != nil {
			return err
		}

		text = text[end:]
		if len(text) > 0 {
			time.Sleep(time.Second) // 避免發送太快
		}
	}

	return nil
}

// formatDate 格式化日期
func formatDate(dateStr string) string {
	if dateStr == "" || dateStr == "-" {
		return "-"
	}
	return dateStr
}
