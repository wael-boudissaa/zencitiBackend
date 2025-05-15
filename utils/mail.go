
package utils

import (
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(email string) {
	// Create a new message
	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", "jw_boudissa@esi.dz")          // You can put any sender here (Mailtrap accepts fake ones for testing)
	m.SetHeader("To", email)      // Receiver's email
	m.SetHeader("Subject", "Hello from Go!")
	m.SetBody("text/plain", "This is a test email sent using Go and Mailtrap!")

	// Create a dialer with your Mailtrap credentials
	d := gomail.NewDialer("smtp.gmail.com"	, 587, "jw_boudissa@esi.dz", "iwin zgse sjps sand")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Could not send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
