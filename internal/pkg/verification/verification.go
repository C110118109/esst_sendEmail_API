package verification

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// VerificationCode 驗證碼結構
type VerificationCode struct {
	Code      string
	ExpiresAt time.Time
	Email     string
}

// VerificationService 驗證碼服務介面
type VerificationService interface {
	GenerateCode(email string) (string, error)
	VerifyCode(email, code string) bool
	CleanExpiredCodes()
}

type verificationService struct {
	codes map[string]*VerificationCode
	mu    sync.RWMutex
}

var (
	instance *verificationService
	once     sync.Once
)

// New 建立驗證碼服務單例
func New() VerificationService {
	once.Do(func() {
		instance = &verificationService{
			codes: make(map[string]*VerificationCode),
		}
		// 啟動清理過期驗證碼的 goroutine
		go instance.startCleanupRoutine()
	})
	return instance
}

// GenerateCode 生成 6 位數驗證碼
func (s *verificationService) GenerateCode(email string) (string, error) {
	// 生成 6 位數隨機驗證碼
	code := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n.Int64())
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 儲存驗證碼,有效期 5 分鐘
	s.codes[email] = &VerificationCode{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Email:     email,
	}

	return code, nil
}

// VerifyCode 驗證驗證碼
func (s *verificationService) VerifyCode(email, code string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	vc, exists := s.codes[email]
	if !exists {
		return false
	}

	// 檢查是否過期
	if time.Now().After(vc.ExpiresAt) {
		return false
	}

	// 驗證碼匹配
	if vc.Code == code {
		// 驗證成功後刪除驗證碼(一次性使用)
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.codes, email)
		s.mu.Unlock()
		s.mu.RLock()
		return true
	}

	return false
}

// CleanExpiredCodes 清理過期的驗證碼
func (s *verificationService) CleanExpiredCodes() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for email, vc := range s.codes {
		if now.After(vc.ExpiresAt) {
			delete(s.codes, email)
		}
	}
}

// startCleanupRoutine 啟動定期清理過期驗證碼的 goroutine
func (s *verificationService) startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.CleanExpiredCodes()
	}
}
