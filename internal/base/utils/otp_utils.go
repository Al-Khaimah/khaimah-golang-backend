package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	redisClient "github.com/Al-Khaimah/khaimah-golang-backend/internal/base/redis"
	"github.com/redis/go-redis/v9"
)

const (
	OTPLength     = 4
	OTPTTLMinutes = 5
	OTPKeyPrefix  = "otp:"
)

// GenerateOTP generates a random 4-digit OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(9000) + 1000
	return fmt.Sprintf("%04d", otp)
}

// StoreOTP stores the OTP in Redis with a TTL
func StoreOTP(ctx context.Context, identifier, otp string) error {
	key := fmt.Sprintf("%s%s", OTPKeyPrefix, identifier)
	return redisClient.SetWithTTL(ctx, key, otp, OTPTTLMinutes*time.Minute)
}

// VerifyOTP verifies the OTP against the stored value
func VerifyOTP(ctx context.Context, identifier, otp string) (bool, error) {
	key := fmt.Sprintf("%s%s", OTPKeyPrefix, identifier)
	storedOTP, err := redisClient.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return storedOTP == otp, nil
}

// DeleteOTP deletes the OTP from Redis
func DeleteOTP(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("%s%s", OTPKeyPrefix, identifier)
	return redisClient.Delete(ctx, key)
}

// SendEmailOTP sends an OTP via email using Resend API
func SendEmailOTP(email, otp string) error {
	apiKey := config.GetEnv("RESEND_API_KEY")
	senderEmail := config.GetEnv("RESEND_SENDER_EMAIL")
	senderName := "الخيمة"
	endpoint := "https://api.resend.com/emails"
	user, err := NewUserRepository(config.GetDB()).FindOneByEmail(email)
	if err != nil {
		fmt.Errorf("error finding user: %w", err)
	}

	firstName := ""
	if user != nil {
		firstName = user.FirstName
	}
	payload := map[string]interface{}{
		"from":    fmt.Sprintf("%s <%s>", senderName, senderEmail),
		"to":      []string{email},
		"subject": fmt.Sprintf("رمز التحقق هو %s", otp),
		"html": fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="ar" dir="rtl">
		<head>
			<meta charset="UTF-8">
			<title>رمز التحقق</title>
			<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans+Arabic:wght@400;700&display=swap" rel="stylesheet">
			<style>
				body { font-family: 'IBM Plex Sans Arabic', Tahoma, Arial, sans-serif; background: #fff; color: #1b1b1b; }
				.btn {
					display: inline-block;
					padding: 8px 16px;
					background: #098445;
					color: #fff;
					text-decoration: none;
					border-radius: 5px;
					font-weight: bold;
				}
			</style>
		</head>
		<body>
			<h2 style="color: #3D8361;">أرحب يا %s 👋</h2>
			<p>يا هلا بك، نورت الخيمة والله!✨</p>
			<p style="font-size: 22px; letter-spacing: 2px; font-weight: bold;">
				رمز التحقق حقك (شفه فوق كل شي):<br>
				<span style="background:#f4f4f4; border-radius:7px; padding:7px 15px; color:#098445;">%s</span>
			</p>
			<br>
			<p>
				<span style="font-size:17px; color:#9b6a18;">
					<br>عندك خوي مسوي مشغول وما عنده وقت يقرا؟ او ما يحب تويتر؟ أو شايب الجرايد معد صاروا يوصلون له؟ شاركهم الرسالة، يحملون التطبيق ويعيشون معنا الأجواء!
					شاركهم التطبيق يحملونه واستمتعوا سوا. واذا قدرت تعطينا تقييم حلو بالاب ستور، تسعدنا! 🌟
				</span>
			</p>
			<p>
				<a class="btn" href="https://apps.apple.com/sa/app/id6745527443">حمل تطبيق الخيمة من هنا</a>
			</p>
			<br>
			<p>
				واجهتك مشكلة؟ عندك سؤال؟ تواصل معنا على الواتساب <b>0591434366</b>، وحنّا بالخدمة دايمًا!
			</p>
			<p style="color:#aaa; font-size:12px;">ودنا نسمع منك، فريق الخيمة 🤠</p>
		</body>
		</html>
	`, firstName, otp),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send OTP via email: status code %d", resp.StatusCode)
	}

	return nil
}

// SendMobileOTP sends an OTP via WhatsApp using SendMsg.dev
func SendMobileOTP(mobile, otp string) error {
	formattedMobile := FormatMobileNumber(mobile)

	apiToken := config.GetEnv("SENDMSG_API_TOKEN")
	endpoint := "https://sendmsg.dev/message"
	user, err := NewUserRepository(config.GetDB()).FindOneByMobile(formattedMobile)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	firstName := ""
	if user != nil {
		firstName = user.FirstName
	}

	payload := map[string]interface{}{
		"to": []string{formattedMobile},
		"message": fmt.Sprintf(
			`رمز التحقق حقك: %s

			أهلًا يا %s، نورت الخيمة! 🌵✨

			حسابك جاهز، تقدر تبدأ تستمع للبودكاستات وتعيش الجو.

			عندك خوي مسوي مشغول وما عنده وقت يقرا؟ او ما يحب تويتر؟ أو شايب الجرايد معد صاروا يوصلون له؟ شاركهم الرسالة، يحملون التطبيق ويعيشون معنا الأجواء!

			وإذا عجبك التطبيق، لا تنسى تعطينا تقييم في الاب ستور:
			https://apps.apple.com/sa/app/id6745527443

			أي استفسار أو واجهتك مشكلة؟ كلمنا مباشرة على هالواتساب: 0591434366 (وتقدر ترد على نفس الرسالة).`,
			otp,
			firstName,
		),
		"token": apiToken,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		bodyString := ""
		if readErr == nil {
			bodyString = string(bodyBytes)
		}
		return fmt.Errorf("failed to send OTP via WhatsApp: status code %d - [%s]", resp.StatusCode, bodyString)
	}

	return nil
}
