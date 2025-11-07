package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"time"

	"esst_sendEmail/internal/pkg/log"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendProjectStep1Email(data *ProjectStep1Data) error
	SendProjectStep2Email(data *ProjectStep2Data) error
}

type emailService struct {
	smtpHost     string
	smtpPort     int
	smtpUser     string
	smtpPassword string
	fromEmail    string
	toEmail      string
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
	Equipments             []Equipment // 設備清單
	UpdatedTime            time.Time
}

// Equipment 設備資料
type Equipment struct {
	PartNumber  string
	Quantity    int64
	Description string
}

func New() EmailService {
	return &emailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     587, // 可以從環境變數讀取
		smtpUser:     os.Getenv("SMTP_USER"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("EMAIL_FROM"),
		toEmail:      os.Getenv("EMAIL_TO"),
	}
}

// SendProjectStep1Email 發送第一階段專案報備通知
func (s *emailService) SendProjectStep1Email(data *ProjectStep1Data) error {
	subject := fmt.Sprintf("【專案報備通知】%s - 第一階段", data.ProjectName)

	htmlBody, err := s.renderStep1Template(data)
	if err != nil {
		log.Error("Failed to render email template:", err)
		return err
	}

	return s.sendEmail(subject, htmlBody)
}

// SendProjectStep2Email 發送第二階段專案報備通知
func (s *emailService) SendProjectStep2Email(data *ProjectStep2Data) error {
	subject := fmt.Sprintf("【專案報備通知】%s - 第二階段完成", data.ProjectName)

	htmlBody, err := s.renderStep2Template(data)
	if err != nil {
		log.Error("Failed to render email template:", err)
		return err
	}

	return s.sendEmail(subject, htmlBody)
}

