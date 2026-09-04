// SendTokenSingle - ارسال قالب به یک شماره، با بدنه JSON.
//
// مسیر معمول رمز یک‌بارمصرف. خط فرستنده ورودی ندارد؛ سامانه آن را از روی خود
// قالب برمی‌دارد. همین واریانت POST را به کار ببرید، نه GET: در GET هم کلید
// حساب و هم خود رمز داخل نشانی و لاگ وب‌سرور می‌نشینند.
//
// جز کتابخانه استاندارد Go به چیزی وابسته نیست. کپی کنید و در پروژه خودتان
// اجرا کنید.
//
//   PAYAM_RESAN_API_KEY=... go run examples/v3/send-token-single.go

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
		"ApiKey":      os.Getenv("PAYAM_RESAN_API_KEY"),
		"TemplateKey": "verifycode",
		"Destination": 9121112222,
		"p1":          "123456",
	})

	client := &http.Client{Timeout: 30 * time.Second}

	answer, err := client.Post(
		"https://api.sms-webservice.com/api/V3/SendTokenSingle",
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

	// این متد UserTraceId در ورودی ندارد، پس در پاسخ null برمی‌گردد. اگر شناسه
	// پی‌گیری لازم دارید، SendTokenMulti را حتی برای یک گیرنده هم می‌شود به کار برد.
	result, _ := response["Result"].([]any)
	for _, item := range result {
		message, _ := item.(map[string]any)
		fmt.Printf("شناسه %v از خط %v\n", message["Id"], message["Sender"])
		fmt.Println("متن نهایی:", message["FinalText"])
	}
}

// docs:end
