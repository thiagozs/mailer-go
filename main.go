package mailer

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// SDKConfigSengrid cfg SDKs
type SDKConfigSengrid struct {
	SendGridAPIKey string
	SDKName        string
	Delay          time.Duration
	ConfigEmail    ConfigEmailSendgrid
}

// SDKConfigMailGun cfg SDKs
type SDKConfigMailGun struct {
	MailGunDomain string
	MailGunAPIKey string
	MailGunPUBKey string
	SDKName       string
	Delay         time.Duration
	ConfigEmail   ConfigEmailMailGun
}

// SDKConfigGmail cfg SDKs
type SDKConfigGmail struct {
	User        string
	Password    string
	SDKName     string
	Delay       time.Duration
	ConfigEmail ConfigEmailGmail
}

// SDK wrapper of APIs
type SDK struct {
	Mailgun  mailgun.Mailgun
	Sendgrid *sendgrid.Client
	Gmail    GmailLogin
}

// GmailLogin username and password for gmail
type GmailLogin struct {
	User     string
	Password string
}

// ConfigEmailSendgrid configuration of send
type ConfigEmailSendgrid struct {
	EmailTo          string
	EmailFrom        string
	EmailToName      string
	EmailFromName    string
	ContentHTML      string
	ContentPlainText string
	Subject          string
}

// ConfigEmailMailGun configuration of send
type ConfigEmailMailGun struct {
	EmailTo          string
	EmailFrom        string
	ContentPlainText string
	Subject          string
}

// ConfigEmailGmail configuration of send
type ConfigEmailGmail struct {
	EmailTo     string
	EmailFrom   string
	ContentHTML string
	Subject     string
}

// newSDKSendgrid get a SDKs
func (cfg SDKConfigSengrid) newSDKSendgrid() *SDK {
	return &SDK{
		Sendgrid: sendgrid.NewSendClient(cfg.SendGridAPIKey),
	}
}

// newSDKMailGun get a SDKs
func (cfg SDKConfigMailGun) newSDKMailGun() *SDK {
	return &SDK{
		Mailgun: mailgun.NewMailgun(
			cfg.MailGunDomain,
			cfg.MailGunAPIKey,
			cfg.MailGunPUBKey,
		),
	}
}

// newSDKMailGun get a SDKs
func (cfg SDKConfigGmail) newSDKGmail() *SDK {
	return &SDK{
		Gmail: GmailLogin{
			cfg.User,
			cfg.Password,
		},
	}
}

// NewMailerSendGrid new instance of SendGrid
func NewMailerSendGrid(apikey string) *SDKConfigSengrid {
	return &SDKConfigSengrid{
		SendGridAPIKey: apikey,
		SDKName:        "sendgrid",
	}
}

// NewMailerMailGun new instance of MailGun
func NewMailerMailGun(domain string, apikey string, pubkey string) *SDKConfigMailGun {
	return &SDKConfigMailGun{
		MailGunDomain: domain,
		MailGunAPIKey: apikey,
		MailGunPUBKey: pubkey,
		SDKName:       "mailgun",
	}
}

// NewMailerGmail new instace for Gmail login
func NewMailerGmail(user string, password string) *SDKConfigGmail {
	return &SDKConfigGmail{
		User:     user,
		Password: password,
		SDKName:  "gmail",
	}
}