// sendEmail 實際發送 Email 的函數
func (s *emailService) sendEmail(subject, htmlBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", s.toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)

	// 針對某些 SMTP 伺服器,可能需要跳過 TLS 驗證
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Error("Failed to send email:", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Info("Email sent successfully to:", s.toEmail)
	return nil
}

// renderStep1Template 渲染第一階段 Email 範本
func (s *emailService) renderStep1Template(data *ProjectStep1Data) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { 
            font-family: 'Microsoft JhengHei', Arial, sans-serif; 
            line-height: 1.8; 
            color: #333; 
            background: #f5f5f5; 
            margin: 0; 
            padding: 20px; 
        }
        .container { 
            max-width: 650px; 
            margin: 0 auto; 
            background: white; 
            border: 1px solid #e0e0e0; 
            border-radius: 4px; 
        }
        .header { 
            background: #2c3e50; 
            color: white; 
            padding: 25px 30px; 
            border-bottom: 3px solid #34495e; 
        }
        .header h1 { 
            margin: 0 0 8px 0; 
            font-size: 22px; 
            font-weight: 600; 
        }
        .header p { 
            margin: 0; 
            font-size: 14px; 
            opacity: 0.9; 
        }
        .content { 
            padding: 30px; 
        }
        .section { 
            margin-bottom: 25px; 
            padding-bottom: 20px; 
            border-bottom: 1px solid #e8e8e8; 
        }
        .section:last-child { 
            border-bottom: none; 
        }
        .section-title { 
            font-size: 16px; 
            font-weight: 600; 
            margin-bottom: 15px; 
            color: #2c3e50; 
        }
        .info-row { 
            margin: 10px 0; 
            font-size: 14px; 
        }
        .label { 
            display: inline-block; 
            width: 120px; 
            font-weight: 500; 
            color: #666; 
        }
        .value { 
            color: #333; 
        }
        table { 
            width: 100%; 
            border-collapse: collapse; 
            margin-top: 12px; 
            font-size: 14px; 
        }
        th, td { 
            padding: 12px 10px; 
            text-align: left; 
            border-bottom: 1px solid #e8e8e8; 
        }
        th { 
            background-color: #f8f9fa; 
            font-weight: 600; 
            color: #2c3e50; 
        }
        tbody tr:hover { 
            background-color: #fafafa; 
        }
        .footer { 
            background: #f8f9fa; 
            padding: 20px 30px; 
            text-align: center; 
            font-size: 12px; 
            color: #999; 
            border-top: 1px solid #e8e8e8; 
        }
        .footer p { 
            margin: 5px 0; 
        }
        .badge { 
            display: inline-block; 
            padding: 4px 10px; 
            background: #f0f0f0; 
            color: #666; 
            border-radius: 3px; 
            font-size: 13px; 
            font-weight: 500; 
        }
        .notice-box { 
            background: #f8f9fa; 
            border-left: 3px solid #7f8c8d; 
            padding: 15px; 
            margin-top: 20px; 
        }
        .notice-box p { 
            margin: 0; 
            color: #555; 
            font-size: 14px; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>專案報備通知 - 第一階段</h1>
            <p>新專案已建立並完成第一階段報備</p>
        </div>
        
        <div class="content">
            <div class="section">
                <div class="section-title">專案基本資訊</div>
                <div class="info-row">
                    <span class="label">專案編號</span>
                    <span class="value">{{.ProjectID}}</span>
                </div>
                <div class="info-row">
                    <span class="label">專案名稱</span>
                    <span class="value">{{.ProjectName}}</span>
                </div>
                <div class="info-row">
                    <span class="label">建立時間</span>
                    <span class="value">{{.CreatedTime.Format "2006-01-02 15:04:05"}}</span>
                </div>
                <div class="info-row">
                    <span class="label">目前狀態</span>
                    <span class="badge">第一階段</span>
                </div>
            </div>

            <div class="section">
                <div class="section-title">聯絡人資訊</div>
                <div class="info-row">
                    <span class="label">聯絡人</span>
                    <span class="value">{{.ContactName}}</span>
                </div>
                <div class="info-row">
                    <span class="label">聯絡電話</span>
                    <span class="value">{{.ContactPhone}}</span>
                </div>
                <div class="info-row">
                    <span class="label">聯絡信箱</span>
                    <span class="value">{{.ContactEmail}}</span>
                </div>
                <div class="info-row">
                    <span class="label">雙欣負責人</span>
                    <span class="value">{{.Owner}}</span>
                </div>
            </div>

            {{if .Equipments}}
            <div class="section">
                <div class="section-title">設備清單</div>
                <table>
                    <thead>
                        <tr>
                            <th>料號</th>
                            <th>數量</th>
                            <th>說明</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Equipments}}
                        <tr>
                            <td>{{.PartNumber}}</td>
                            <td>{{.Quantity}}</td>
                            <td>{{if .Description}}{{.Description}}{{else}}-{{end}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{end}}

            {{if .Remark}}
            <div class="section">
                <div class="section-title">備註</div>
                <p style="margin: 0; color: #555; font-size: 14px;">{{.Remark}}</p>
            </div>
            {{end}}

            <div class="notice-box">
                <p><strong>提醒：</strong>專案得標後，請記得填寫第二階段交貨資訊。</p>
            </div>
        </div>
        
        <div class="footer">
            <p>此為系統自動發送的通知郵件，請勿直接回覆</p>
            <p>專案報備系統 © 2025</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// renderStep2Template 渲染第二階段 Email 範本
