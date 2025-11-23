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
	SendVerificationCode(email, code, username string) error // 新增
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
	Equipments             []Equipment
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
		smtpPort:     587,
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
        body { font-family: 'Microsoft JhengHei', Arial, sans-serif; line-height: 1.8; color: #333; background: #f5f5f5; margin: 0; padding: 20px; }
        .container { max-width: 650px; margin: 0 auto; background: white; border: 1px solid #e0e0e0; border-radius: 4px; }
        .header { background: #2c3e50; color: white; padding: 25px 30px; border-bottom: 3px solid #34495e; }
        .header h1 { margin: 0 0 8px 0; font-size: 22px; font-weight: 600; }
        .content { padding: 30px; }
        .section { margin-bottom: 25px; padding-bottom: 20px; border-bottom: 1px solid #e8e8e8; }
        .section-title { font-size: 16px; font-weight: 600; margin-bottom: 15px; color: #2c3e50; }
        .info-row { margin: 10px 0; font-size: 14px; }
        .label { display: inline-block; width: 120px; font-weight: 500; color: #666; }
        .footer { background: #f8f9fa; padding: 20px 30px; text-align: center; font-size: 12px; color: #999; border-top: 1px solid #e8e8e8; }
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
                <div class="info-row"><span class="label">專案編號</span>{{.ProjectID}}</div>
                <div class="info-row"><span class="label">專案名稱</span>{{.ProjectName}}</div>
                <div class="info-row"><span class="label">建立時間</span>{{.CreatedTime.Format "2006-01-02 15:04:05"}}</div>
            </div>
            <div class="section">
                <div class="section-title">聯絡人資訊</div>
                <div class="info-row"><span class="label">聯絡人</span>{{.ContactName}}</div>
                <div class="info-row"><span class="label">聯絡電話</span>{{.ContactPhone}}</div>
                <div class="info-row"><span class="label">聯絡信箱</span>{{.ContactEmail}}</div>
            </div>
        </div>
        <div class="footer"><p>此為系統自動發送的通知郵件，請勿直接回覆</p></div>
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
        body { font-family: 'Microsoft JhengHei', Arial, sans-serif; line-height: 1.8; color: #333; background: #f5f5f5; margin: 0; padding: 20px; }
        .container { max-width: 650px; margin: 0 auto; background: white; border: 1px solid #e0e0e0; border-radius: 4px; }
        .header { background: #2c3e50; color: white; padding: 25px 30px; border-bottom: 3px solid #34495e; }
        .header h1 { margin: 0 0 8px 0; font-size: 22px; font-weight: 600; }
        .content { padding: 30px; }
        .section { margin-bottom: 25px; padding-bottom: 20px; border-bottom: 1px solid #e8e8e8; }
        .section-title { font-size: 16px; font-weight: 600; margin-bottom: 15px; color: #2c3e50; }
        .info-row { margin: 10px 0; font-size: 14px; }
        .label { display: inline-block; width: 120px; font-weight: 500; color: #666; }
        .footer { background: #f8f9fa; padding: 20px 30px; text-align: center; font-size: 12px; color: #999; border-top: 1px solid #e8e8e8; }
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
                <div class="info-row"><span class="label">專案編號</span>{{.ProjectID}}</div>
                <div class="info-row"><span class="label">專案名稱</span>{{.ProjectName}}</div>
            </div>
            <div class="section">
                <div class="section-title">交貨資訊</div>
                <div class="info-row"><span class="label">預計交貨期</span>{{.ExpectedDeliveryPeriod}}</div>
                <div class="info-row"><span class="label">預計交貨日</span>{{.ExpectedDeliveryDate}}</div>
            </div>
        </div>
        <div class="footer"><p>此為系統自動發送的通知郵件，請勿直接回覆</p></div>
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

// SendVerificationCode 發送登入驗證碼郵件
func (s *emailService) SendVerificationCode(email, code, username string) error {
	subject := "【登入驗證碼】專案報備系統"

	htmlBody := s.renderVerificationCodeTemplate(email, code, username)

	// 使用 gomail 發送郵件到指定信箱
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", email) // 發送到使用者的信箱
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUser, s.smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Error("Failed to send verification code email:", err)
		return fmt.Errorf("failed to send verification code email: %v", err)
	}

	log.Info("Verification code email sent successfully to:", email)
	return nil
}

// renderVerificationCodeTemplate 渲染驗證碼郵件範本
func (s *emailService) renderVerificationCodeTemplate(email, code, username string) string {
	return `
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
            max-width: 600px; 
            margin: 0 auto; 
            background: white; 
            border: 1px solid #e0e0e0; 
            border-radius: 8px; 
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white; 
            padding: 30px; 
            text-align: center;
        }
        .header h1 { 
            margin: 0 0 10px 0; 
            font-size: 24px; 
            font-weight: 600; 
        }
        .content { 
            padding: 40px 30px; 
            text-align: center;
        }
        .greeting {
            font-size: 18px;
            color: #2c3e50;
            margin-bottom: 30px;
        }
        .code-box {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 30px;
            border-radius: 8px;
            margin: 30px 0;
        }
        .code {
            font-size: 42px;
            font-weight: bold;
            color: white;
            letter-spacing: 8px;
            margin: 0;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }
        .code-label {
            color: rgba(255,255,255,0.9);
            font-size: 14px;
            margin-top: 15px;
        }
        .notice {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px 20px;
            margin: 30px 0;
            text-align: left;
            border-radius: 4px;
        }
        .notice-title {
            font-weight: 600;
            color: #856404;
            margin-bottom: 8px;
        }
        .notice-item {
            color: #856404;
            font-size: 14px;
            margin: 5px 0;
            padding-left: 20px;
        }
        .footer { 
            background: #f8f9fa; 
            padding: 20px 30px; 
            text-align: center; 
            font-size: 12px; 
            color: #999; 
            border-top: 1px solid #e8e8e8; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 登入驗證碼</h1>
            <p>專案報備系統</p>
        </div>
        
        <div class="content">
            <div class="greeting">
                <strong>` + username + `</strong>，您好！
            </div>
            
            <p style="color: #666; font-size: 15px; margin-bottom: 20px;">
                您正在嘗試登入專案報備系統<br>
                請使用以下驗證碼完成登入驗證：
            </p>

            <div class="code-box">
                <div class="code">` + code + `</div>
                <div class="code-label">請在登入頁面輸入此驗證碼</div>
            </div>

            <div class="notice">
                <div class="notice-title">⚠️ 重要提醒：</div>
                <div class="notice-item">• 驗證碼有效期為 <strong>5 分鐘</strong></div>
                <div class="notice-item">• 驗證碼僅可使用一次</div>
                <div class="notice-item">• 如非本人操作，請忽略此郵件</div>
            </div>

            <p style="color: #999; font-size: 13px; margin-top: 30px;">
                如果您沒有嘗試登入，請忽略此郵件。<br>
                為了您的帳號安全，請勿將驗證碼分享給任何人。
            </p>
        </div>
        
        <div class="footer">
            <p>此為系統自動發送的驗證郵件，請勿直接回覆</p>
            <p>專案報備系統 © 2025</p>
        </div>
    </div>
</body>
</html>
`
}
