// SendTokenMulti - یک قالب، چند گیرنده، مقادیر متفاوت.
//
// پارامترها اینجا آرایه‌اند، نه p1 تا p10. درایه اول به {1} می‌نشیند، دومی به
// {2} و همین‌طور تا آخر: ترتیب از شماره جای‌گاه می‌آید، نه از جایی که در متن
// قالب دیده می‌شود.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... go run examples/v3/send-token-multi.go

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
	// قالب نمونه: «مرسوله شما از {2} تحویل پست شد. بارکد مرسوله پستی: {1}»
	payload, _ := json.Marshal(map[string]any{
		"ApiKey":      os.Getenv("PAYAM_RESAN_API_KEY"),
		"TemplateKey": "postcode",
		"Recipients": []map[string]any{
			{
				"Destination": 9121112222,
				"UserTraceId": 1001,
				"Parameters":  []string{"BARCODE-AAA", "شیراز"},
			},
			{
				"Destination": 9121113333,
				"UserTraceId": 1002,
				"Parameters":  []string{"BARCODE-BBB", "تبریز"},
			},
		},
	})

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Post(
		"https://api.sms-webservice.com/api/V3/SendTokenMulti",
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
		fmt.Printf("%v => %v\n", message["UserTraceId"], message["FinalText"])
	}
}

// docs:end
