package mailserver

import (
	"asmr/alerts"
	"fmt"
	"os"

	gomail "gopkg.in/mail.v2"
)



func SendEmail(alert alerts.Alerts, Error error) error {

	
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "alertssim@gmail.com")
	mailer.SetHeader("To", "alertssim@gmail.com")
	mailer.SetHeader("Subject", "Alerts for Node")
	// fmt.Println("In mailserver", alert)
	body := "Alerts for Node:\n"
	// fmt.Println(alertsnew)
	body += fmt.Sprintf("ID: %s \n, NodeID: %s\n, Message: %s\n Severity: %s\n Source:%s\n Time of Creation: %s\n", alert.ID, alert.NodeID, alert.Description, alert.Severity, alert.Source, alert.CreatedAt)
	
	// fmt.Println(body)
	mailer.SetBody("text/plain", body)

	// // SMTP server settings
	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("MAIL_ADDR"), os.Getenv("APP_PWD"))

	err := dialer.DialAndSend(mailer)
	// fmt.Println(err)

	if err != nil {
		fmt.Println("Email sending failed!")
		return fmt.Errorf("error sending email: %v", err)
	}
	fmt.Println("Email sent successfully!")
	return nil

}

