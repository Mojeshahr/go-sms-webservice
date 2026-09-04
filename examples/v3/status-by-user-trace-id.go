// StatusByUserTraceId - وضعیت پیامک با شناسه‌هایی که خودتان داده‌اید.
//
// اگر UserTraceId را کلید رکورد پایگاه داده خودتان بگذارید، دیگر لازم نیست Id
// سامانه را ذخیره کنید. این متد راه امن تشخیص ارسال تکراری هم هست: بعد از قطع
// ارتباط، اول اینجا بپرسید ثبت شده یا نه.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... go run examples/v3/status-by-user-trace-id.go

// docs:start
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	payload, _ := json.Marshal(map[string]any{
		"ApiKey":       os.Getenv("PAYAM_RESAN_API_KEY"),
		"UserTraceIds": []int64{1001, 1002},
	})

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Post(
		"https://api.sms-webservice.com/api/V3/StatusByUserTraceId",
		"application/json; charset=utf-8",
		bytes.NewReader(payload),
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

		// کد ۸ یعنی این شناسه در حساب شما نیست. بعد از یک timeout، همین یعنی
		// ارسال ثبت نشده و می‌توانید با خیال راحت دوباره بفرستید.
		if fmt.Sprint(message["StatusCode"]) == "8" {
			fmt.Printf("%v: ثبت نشده\n", message["UserTraceId"])
			continue
		}

		fmt.Printf("%v: %v\n", message["UserTraceId"], message["Status"])
	}
}

// docs:end
