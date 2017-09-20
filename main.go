package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

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

// SDKConfigAWSSES cfg SDKs
type SDKConfigAWSSES struct {
	SecretKey   string
	AccessKey   string
	Region      string
	SDKName     string
	Delay       time.Duration
	ConfigEmail ConfigEmailAWSSES
}

// SDKConfigSMTPSSL cfg SDKs
type SDKConfigSMTPSSL struct {
	User        string
	Password    string
	Server      string
	Port        string
	SDKName     string
	Delay       time.Duration
	ConfigEmail ConfigEmailSMTPSSL
}

// SDK wrapper of APIs
type SDK struct {
	Mailgun  mailgun.Mailgun
	Sendgrid *sendgrid.Client
	Gmail    GmailLogin
	AWSSES   *ses.SES
	SMTPSSL  SMTPLoginSSL
}

// SMTPLoginSSL username and password for gmail
type SMTPLoginSSL struct {
	User     string
	Password string
	Server   string
	Port     string
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

// ConfigEmailAWSSES configuration of send.
type ConfigEmailAWSSES struct {
	EmailFrom        string
	EmailTo          string
	Subject          string
	ContentPlainText string
	ContentHTML      string
}

// ConfigEmailSMTPSSL configuration of send.
type ConfigEmailSMTPSSL struct {
	EmailFrom        string
	EmailTo          string
	Subject          string
	ContentPlainText string
	ContentHTML      string
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

// newSDKAWSSES get a SDKs
func (cfg SDKConfigAWSSES) newSDKAWSSES() *SDK {
	cred := credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, "")
	return &SDK{
		AWSSES: ses.New(session.New(aws.NewConfig().WithRegion(cfg.Region).WithCredentials(cred))),
	}
}

// newSDKSMTPSSL get a SDKs
func (cfg SDKConfigSMTPSSL) newSDKSMTPSSL() *SDK {
	return &SDK{
		SMTPSSL: SMTPLoginSSL{
			User:     cfg.User,
			Password: cfg.Password,
			Port:     cfg.Port,
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

// NewMailerAWSSES new instance of AWS
func NewMailerAWSSES(accesskey string, secretkey string, region string) *SDKConfigAWSSES {
	return &SDKConfigAWSSES{
		AccessKey: accesskey,
		SecretKey: secretkey,
		Region:    region,
		SDKName:   "awsses",
	}
}

// NewMailerSMTPSSL new instace for SMTP login
func NewMailerSMTPSSL(user string, password string, server string, port string) *SDKConfigSMTPSSL {
	return &SDKConfigSMTPSSL{
		User:     user,
		Password: password,
		Server:   server,
		Port:     port,
		SDKName:  "smtpssl",
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
	// fmt.Printf("Gmail %t\n", ok1)
	if ok3 {
		if len(gm.ConfigEmail.ContentHTML) == 0 ||
			len(gm.ConfigEmail.EmailFrom) == 0 ||
			len(gm.ConfigEmail.EmailTo) == 0 ||
			len(gm.ConfigEmail.Subject) == 0 {
			return true
		}
	}

	smtp, ok4 := cfg.(*SDKConfigSMTPSSL)
	// fmt.Printf("SMTPSSL %t\n", ok1)
	if ok4 {
		if len(smtp.ConfigEmail.ContentHTML) == 0 ||
			len(smtp.ConfigEmail.EmailFrom) == 0 ||
			len(smtp.ConfigEmail.EmailTo) == 0 ||
			len(smtp.ConfigEmail.Subject) == 0 {
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

// SendMail sendemail
func (cfg *SDKConfigAWSSES) SendMail() error {
	if cfg.ConfigEmail.ContentHTML == "" {
		cfg.ConfigEmail.ContentHTML = cfg.ConfigEmail.ContentPlainText
	}

	sdk := cfg.newSDKAWSSES()

	msg := &ses.Message{
		Subject: &ses.Content{
			Charset: aws.String("utf-8"),
			Data:    &cfg.ConfigEmail.Subject,
		},
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &cfg.ConfigEmail.ContentHTML,
			},
			Text: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &cfg.ConfigEmail.ContentPlainText,
			},
		},
	}

	dest := &ses.Destination{
		ToAddresses: aws.StringSlice([]string{cfg.ConfigEmail.EmailTo}),
	}

	_, err := sdk.AWSSES.SendEmail(&ses.SendEmailInput{
		Source:      &cfg.ConfigEmail.EmailFrom,
		Destination: dest,
		Message:     msg,
		// ReplyToAddresses: aws.StringSlice(cfg.ConfigEmail.ReplyTo),
	})

	if err != nil {
		return err
	}

	return nil
}

// SendMail sendemail
func (cfg *SDKConfigSMTPSSL) SendMail() error {
	sdk := cfg.newSDKSMTPSSL()

	if CheckIsEmptyCfg(cfg) {
		return fmt.Errorf("Empty fields on ConfigEmail")
	}

	if cfg.Delay > 0 {
		time.Sleep(cfg.Delay)
	}

	//Default server
	SMTPServerWithPort := fmt.Sprintf("%s:%s", cfg.Server, cfg.Port)
	SMTPServerNoPort := cfg.Server

	headerConf := make(map[string]string)
	headerMail := make(map[string]string)

	message := ""

	headerConf["Content-Type"] = "text/html; charset=\"UTF-8\";"
	headerConf["MIME-Version"] = "1.0;"

	for key, value := range headerConf {
		message += fmt.Sprintf("%s: %s\n", key, value)
	}

	headerMail["From"] = sdk.SMTPSSL.User
	headerMail["To"] = cfg.ConfigEmail.EmailTo
	headerMail["Subject"] = cfg.ConfigEmail.Subject

	for key, value := range headerMail {
		message += fmt.Sprintf("%s: %s\n", key, value)
	}

	message += "\n" + cfg.ConfigEmail.ContentHTML

	auth := smtp.PlainAuth("", sdk.SMTPSSL.User, sdk.SMTPSSL.Password, SMTPServerNoPort)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTPServerNoPort,
	}

	// Conn
	conn, err := tls.Dial("tcp", SMTPServerWithPort, tlsconfig)
	if err != nil {
		return err
	}

	// New Client smtp ssl
	c, err := smtp.NewClient(conn, SMTPServerNoPort)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(cfg.User); err != nil {
		return err
	}

	if err = c.Rcpt(cfg.ConfigEmail.EmailTo); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	// Write messsage
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	// Close Write
	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil

}