// SendMail sendemail
func (cfg *SDKConfigSengrid) SendMail() error {
	sdk := cfg.newSDKSendgrid()

	if CheckIsEmptyCfg(cfg) {
		return fmt.Errorf("Empty fields on ConfigEmail")
	}

	if cfg.Delay > 0 {
		time.Sleep(cfg.Delay)
	}

	from := mail.NewEmail(cfg.ConfigEmail.EmailFromName, cfg.ConfigEmail.EmailFrom)
	to := mail.NewEmail(cfg.ConfigEmail.EmailToName, cfg.ConfigEmail.EmailTo)
	msg := mail.NewSingleEmail(from,
		cfg.ConfigEmail.Subject,
		to, cfg.ConfigEmail.ContentPlainText,
		cfg.ConfigEmail.ContentHTML,
	)
	_, err := sdk.Sendgrid.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// SendMail sendemail
func (cfg *SDKConfigMailGun) SendMail() error {
	sdk := cfg.newSDKMailGun()

	if CheckIsEmptyCfg(cfg) {
		return fmt.Errorf("Empty fields on ConfigEmail")
	}

	if cfg.Delay > 0 {
		time.Sleep(cfg.Delay)
	}

	msg := sdk.Mailgun.NewMessage(
		cfg.ConfigEmail.EmailFrom,
		cfg.ConfigEmail.Subject,
		cfg.ConfigEmail.ContentPlainText,
		cfg.ConfigEmail.EmailTo,
	)
	_, _, err := sdk.Mailgun.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// CheckIsEmptyCfg check if config is empty
func CheckIsEmptyCfg(cfg interface{}) bool {

	sg, ok1 := cfg.(*SDKConfigSengrid)
	// fmt.Printf("Sendgrid %t\n", ok1)
	if ok1 {
		if len(sg.ConfigEmail.ContentHTML) == 0 ||
			len(sg.ConfigEmail.ContentPlainText) == 0 ||
			len(sg.ConfigEmail.EmailFrom) == 0 ||
			len(sg.ConfigEmail.EmailFromName) == 0 ||
			len(sg.ConfigEmail.EmailTo) == 0 ||
			len(sg.ConfigEmail.EmailToName) == 0 ||
			len(sg.ConfigEmail.Subject) == 0 {
			return true
		}
	}

	mg, ok2 := cfg.(*SDKConfigMailGun)
	// fmt.Printf("MailGun %t\n", ok1)
	if ok2 {
		if len(mg.ConfigEmail.ContentPlainText) == 0 ||
			len(mg.ConfigEmail.EmailFrom) == 0 ||
			len(mg.ConfigEmail.EmailTo) == 0 ||
			len(mg.ConfigEmail.Subject) == 0 {
			return true
		}
	}

	gm, ok3 := cfg.(*SDKConfigGmail)
	// fmt.Printf("MailGun %t\n", ok1)
	if ok3 {
		if len(gm.ConfigEmail.ContentHTML) == 0 ||
			len(gm.ConfigEmail.EmailFrom) == 0 ||
			len(gm.ConfigEmail.EmailTo) == 0 ||
			len(gm.ConfigEmail.Subject) == 0 {
			return true
		}
	}

	return false
}

// SendMail sendemail
func (cfg *SDKConfigGmail) SendMail() error {
	sdk := cfg.newSDKGmail()

	if CheckIsEmptyCfg(cfg) {
		return fmt.Errorf("Empty fields on ConfigEmail")
	}

	if cfg.Delay > 0 {
		time.Sleep(cfg.Delay)
	}

	//Default server
	SMTPServerWithPort := "smtp.gmail.com:587"
	SMTPServerNoPort := "smtp.gmail.com"

	headerConf := make(map[string]string)
	headerMail := make(map[string]string)

	message := ""

	headerConf["Content-Type"] = "text/html; charset=\"UTF-8\";"
	headerConf["MIME-Version"] = "1.0;"

	for key, value := range headerConf {
		message += fmt.Sprintf("%s: %s\n", key, value)
	}

	headerMail["From"] = sdk.Gmail.User
	headerMail["To"] = cfg.ConfigEmail.EmailTo
	headerMail["Subject"] = cfg.ConfigEmail.Subject

	for key, value := range headerMail {
		message += fmt.Sprintf("%s: %s\n", key, value)
	}

	message += "\n" + cfg.ConfigEmail.ContentHTML

	auth := smtp.PlainAuth("", sdk.Gmail.User, sdk.Gmail.Password, SMTPServerNoPort)
	err := smtp.SendMail(SMTPServerWithPort,
		auth,
		sdk.Gmail.User,
		[]string{cfg.ConfigEmail.EmailTo},
		[]byte(message),
	)

	if err != nil {
		return err
	}

	return nil

}
