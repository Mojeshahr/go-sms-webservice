<div align="center">

<a href="https://payam-resan.com">
  <img src=".github/assets/logo.svg" width="64" height="64" alt="پیام رسان">
</a>

<h1>نمونه‌کدهای Go وب‌سرویس پیام رسان</h1>

اتصال به وب‌سرویس <a href="https://payam-resan.com"><b>پنل پیامکی پیام رسان</b></a> با Go<br>
یک فایل قابل اجرا به‌ازای هر متد سرویس، بدون هیچ وابستگی

[![API](https://img.shields.io/badge/API-V3-0a7cbd)](https://payam-resan.com)
[![Go](https://img.shields.io/badge/Go-1.21%2B-00add8)](https://go.dev)
[![Dependencies](https://img.shields.io/badge/dependencies-none-2ea44f)](#شروع-سریع)
[![License](https://img.shields.io/badge/license-MIT-6e7781)](LICENSE)

<b>فارسی</b> · <a href="README.en.md">English</a>

</div>

<sub>دنبال زبان دیگری هستید؟ همین نمونه‌ها برای زبان‌های دیگر هم در
[github.com/Mojeshahr](https://github.com/Mojeshahr) هست.</sub>

---

## شروع سریع

```bash
git clone https://github.com/Mojeshahr/go-sms-webservice.git
cd go-sms-webservice

export PAYAM_RESAN_API_KEY='123456-XXXXXXXXXXXXXXX'
export PAYAM_RESAN_SENDER='30004040'

go run examples/v3/account-info.go
```

نه `go mod` لازم است و نه هیچ بسته‌ای. `net/http` و `encoding/json` هر دو در
کتابخانه استاندارد Go هستند.

با `account-info.go` شروع کنید: چیزی ارسال نمی‌کند، اعتباری مصرف نمی‌کند، و اگر
جواب داد یعنی کلید و اتصال هر دو سالم‌اند.

## پیش از ارسال واقعی

یک سرور آزمایشی هست که مثل سرور عملیاتی جواب می‌دهد ولی پیامکی نمی‌فرستد و
اعتباری مصرف نمی‌کند. کافی است `V3` در نشانی را با `V3SandBox` عوض کنید. تنها
استثنا `TokenList` است که روی آن سرور پیاده نشده.

## متدها

<div dir="rtl">

| نمونه | متد | کار |
|---|---|---|
| [account-info.go](examples/v3/account-info.go) | `AccountInfo` | اعتبار و خطوط فعال |
| [send.go](examples/v3/send.go) | `Send` | ارسال ساده با `GET` |
| [send-bulk.go](examples/v3/send-bulk.go) | `SendBulk` | یک متن به چند گیرنده، با شناسه پی‌گیری |
| [send-multiple.go](examples/v3/send-multiple.go) | `SendMultiple` | متن جدا برای هر گیرنده |
| [token-list.go](examples/v3/token-list.go) | `TokenList` | فهرست قالب‌ها |
| [send-token-single.go](examples/v3/send-token-single.go) | `SendTokenSingle` | ارسال قالب به یک شماره |
| [send-token-single-get.go](examples/v3/send-token-single-get.go) | `SendTokenSingle` | همان، با `GET` |
| [send-token-multi.go](examples/v3/send-token-multi.go) | `SendTokenMulti` | یک قالب، چند گیرنده |
| [status-by-id.go](examples/v3/status-by-id.go) | `StatusById` | وضعیت با شناسه سامانه |
| [status-by-user-trace-id.go](examples/v3/status-by-user-trace-id.go) | `StatusByUserTraceId` | وضعیت با شناسه خودتان |
| [get-inbox.go](examples/v3/get-inbox.go) | `GetInbox` | پیامک‌های رسیده |

</div>

## چرا go.mod نیست

هر فایل یک `package main` مستقل است و با `go run <file>` اجرا می‌شود. چند فایل
با `func main()` در یک پوشه مشکلی ندارند، چون `go run` فقط همان فایلی را که
نام برده‌اید کامپایل می‌کند.

اگر پروژه خودتان ماژول دارد، بدنه `main` را داخل تابع خودتان کپی کنید؛ کد بدون
تغییر کار می‌کند.

## چند نکته که وقت‌تان را می‌خرد

**`UseNumber` را حذف نکنید.** Go هر عدد JSON را به `float64` تبدیل می‌کند و
شناسه‌های بلند سامانه آن وقت به شکل نمایی چاپ می‌شوند، مثل
`6.8733095017201e+14` به‌جای `687330950172013`. با `UseNumber` هر عدد همان‌طور
که سرویس فرستاده می‌ماند.

**کد وضعیت HTTP را نخوانید.** سرویس همیشه `200` برمی‌گرداند، حتی وقتی کلید
اشتباه است، پس `answer.StatusCode` چیزی ثابت نمی‌کند. تصمیم را از فیلد `Success`
بگیرید.

**شماره گیرنده صفر ابتدایی ندارد.** یعنی `9121112222` یا با کد کشور
`989121112222`. شماره‌ای که با `9` یا `989` شروع نشود کد خطای `13` می‌گیرد.

**متن را دوباره encode نکنید.** در `send.go` متد `Encode` روی `url.Values` خودش
یک بار این کار را می‌کند. اگر پیش از آن هم `QueryEscape` کنید، پیامک با
نویسه‌های `%D8` به گوشی می‌رسد.

**برای هر گیرنده یک `UserTraceId` یکتا بفرستید.** بعد از یک timeout، این تنها
راه فهمیدن این است که پیامک ثبت شده یا نه.

## امنیت کلید

کلید یک راز است. در مخزن کد، در جاوااسکریپت مرورگر و در بسته اپلیکیشن موبایل
نباید قرار بگیرد. جای آن متغیر محیطی است، همان‌طور که همه نمونه‌ها می‌خوانندش.

اگر کلیدی لو رفت، از پنل یکی تازه بسازید. کلید حذف‌شده برنمی‌گردد.

## ساختار

<div dir="rtl">

| مسیر | چه چیزی دارد |
|---|---|
| `examples/v3/` | یک نمونه مستقل به‌ازای هر عملیات سرویس |
| `.env.example` | نمونه متغیرهای محیطی |

</div>

عدد `v3` در مسیر عمدی است. نسخه تازه سرویس یعنی پوشه `examples/v<n>/` تازه، و
پوشه موجود دست‌نخورده می‌ماند.

## مستندات و پشتیبانی

راهنمای کامل وب‌سرویس در [docs.payam-resan.com](https://docs.payam-resan.com)
است. توصیف ماشین‌خوان OpenAPI هم در
[sms-webservice-spec](https://github.com/Mojeshahr/sms-webservice-spec).

## مجوز

MIT. متن کامل در [`LICENSE`](LICENSE).
