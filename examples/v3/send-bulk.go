// SendBulk - یک متن به چند گیرنده، هر کدام با شناسه پی‌گیری خودتان.
//
// روش پیشنهادی برای ارسال عملیاتی. کلید در بدنه درخواست می‌رود نه در نشانی،
// و برای هر گیرنده UserTraceId می‌پذیرد تا گزارش تحویل را بدون نگه‌داشتن Id
// سامانه بگیرید.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... PAYAM_RESAN_SENDER=... go run examples/v3/send-bulk.go

// docs:start
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	sender, err := strconv.ParseInt(os.Getenv("PAYAM_RESAN_SENDER"), 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "PAYAM_RESAN_SENDER یک عدد نیست")
		os.Exit(1)
	}

	payload, _ := json.Marshal(map[string]any{
		"ApiKey": os.Getenv("PAYAM_RESAN_API_KEY"),
		"Sender": sender,
		"Text":   "سفارش شما ثبت شد.",
		"Recipients": []map[string]any{
			{"Destination": 9121112222, "UserTraceId": 1001},
			{"Destination": 9121113333, "UserTraceId": 1002},
		},
	})

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Post(
		"https://api.sms-webservice.com/api/V3/SendBulk",
		"application/json; charset=utf-8",
		bytes.NewReader(payload),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "خطای شبکه:", err)
		os.Exit(1)
	}
	defer answer.Body.Close()

	// UseNumber وگرنه هر عدد JSON به float64 تبدیل می‌شود و شناسه‌های بلند
	// سامانه به شکل نمایی چاپ می‌شوند.
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
		fmt.Printf("%v => شناسه %v\n", message["UserTraceId"], message["Id"])
	}
}

// docs:end
