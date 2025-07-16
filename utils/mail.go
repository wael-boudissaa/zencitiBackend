package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

func SendEmail(email string) {
	// Create a new message
	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", "jw_boudissa@esi.dz") // You can put any sender here (Mailtrap accepts fake ones for testing)
	m.SetHeader("To", email)                  // Receiver's email
	m.SetHeader("Subject", "Hello from Go!")
	m.SetBody("text/plain", "This is a test email sent using Go and Mailtrap!")

	// Create a dialer with your Mailtrap credentials
	d := gomail.NewDialer("smtp.gmail.com", 587, "jw_boudissa@esi.dz", "iwin zgse sjps sand")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Could not send email: %v", err)
	}

	log.Println("Email sent successfully!")
}

func SendRestaurantAdminWelcomeEmail(email, firstName, lastName, password, restaurantName string) error {
	// Use the existing credentials from your SendEmail function
	from := "jw_boudissa@esi.dz"
	smtpPassword := "iwin zgse sjps sand"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Create the email content
	subject := "Welcome to Zenciti - Restaurant Admin Account Created!"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Zenciti</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 3px solid #e74c3c;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 32px;
            font-weight: bold;
            color: #e74c3c;
            margin-bottom: 10px;
        }
        .welcome-text {
            font-size: 24px;
            color: #2c3e50;
            margin-bottom: 20px;
        }
        .credentials-box {
            background: #f8f9fa;
            border: 2px solid #e74c3c;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        .credential-item {
            margin: 10px 0;
            padding: 8px;
            background: white;
            border-radius: 4px;
            border-left: 4px solid #e74c3c;
        }
        .credential-label {
            font-weight: bold;
            color: #2c3e50;
        }
        .credential-value {
            font-family: 'Courier New', monospace;
            color: #e74c3c;
            font-size: 16px;
        }
        .restaurant-info {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            text-align: center;
        }
        .next-steps {
            background: #e8f5e8;
            border-left: 4px solid #27ae60;
            padding: 15px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        .button {
            display: inline-block;
            background: #e74c3c;
            color: white;
            padding: 12px 25px;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 0;
            font-weight: bold;
        }
        .warning {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 5px;
            margin: 15px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🍽️ ZENCITI</div>
            <div style="color: #666;">Restaurant Management Platform</div>
        </div>
        
        <div class="welcome-text">
            Welcome to the Team, %s! 🎉
        </div>
        
        <p>Congratulations! You have been successfully registered as a <strong>Restaurant Administrator</strong> for <strong>%s</strong> on the Zenciti platform.</p>
        
        <div class="restaurant-info">
            <h3>🏪 Your Restaurant</h3>
            <h2>%s</h2>
            <p>You now have full administrative access to manage your restaurant!</p>
        </div>
        
        <div class="credentials-box">
            <h3 style="color: #e74c3c; margin-top: 0;">🔐 Your Login Credentials</h3>
            <div class="credential-item">
                <span class="credential-label">Email:</span><br>
                <span class="credential-value">%s</span>
            </div>
            <div class="credential-item">
                <span class="credential-label">Password:</span><br>
                <span class="credential-value">%s</span>
            </div>
        </div>
        
        <div class="warning">
            <strong>⚠️ Important Security Notice:</strong><br>
            Please change your password after your first login for security purposes.
        </div>
        
        <div class="next-steps">
            <h3 style="color: #27ae60; margin-top: 0;">🚀 What's Next?</h3>
            <ul style="margin: 0; padding-left: 20px;">
                <li>Log in to your admin dashboard</li>
                <li>Complete your restaurant profile</li>
                <li>Set up your menu and food items</li>
                <li>Configure your table layout</li>
                <li>Start managing reservations and orders</li>
            </ul>
        </div>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="#" class="button">🚀 Access Your Dashboard</a>
        </div>
        
        <div style="background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0;">
            <h4 style="color: #2c3e50; margin-top: 0;">📞 Need Help?</h4>
            <p style="margin: 5px 0;">If you have any questions or need assistance getting started, our support team is here to help!</p>
            <p style="margin: 5px 0;">
                📧 Email: <a href="mailto:support@zenciti.com">support@zenciti.com</a><br>
                📱 Phone: +1 (555) 123-4567
            </p>
        </div>
        
        <div class="footer">
            <p><strong>Welcome to Zenciti!</strong></p>
            <p>We're excited to have you on board and look forward to helping you manage your restaurant successfully.</p>
            <p style="font-size: 12px; color: #888;">
                This is an automated message. Please do not reply to this email.<br>
                © 2024 Zenciti. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>`, firstName, restaurantName, restaurantName, email, password)

	// Plain text version for email clients that don't support HTML
	textBody := fmt.Sprintf(`
Welcome to Zenciti, %s!

Congratulations! You have been successfully registered as a Restaurant Administrator for "%s" on the Zenciti platform.

Your Login Credentials:
- Email: %s
- Password: %s

IMPORTANT: Please change your password after your first login for security purposes.

What's Next?
1. Log in to your admin dashboard
2. Complete your restaurant profile
3. Set up your menu and food items
4. Configure your table layout
5. Start managing reservations and orders

Need Help?
Email: support@zenciti.com
Phone: +1 (555) 123-4567

Welcome to Zenciti!
We're excited to have you on board and look forward to helping you manage your restaurant successfully.

© 2024 Zenciti. All rights reserved.
`, firstName, restaurantName, email, password)

	// Set up authentication
	auth := smtp.PlainAuth("", from, smtpPassword, smtpHost)

	// Create message
	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n"+
		"\r\n"+
		"--boundary123\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n"+
		"--boundary123\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n"+
		"--boundary123--\r\n",
		email, subject, textBody, htmlBody)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Welcome email sent successfully to %s for restaurant %s", email, restaurantName)
	return nil
}

func SendActivityAdminWelcomeEmail(email, firstName, lastName, password, restaurantName string) error {
	// Use the existing credentials from your SendEmail function
	from := "jw_boudissa@esi.dz"
	smtpPassword := "iwin zgse sjps sand"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Create the email content
	subject := "Welcome to Zenciti - Restaurant Admin Account Created!"

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Zenciti</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 3px solid #e74c3c;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 32px;
            font-weight: bold;
            color: #e74c3c;
            margin-bottom: 10px;
        }
        .welcome-text {
            font-size: 24px;
            color: #2c3e50;
            margin-bottom: 20px;
        }
        .credentials-box {
            background: #f8f9fa;
            border: 2px solid #e74c3c;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        .credential-item {
            margin: 10px 0;
            padding: 8px;
            background: white;
            border-radius: 4px;
            border-left: 4px solid #e74c3c;
        }
        .credential-label {
            font-weight: bold;
            color: #2c3e50;
        }
        .credential-value {
            font-family: 'Courier New', monospace;
            color: #e74c3c;
            font-size: 16px;
        }
        .restaurant-info {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            text-align: center;
        }
        .next-steps {
            background: #e8f5e8;
            border-left: 4px solid #27ae60;
            padding: 15px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        .button {
            display: inline-block;
            background: #e74c3c;
            color: white;
            padding: 12px 25px;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 0;
            font-weight: bold;
        }
        .warning {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 5px;
            margin: 15px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🍽️ ZENCITI</div>
            <div style="color: #666;">Restaurant Management Platform</div>
        </div>
        
        <div class="welcome-text">
            Welcome to the Team, %s! 🎉
        </div>
        
        <p>Congratulations! You have been successfully registered as a <strong>Restaurant Administrator</strong> for <strong>%s</strong> on the Zenciti platform.</p>
        
        <div class="restaurant-info">
            <h3>🏪 Your Restaurant</h3>
            <h2>%s</h2>
            <p>You now have full administrative access to manage your restaurant!</p>
        </div>
        
        <div class="credentials-box">
            <h3 style="color: #e74c3c; margin-top: 0;">🔐 Your Login Credentials</h3>
            <div class="credential-item">
                <span class="credential-label">Email:</span><br>
                <span class="credential-value">%s</span>
            </div>
            <div class="credential-item">
                <span class="credential-label">Password:</span><br>
                <span class="credential-value">%s</span>
            </div>
        </div>
        
        <div class="warning">
            <strong>⚠️ Important Security Notice:</strong><br>
            Please change your password after your first login for security purposes.
        </div>
        
        <div class="next-steps">
            <h3 style="color: #27ae60; margin-top: 0;">🚀 What's Next?</h3>
            <ul style="margin: 0; padding-left: 20px;">
                <li>Log in to your admin dashboard</li>
                <li>Complete your restaurant profile</li>
                <li>Set up your menu and food items</li>
                <li>Configure your table layout</li>
                <li>Start managing reservations and orders</li>
            </ul>
        </div>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="#" class="button">🚀 Access Your Dashboard</a>
        </div>
        
        <div style="background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0;">
            <h4 style="color: #2c3e50; margin-top: 0;">📞 Need Help?</h4>
            <p style="margin: 5px 0;">If you have any questions or need assistance getting started, our support team is here to help!</p>
            <p style="margin: 5px 0;">
                📧 Email: <a href="mailto:support@zenciti.com">support@zenciti.com</a><br>
                📱 Phone: +1 (555) 123-4567
            </p>
        </div>
        
        <div class="footer">
            <p><strong>Welcome to Zenciti!</strong></p>
            <p>We're excited to have you on board and look forward to helping you manage your restaurant successfully.</p>
            <p style="font-size: 12px; color: #888;">
                This is an automated message. Please do not reply to this email.<br>
                © 2024 Zenciti. All rights reserved.
            </p>
        </div>
    </div>
</body>
</html>`, firstName, restaurantName, restaurantName, email, password)

	// Plain text version for email clients that don't support HTML
	textBody := fmt.Sprintf(`
Welcome to Zenciti, %s!

Congratulations! You have been successfully registered as a Restaurant Administrator for "%s" on the Zenciti platform.

Your Login Credentials:
- Email: %s
- Password: %s

IMPORTANT: Please change your password after your first login for security purposes.

What's Next?
1. Log in to your admin dashboard
2. Complete your restaurant profile
3. Set up your menu and food items
4. Configure your table layout
5. Start managing reservations and orders

Need Help?
Email: support@zenciti.com
Phone: +1 (555) 123-4567

Welcome to Zenciti!
We're excited to have you on board and look forward to helping you manage your restaurant successfully.

© 2024 Zenciti. All rights reserved.
`, firstName, restaurantName, email, password)

	// Set up authentication
	auth := smtp.PlainAuth("", from, smtpPassword, smtpHost)

	// Create message
	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n"+
		"\r\n"+
		"--boundary123\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n"+
		"--boundary123\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n"+
		"--boundary123--\r\n",
		email, subject, textBody, htmlBody)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Welcome email sent successfully to %s for restaurant %s", email, restaurantName)
	return nil
}
