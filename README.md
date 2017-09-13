Wrapper of SDKs(SendGrid, MailGun, Gmail) Mailer
---
[ ![Codeship Status for thiagozs/sendmail-poc-go](https://app.codeship.com/projects/ba8d82b0-7945-0135-b5d7-62a7bc934352/status?branch=master)](https://app.codeship.com/projects/244971) [![Go Report Card](https://goreportcard.com/badge/github.com/thiagozs/sendmail-poc-go)](https://goreportcard.com/report/github.com/thiagozs/sendmail-poc-go)

- Go 1.9
- SDK Sendgrid
- SDK Mailgun
- SDK Gmail
- Glide package manager
- Linux

Just get the package.
```sh
go get -v github.com/thiagozs/mailer-go
```

To run the study just execute the `env.sh` file to setup the path and exports.
```sh
$. ./env.sh
```

In your `{your-project}/src` run the `glide` to download the dependencies of project.
```sh
.../src$ glide up
```

On root of project, create the file `.env` with the contents:
```sh
SENDGRID_API_KEY=...
MAILGUN_DOMAIN=...
MAILGUN_API_KEY=...
MAILGUN_PUB_KEY=...
GMAIL_LOGIN=...
GMAIL_PASSWD=...
EMAILFROM=...
EMAILFROMNAME=...
EMAILTO=...
EMAILTONAME=...
```

The project is open for any upgrade / commit, submit your proof of concept (POC).

**Simple example of code implementation**.
```go
package main

import (
	"fmt"
	mailer "github.com/thiagozs/mailer-go"
	"time"
)

func main() {

	//Sendgrid Instance
	sg := mailer.NewMailerSendGrid(os.Getenv("SENDGRID_API_KEY"))
	sg.Delay = time.Second * 5

	//Mailgun Instance
	mg := mailer.NewMailerMailGun(os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
		os.Getenv("MAILGUN_PUB_KEY"))
	mg.Delay = time.Second * 5

	//Gmail Instance
	gm := mailer.NewMailerGmail(os.Getenv("GMAIL_LOGIN"),
		os.Getenv("GMAIL_PASSWD"))
	gm.Delay = time.Second * 5

	// --------------

	// Sendgrid configs
	sg.ConfigEmail = mailer.ConfigEmailSendgrid{
		ContentHTML:      "<b>Test</b>",
		ContentPlainText: "Test",
		EmailFromName:    "Thiago Zilli",
		EmailFrom:        "yourmail@host.com",
		EmailToName:      "Client Name here",
		EmailTo:          "emailofclient@host.com",
		Subject:          "Test email",
	}

	// Mailgun configs
	mg.ConfigEmail = mailer.ConfigEmailMailGun{
		ContentPlainText: "Test",
		EmailFrom:        "yourmail@host.com",
		EmailTo:          "emailofclient@host.com",
		Subject:          "Test email",
	}

	// Gmail Configs
	gm.ConfigEmail = mailer.ConfigEmailGmail{
		ContentHTML: "<b>Test with Gmail</b>",
		EmailFrom:   "yourmail@host.com",
		EmailTo:     "emailofclient@host.com",
		Subject:     "Test email",
	}

	// --------------

	// Send email by Sendgrid
	if err := sg.SendMail(); err != nil {
		fmt.Printf("Sendgrid Error: %s\n", err.Error())
	}

	// Send email by Mailgun
	if err := mg.SendMail(); err != nil {
		fmt.Printf("Mailgun Error: %s\n", err.Error())
	}

	// Send email by Gmail
	if err := gm.SendMail(); err != nil {
		fmt.Printf("Mailgun Error: %s\n", err.Error())
	}

}
``` 

ToDos
---
- [x] Wrapper Sendgrid
- [x] Wrapper Mailgun
- [x] Wrapper Gmail
- [X] Tests

Autors
---
- Thiago Z S - @thiagozs (Github) / @thiagozs (Twitter)

