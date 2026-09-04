// Send - ساده‌ترین ارسال، یک متن به چند شماره با یک درخواست GET.
//
// برای آزمایش سریع خوب است. در محیط عملیاتی SendBulk را بردارید: کلید را از
// نشانی بیرون می‌برد و برای هر گیرنده شناسه پی‌گیری می‌پذیرد.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... PAYAM_RESAN_SENDER=... go run examples/v3/send.go

// docs:start
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	// Encode دقیقاً یک بار encode می‌کند. اگر متن را خودتان هم پیش از این
	// QueryEscape کنید، پیامک با نویسه‌های %D8 به گوشی می‌رسد.
	query := url.Values{}
	query.Set("ApiKey", os.Getenv("PAYAM_RESAN_API_KEY"))
	query.Set("Sender", os.Getenv("PAYAM_RESAN_SENDER"))
	query.Set("Text", "کد تأیید شما ۱۲۳۴۵۶ است")
	query.Set("Recipients", "9121112222,9121113333")

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Get("https://api.sms-webservice.com/api/V3/Send?" + query.Encode())
	if err != nil {
		fmt.Fprintln(os.Stderr, "خطای شبکه:", err)
		os.Exit(1)
	}
	defer answer.Body.Close()

	decoder := json.NewDecoder(answer.Body)
	decoder.UseNumber()

	var response map[string]any
	if err := decoder.Decode(&response); err != nil {
		fmt.Fprintln(os.Stderr, "پاسخ خوانده نشد:", err)
		os.Exit(1)
	}

	if response["Success"] != true {
		fmt.Fprintf(os.Stderr, "ناموفق. کد %v: %v\n", response["ErrorCode"], response["Error"])
		os.Exit(1)
	}

	result, _ := response["Result"].([]any)
	for _, item := range result {
		message, _ := item.(map[string]any)
		fmt.Println("شناسه", message["Id"])
	}
}

// docs:end
