package customer_service

import "delivery-backend/internal/app"

type Application struct {
	Status      int8
	Description string
	Email       string
	PhoneNumber string
}

func init() {
	application_rules := map[string]string{
		"Status":      "gte=1,lte=3",
		"Description": "min=1,max=300",
		"License":     "min=1,max=200",
		"Email":       "required,email",
		"PhoneNumber": "required,e164",
	}
	app.RegisterValidation(Application{}, application_rules)
}
