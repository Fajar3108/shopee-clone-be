package authaction

import (
	"bytes"
	"path"
	"text/template"

	"github.com/Fajar3108/mafi-course-be/internal/model"
	"github.com/Fajar3108/mafi-course-be/pkg/mail"
	"github.com/gofiber/fiber/v2/log"
)

func SendWelcomeEmail(user *model.User) {
	filepath := path.Join("template", "welcome.html")
	tmpl, err := template.ParseFiles(filepath)

	if err != nil {
		log.Errorf("Failed to send an email: %v", err)
	}

	var body bytes.Buffer

	data := map[string]string{
		"Name": user.Name,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		log.Errorf("Failed to send an email: %v", err)
	}

	if err := mail.SendMail(
		user.Email,
		"Welcome to Mafi Course",
		body.String(),
	); err != nil {
		log.Errorf("Failed to send an email: %v", err)
	}
}
