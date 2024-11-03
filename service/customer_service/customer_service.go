package customer_service

import "delivery-backend/internal/app"

type Application struct {
	Description string
	Email       string
	PhoneNumber string
	Name        string
}

func init() {
	application_rules := map[string]string{
		"Description": "min=1,max=300",
		"Email":       "required,email,max=50",
		"PhoneNumber": "required,e164",
		"Name":        "min=2,max=20",
	}
	app.RegisterValidation(Application{}, application_rules)
}
