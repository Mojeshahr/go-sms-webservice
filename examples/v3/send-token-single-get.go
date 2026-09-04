// SendTokenSingle با GET - همان ارسال قالب، با ورودی در نشانی.
//
// برای آزمایش دستی مناسب است، برای محیط عملیاتی نه: در GET هم کلید حساب و هم
// مقدار رمز یک‌بارمصرف داخل نشانی می‌نشینند و در لاگ وب‌سرور و هدر Referer
// ثبت می‌شوند. واریانت POST را بردارید.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... go run examples/v3/send-token-single-get.go

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
	query := url.Values{}
	query.Set("ApiKey", os.Getenv("PAYAM_RESAN_API_KEY"))
	query.Set("TemplateKey", "verifycode")
	query.Set("Destination", "9121112222")
	query.Set("p1", "123456")

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Get(
		"https://api.sms-webservice.com/api/V3/SendTokenSingle?" + query.Encode(),
	)
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
		fmt.Printf("شناسه %v، متن نهایی: %v\n", message["Id"], message["FinalText"])
	}
}

// docs:end
