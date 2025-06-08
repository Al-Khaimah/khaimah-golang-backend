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
    <title>رمز التحقق - الخيمة</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body, table, td, p, a, li, blockquote {
            -webkit-text-size-adjust: 100%%;
            -ms-text-size-adjust: 100%%;
            font-family: 'IBM Plex Sans Arabic', Tahoma, Arial, sans-serif;
        }
        body {
            direction: rtl;
            text-align: right;
            background-color: #f5f5f5;
            margin: 0;
            padding: 0;
        }
        .main-table {
            background-color: #fff;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin: 0 auto;
        }
        .header {
            color: #C13144;
            font-size: 32px;
            font-weight: bold;
            text-align: center;
            border-bottom: 3px solid #C13144;
            padding: 30px 30px 20px 30px;
            font-family: 'IBM Plex Sans Arabic', Tahoma, Arial, sans-serif;
        }
        .welcome {
            color: #619781;
            font-size: 24px;
            font-weight: bold;
            margin: 0 0 15px 0;
        }
        .desc {
            color: #323334;
            font-size: 16px;
            margin: 0 0 20px 0;
            line-height: 1.6;
        }
        .otp-box {
            background-color: #fff;
            border: 2px solid #C13144;
            border-radius: 6px;
            padding: 15px 25px;
            display: inline-block;
            margin: 10px 0;
            color: #C13144;
            font-size: 28px;
            font-weight: bold;
            letter-spacing: 3px;
            font-family: 'Courier New', monospace;
        }
        .section {
            background-color: #f8f9fa;
            border-radius: 8px;
            margin-bottom: 16px;
            padding: 25px;
        }
        .share {
            color: #619781;
            font-size: 17px;
        }
        .download-btn {
            display: inline-block;
            padding: 15px 30px;
            background-color: #C13144;
            color: #fff;
            text-decoration: none;
            border-radius: 8px;
            font-weight: bold;
            font-size: 16px;
            border: none;
            margin: 25px 0 0 0;
        }
        .contact-section {
            background-color: #323334;
            border-radius: 6px;
            text-align: center;
            padding: 20px;
        }
        .contact-phone {
            background-color: #fff;
            color: #C13144;
            padding: 4px 8px;
            border-radius: 4px;
            font-weight: bold;
            display: inline-block;
        }
        .footer {
            text-align: center;
            padding: 20px 30px 30px 30px;
            border-top: 1px solid #e9ecef;
        }
        .footer-title {
            color: #619781;
            font-size: 16px;
            font-weight: bold;
        }
        .footer-contact {
            color: #6c757d;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <table width="100%%" style="background-color: #f5f5f5; direction:rtl;">
        <tr>
            <td style="padding: 20px 0;">
                <table width="600" class="main-table">
                    <!-- Header -->
                    <tr>
                        <td class="header">الخيمة</td>
                    </tr>
                    
                    <!-- Welcome -->
                    <tr>
                        <td style="padding: 30px;">
                            <h2 class="welcome">أرحب يا %s 👋</h2>
                            <p class="desc">يا هلا بك، <span style="color: #C13144; font-weight: bold;"> تو ما نورت الخيمة</span> والله! ⛺</p>
                        </td>
                    </tr>
                    
                    <!-- OTP -->
                    <tr>
                        <td style="padding: 0 30px 30px 30px;">
                            <div class="section" style="text-align: center;">
                                <p style="margin: 0 0 15px 0; color: #323334; font-size: 18px; font-weight: bold;">رمز التحقق حقك:</p>
                                <div class="otp-box">%s</div>
                                <p style="margin: 15px 0 0 0; color: #619781; font-size: 14px;">استخدم هذا الرمز لتفعيل حسابك</p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Share Section -->
                    <tr>
                        <td style="padding: 0 30px 30px 30px;">
                            <div class="section">
							<p class="share">
								عندك خوي مسوي مشغول وما عنده وقت يقرا؟ 🤷‍♂️<br>
								أو ما يحب تويتر؟ 🐦🚫<br>
								أو شايب الجرايد معد صاروا يوصلون له؟ 👴📰<br>
								<br>
								<br>
								<span style="color: #C13144; font-weight: bold;">شاركهم التطبيق</span> وخلهم يسمعون الأخبار اللي تهمهم بضغطة زر وحده!<br>
								<br>
								<br>
								إذا جازلتلك الخيمة، قيمنا في الاب ستور ❤️🌟
							</p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Download Button -->
                    <tr>
                        <td style="padding: 0 30px 30px 30px; text-align: center;">
                            <a class="download-btn" href="https://apps.apple.com/sa/app/id6745527443">
                                📱 حمل تطبيق الخيمة من هنا
                            </a>
                        </td>
                    </tr>
                    
                    <!-- Contact Section -->
                    <tr>
                        <td style="padding: 0 30px 30px 30px;">
                            <div class="contact-section">
                                <p style="margin: 0; color: #fff; font-size: 16px; line-height: 1.6;">
                                    <strong>واجهتك مشكلة؟ عندك سؤال؟</strong><br>
                                    تواصل معنا على الواتساب <span class="contact-phone">0591434366</span><br>
                                    <span style="color: #C13144; font-weight: bold;">وحنّا بالخدمة دايمًا!</span>
                                </p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td class="footer">
                            <p class="footer-title">ودنا نسمع منك، <span style="color: #C13144;">فريق الخيمة</span> 🤠</p>
                            <p class="footer-contact">
                                AlKhimaPlatform@outlook.com | 0506054839
                            </p>
                        </td>
                    </tr>
                    
                </table>
            </td>
        </tr>
    </table>
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
		"to":      []string{formattedMobile},
		"message": fmt.Sprintf("رمز التحقق حقك: %s\n\nأرحب %s، تو ما نورت الخيمة! ⛺ \n\nحسابك جاهز، تقدر تبدأ تستمع للبودكاستات وتعيش الجو.\n\nعندك خوي مسوي مشغول وما عنده وقت يقرا؟ 🤷‍♂️\nأو ما يحب تويتر؟ 🐦🚫\nأو شايب الجرايد معد صاروا يوصلون له؟ 👴📰\n\nشاركهم التطبيق وخلهم يسمعون الأخبار اللي تهمهم بضغطة زر وحده!\n\nإذا جازلتلك الخيمة، قيمنا في الاب ستور ❤️🌟\n:https://apps.apple.com/sa/app/id6745527443\n\nأي استفسار أو واجهتك مشكلة؟ كلمنا مباشرة على هالواتساب: 0591434366 (وتقدر ترد على نفس الرسالة).", otp, firstName),
		"token":   apiToken,
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