func (s *emailService) renderStep2Template(data *ProjectStep2Data) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { 
            font-family: 'Microsoft JhengHei', Arial, sans-serif; 
            line-height: 1.8; 
            color: #333; 
            background: #f5f5f5; 
            margin: 0; 
            padding: 20px; 
        }
        .container { 
            max-width: 650px; 
            margin: 0 auto; 
            background: white; 
            border: 1px solid #e0e0e0; 
            border-radius: 4px; 
        }
        .header { 
            background: #2c3e50; 
            color: white; 
            padding: 25px 30px; 
            border-bottom: 3px solid #34495e; 
        }
        .header h1 { 
            margin: 0 0 8px 0; 
            font-size: 22px; 
            font-weight: 600; 
        }
        .header p { 
            margin: 0; 
            font-size: 14px; 
            opacity: 0.9; 
        }
        .content { 
            padding: 30px; 
        }
        .section { 
            margin-bottom: 25px; 
            padding-bottom: 20px; 
            border-bottom: 1px solid #e8e8e8; 
        }
        .section:last-child { 
            border-bottom: none; 
        }
        .section-title { 
            font-size: 16px; 
            font-weight: 600; 
            margin-bottom: 15px; 
            color: #2c3e50; 
        }
        .info-row { 
            margin: 10px 0; 
            font-size: 14px; 
        }
        .label { 
            display: inline-block; 
            width: 120px; 
            font-weight: 500; 
            color: #666; 
        }
        .value { 
            color: #333; 
        }
        .footer { 
            background: #f8f9fa; 
            padding: 20px 30px; 
            text-align: center; 
            font-size: 12px; 
            color: #999; 
            border-top: 1px solid #e8e8e8; 
        }
        .footer p { 
            margin: 5px 0; 
        }
        .badge { 
            display: inline-block; 
            padding: 4px 10px; 
            background: #f0f0f0; 
            color: #666; 
            border-radius: 3px; 
            font-size: 13px; 
            font-weight: 500; 
        }
        .success-box { 
            background: #f8f9fa; 
            border-left: 3px solid #7f8c8d; 
            padding: 15px; 
            margin-top: 20px; 
        }
        .success-box p { 
            margin: 0; 
            color: #555; 
            font-size: 14px; 
        }
        .text-content {
            color: #555;
            font-size: 14px;
            line-height: 1.6;
            margin: 0;
            padding: 10px 0;
        }
        table { 
            width: 100%; 
            border-collapse: collapse; 
            margin-top: 12px; 
            font-size: 14px; 
        }
        th, td { 
            padding: 12px 10px; 
            text-align: left; 
            border-bottom: 1px solid #e8e8e8; 
        }
        th { 
            background-color: #f8f9fa; 
            font-weight: 600; 
            color: #2c3e50; 
        }
        tbody tr:hover { 
            background-color: #fafafa; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>專案報備通知 - 第二階段完成</h1>
            <p>專案交貨資訊已完成填寫</p>
        </div>
        
        <div class="content">
            <div class="section">
                <div class="section-title">專案基本資訊</div>
                <div class="info-row">
                    <span class="label">專案編號</span>
                    <span class="value">{{.ProjectID}}</span>
                </div>
                <div class="info-row">
                    <span class="label">專案名稱</span>
                    <span class="value">{{.ProjectName}}</span>
                </div>
                <div class="info-row">
                    <span class="label">聯絡人</span>
                    <span class="value">{{.ContactName}}</span>
                </div>
                <div class="info-row">
                    <span class="label">更新時間</span>
                    <span class="value">{{.UpdatedTime.Format "2006-01-02 15:04:05"}}</span>
                </div>
                <div class="info-row">
                    <span class="label">目前狀態</span>
                    <span class="badge">第二階段完成</span>
                </div>
            </div>

            <div class="section">
                <div class="section-title">交貨資訊</div>
                <div class="info-row">
                    <span class="label">預計交貨期</span>
                    <span class="value">{{.ExpectedDeliveryPeriod}}</span>
                </div>
                <div class="info-row">
                    <span class="label">預計交貨日</span>
                    <span class="value">{{.ExpectedDeliveryDate}}</span>
                </div>
                <div class="info-row">
                    <span class="label">預計履約期</span>
                    <span class="value">{{.ExpectedContractPeriod}}</span>
                </div>
                {{if .ContractStartDate}}
                <div class="info-row">
                    <span class="label">履約開始日</span>
                    <span class="value">{{.ContractStartDate}}</span>
                </div>
                {{end}}
                {{if .ContractEndDate}}
                <div class="info-row">
                    <span class="label">履約結束日</span>
                    <span class="value">{{.ContractEndDate}}</span>
                </div>
                {{end}}
            </div>

            {{if .Equipments}}
            <div class="section">
                <div class="section-title">設備清單</div>
                <table>
                    <thead>
                        <tr>
                            <th>料號</th>
                            <th>數量</th>
                            <th>說明</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Equipments}}
                        <tr>
                            <td>{{.PartNumber}}</td>
                            <td>{{.Quantity}}</td>
                            <td>{{if .Description}}{{.Description}}{{else}}-{{end}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{end}}

            {{if .DeliveryAddress}}
            <div class="section">
                <div class="section-title">交貨地址</div>
                <p class="text-content">{{.DeliveryAddress}}</p>
            </div>
            {{end}}

            {{if .SpecialRequirements}}
            <div class="section">
                <div class="section-title">特殊需求</div>
                <p class="text-content">{{.SpecialRequirements}}</p>
            </div>
            {{end}}

            <div class="success-box">
                <p><strong>完成：</strong>專案第二階段交貨資訊已完整填寫。</p>
            </div>
        </div>
        
        <div class="footer">
            <p>此為系統自動發送的通知郵件，請勿直接回覆</p>
            <p>專案報備系統 © 2025</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
