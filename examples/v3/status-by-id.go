// StatusById - وضعیت پیامک با شناسه‌هایی که متد ارسال برگردانده است.
//
// دسته‌ای بپرسید، نه یکی‌یکی. فاصله استعلام‌ها را هم کمتر از چند دقیقه
// نگذارید، وگرنه به خطای ۲۰ می‌خورید.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... go run examples/v3/status-by-id.go

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

// شرط را روی StatusCode بگذارید، نه روی متن Status. این پنج کد یعنی هنوز در
// راه است و باید بعداً دوباره استعلام کنید، نه اینکه دوباره بفرستید.
var pending = map[string]bool{"0": true, "1": true, "2": true, "3": true, "10": true}

func main() {
	payload, _ := json.Marshal(map[string]any{
		"ApiKey": os.Getenv("PAYAM_RESAN_API_KEY"),
		"Ids":    []int64{9903211, 9903212},
	})

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Post(
		"https://api.sms-webservice.com/api/V3/StatusById",
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

		again := ""
		if pending[fmt.Sprint(message["StatusCode"])] {
			again = " (بعداً دوباره بپرسید)"
		}

		fmt.Printf("%v: %v%s\n", message["Id"], message["Status"], again)
	}
}

// docs:end
