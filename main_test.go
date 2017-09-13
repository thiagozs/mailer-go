package mailer_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/thiagozs/mailer-go"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file! Try get a path...")
		log.Printf("PATH = " + getPath() + "/.env")
		if err1 := godotenv.Load(getPath() + "/.env"); err1 != nil {
			log.Printf("Fail...")
			os.Exit(1)
		}
	}
}

func getPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func TestEmptyConfigSendrid(t *testing.T) {
	t.Log("SendMail(Sendgrid) without config... (expected some err)")
	sg := mailer.NewMailerSendGrid("")

	sg.ConfigEmail = mailer.ConfigEmailSendgrid{}

	err := sg.SendMail()
	if err != nil {
		return
	}
	t.Errorf("SendMail(Sendgrid) Fail, got: %s", err)

}

func TestEmptyConfigMailGun(t *testing.T) {
	t.Log("SendMail(MailGun) without config... (expected some err)")
	mg := mailer.NewMailerMailGun("", "", "")

	mg.ConfigEmail = mailer.ConfigEmailMailGun{}

	err := mg.SendMail()
	if err != nil {
		return
	}
	t.Errorf("SendMail(MailGun) Fail, got: %s", err)
}

func TestEmptyConfigGmail(t *testing.T) {
	t.Log("SendMail(Gmail) without config... (expected some err)")
	mg := mailer.NewMailerGmail("", "")

	mg.ConfigEmail = mailer.ConfigEmailGmail{}

	err := mg.SendMail()
	if err != nil {
		return
	}
	t.Errorf("SendMail(Gmail) Fail, got: %s", err)
}

func TestConfigFromSendgrid(t *testing.T) {
	t.Log("SendMail(Sendgrid) with config... (NOT expected some err)")
	mg := mailer.NewMailerSendGrid(os.Getenv("SENDGRID_API_KEY"))

	mg.ConfigEmail = mailer.ConfigEmailSendgrid{
		ContentPlainText: "test",
		ContentHTML:      "<b>test</b>",
		EmailFrom:        "sender@host.com",
		EmailFromName:    "Sender name",
		EmailTo:          "client@host.com",
		EmailToName:      "Client name",
		Subject:          "Test",
	}

	mg.Delay = time.Second * 1

	if mailer.CheckIsEmptyCfg(mg) {
		t.Errorf("SendMail(Sendgrid) Config Empty")
	}
}

func TestConfigFromMailGun(t *testing.T) {
	t.Log("SendMail(MailGun) with config... (NOT expected some err)")
	mg := mailer.NewMailerMailGun("", "", "")

	mg.ConfigEmail = mailer.ConfigEmailMailGun{
		ContentPlainText: "test",
		EmailFrom:        "sender@host.com",
		EmailTo:          "client@host.com",
		Subject:          "Test",
	}

	mg.Delay = time.Second * 1

	if mailer.CheckIsEmptyCfg(mg) {
		t.Errorf("SendMail(MailGun) Config Empty")
	}
}

func TestConfigFromGmail(t *testing.T) {
	t.Log("SendMail(Gmail) with config... (NOT expected some err)")
	mg := mailer.NewMailerGmail("", "")

	mg.ConfigEmail = mailer.ConfigEmailGmail{
		ContentHTML: "test",
		EmailFrom:   "sender@host.com",
		EmailTo:     "client@host.com",
		Subject:     "Test",
	}

	mg.Delay = time.Second * 1

	if mailer.CheckIsEmptyCfg(mg) {
		t.Errorf("SendMail(Gmail) Config Empty")
	}
}

func TestConfigNilFromSendgrid(t *testing.T) {
	t.Log("SendMail(Sendgrid) with NIL config... (expected some err)")
	mg := mailer.NewMailerSendGrid(os.Getenv("SENDGRID_API_KEY"))

	if mailer.CheckIsEmptyCfg(mg) {
		return
	}

	t.Errorf("SendMail(Sendgrid) Config is Nil")
}

func TestConfigNilFromMailgun(t *testing.T) {
	t.Log("SendMail(Mailgun) with NIL config... (expected some err)")
	mg := mailer.NewMailerMailGun(os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
		os.Getenv("MAILGUN_PUB_KEY"))

	if mailer.CheckIsEmptyCfg(mg) {
		return
	}

	t.Errorf("SendMail(Mailgun) Config is Nil")
}

func TestConfigNilFromGmail(t *testing.T) {
	t.Log("SendMail(Gmail) with NIL config... (expected some err)")
	mg := mailer.NewMailerGmail(os.Getenv("GMAIL_USER"),
		os.Getenv("GMAIL_PASSWD"))

	if mailer.CheckIsEmptyCfg(mg) {
		return
	}

	t.Errorf("SendMail(Gmail) Config is Nil")
}

func TestSendMailWithErrorFromSendgrid(t *testing.T) {
	t.Log("SendMail(Sendgrid)... (expected some err)")
	mg := mailer.NewMailerSendGrid(os.Getenv("SENDGRID_API_KEY"))

	mg.ConfigEmail = mailer.ConfigEmailSendgrid{
		ContentPlainText: "test",
		ContentHTML:      "<b>test</b>",
		EmailFrom:        os.Getenv("EMAILFROM"),
		EmailFromName:    os.Getenv("EMAILFROMNAME"),
		EmailTo:          os.Getenv("EMAILTO"),
		EmailToName:      os.Getenv("EMAILTONAME"),
		Subject:          "Test from Sendgrid",
	}

	mg.Delay = time.Second * 1

	err := mg.SendMail()
	if err == nil {
		return
	}

	t.Errorf("SendMail got: %s", err)
}

func TestSendMailWithErrorFromMailGun(t *testing.T) {
	t.Log("SendMail(MailGun)... (expected some err)")
	mg := mailer.NewMailerMailGun(os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
		os.Getenv("MAILGUN_PUB_KEY"))

	mg.ConfigEmail = mailer.ConfigEmailMailGun{
		ContentPlainText: "Test",
		EmailFrom:        os.Getenv("EMAILFROM"),
		EmailTo:          os.Getenv("EMAILTO"),
		Subject:          "Test from Mailgun",
	}

	mg.Delay = time.Second * 1

	err := mg.SendMail()
	if err == nil {
		return
	}

	t.Errorf("SendMail got: %s", err)
}

func TestSendMailWithErrorFromGmail(t *testing.T) {
	t.Log("SendMail(Gmail)... (expected some err)")
	mg := mailer.NewMailerGmail(os.Getenv("GMAIL_LOGIN"),
		os.Getenv("GMAIL_PASSWD"))

	mg.ConfigEmail = mailer.ConfigEmailGmail{
		//ContentHTML: "<html><head><body><b>Test with gmail</b></body></head></html>",
		ContentHTML: "<b>Test with gmail</b>",
		EmailFrom:   os.Getenv("GMAIL_LOGIN"),
		EmailTo:     os.Getenv("EMAILTO"),
		Subject:     "Test from gmail",
	}

	mg.Delay = time.Second * 1

	err := mg.SendMail()
	if err == nil {
		return
	}

	t.Errorf("SendMail got: %s", err)
}
