package mailserver

import (
	"asmr/alerts"
	"fmt"

	gomail "gopkg.in/mail.v2"
)



func SendEmail(alertsnew []alerts.Alerts, Error error) error {

	
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "alertssim@gmail.com")
	mailer.SetHeader("To", "")
	mailer.SetHeader("Subject", "Alerts for Node")
	body := "Alerts for Node:\n"
	for _, alert := range alertsnew {
		body += fmt.Sprintf("ID: %s, NodeID: %s, Message: %s\n", alert.ID, alert.NodeID, alert.Description, alert.CreatedAt)
	}
	mailer.SetBody("text/plain", body)

	// SMTP server settings
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "alertssim@gmail.com", "")

	err := dialer.DialAndSend(mailer)
	fmt.Println(err)

	if err != nil {
		fmt.Println("Email sending failed!")
		return fmt.Errorf("error sending email: %v", err)
	}
	fmt.Println("Email sent successfully!")
	return nil

}

