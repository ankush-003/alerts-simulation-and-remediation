package mailserver

import (
	"asmr/rule_engine"
	"fmt"
	"os"
	"time"

	gomail "gopkg.in/mail.v2"
)



func SendEmail(ID string, Category string, CreatedAt time.Time, Handled bool, Source string, Origin string, Params rule_engine.ParamInput, Severity string, Remedy string, err error) error {

	
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "alertssim@gmail.com")
	mailer.SetHeader("To", "alertssim@gmail.com")
	mailer.SetHeader("Subject", "Alerts for Node")
	// fmt.Println("In mailserver", alert)
	body := "Alerts for Node:\n"
	// fmt.Println(Severity, Remedy)
	// fmt.Println(ID, Category, CreatedAt, Handled, Source, Origin, Severity, Remedy)
	body += fmt.Sprintf(
		"%s: %d, %s: %s, %s: %s, %s: %t, %s: %s, %s: %s, %s: %s, %s: %s",
		"ID", ID,
		"Category", Category,
		"CreatedAt", CreatedAt,
		"Handled", Handled,
		"Source", Source,
		"Origin", Origin,
		"Params", Params,
		"Severity", Severity,
		"Remedy", Remedy,
	)	
	// // fmt.Println(body)
	mailer.SetBody("text/plain", body)

	// // SMTP server settings
	dialer := gomail.NewDialer("smtp.gmail.com",587, "alertssim@gmail.com", os.Getenv("APP_PWD"))
	// fmt.Println(dialer)

	err = dialer.DialAndSend(mailer)
	fmt.Println(err)

	if err != nil {
		fmt.Println("Email sending failed!")
		return fmt.Errorf("error sending email: %v", err)
	}
	fmt.Println("Email sent successfully!")
	return nil

}

