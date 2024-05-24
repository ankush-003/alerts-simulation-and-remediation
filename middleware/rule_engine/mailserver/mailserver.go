package mailserver

import (
	"fmt"
	"os"

	"github.com/ankush-003/alerts-simulation-and-remediation/middleware/sim/alerts"

	gomail "gopkg.in/mail.v2"
)

func SendEmail(input alerts.AlertInput, output alerts.AlertOutput, mail string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "alertssim@gmail.com")
	mailer.SetHeader("To", mail)
	mailer.SetHeader("Subject", "Alerts for Node ", input.Origin)
	// fmt.Println("In mailserver", alert)
	// fmt.Println(Severity, Remedy)
	// fmt.Println(ID, Category, CreatedAt, Handled, Source, Origin, Severity, Remedy)
	body := fmt.Sprintf(
		"%s: %s\n, %s: %s\n, %s: %s\n, %s: %s\n, %s: %s\n, %s: %s\n, %s: %s\n",
		"ID", input.ID,
		"Category", input.Category,
		"CreatedAt", input.CreatedAt,
		"Source", input.Source,
		"Origin", input.Origin,
		"Severity", output.Severity,
		"Remedy", output.Remedy,
	)
	// fmt.Println(body)
	mailer.SetBody("text/plain", body)

	// SMTP server settings
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "alertssim@gmail.com", os.Getenv("APP_PWD"))

	err := dialer.DialAndSend(mailer)

	if err != nil {
		fmt.Println("Email sending failed!")
		return fmt.Errorf("error sending email: %v", err)
	}
	fmt.Println("Email sent to ", mail, "successfully!")
	return nil

}
