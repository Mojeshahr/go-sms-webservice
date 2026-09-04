// SendMultiple - متن و خط فرستنده جدا برای هر گیرنده.
//
// برای پیام‌های شخصی‌سازی‌شده که با یک قالب ثابت پوشش داده نمی‌شوند. برخلاف
// SendBulk، اینجا Text و Sender در سطح هر گیرنده تعریف می‌شوند.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... PAYAM_RESAN_SENDER=... go run examples/v3/send-multiple.go

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
		"Recipients": []map[string]any{
			{
				"Sender":      sender,
				"Destination": 9121112222,
				"Text":        "آقای محمدی، سفارش شما ارسال شد.",
				"UserTraceId": 1001,
			},
			{
				"Sender":      sender,
				"Destination": 9121113333,
				"Text":        "خانم رضایی، سفارش شما ارسال شد.",
				"UserTraceId": 1002,
			},
		},
	})

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Post(
		"https://api.sms-webservice.com/api/V3/SendMultiple",
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
		fmt.Printf("%v => شناسه %v\n", message["UserTraceId"], message["Id"])
	}
}

// docs:end
