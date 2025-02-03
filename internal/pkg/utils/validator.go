package utils

import (
	"Backend-trainee-assignment-autumn-2024/internal/model"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var ValidatorInstance = validator.New()

func init() {
	ValidatorInstance.RegisterValidation("bidstatus", func(fl validator.FieldLevel) bool {
		status := fl.Field().String()
		switch status {
		case string(model.BidStatusCreated), string(model.BidStatusPublished), string(model.BidStatusCanceled), string(model.BidStatusApproved), string(model.BidStatusRejected):
			return true
		default:
			return false
		}
	})

	ValidatorInstance.RegisterValidation("bidauthortype", func(fl validator.FieldLevel) bool {
		authorType := fl.Field().String()
		switch authorType {
		case string(model.BidAuthorTypeOrganization), string(model.BidAuthorTypeUser):
			return true
		default:
			return false
		}
	})

	ValidatorInstance.RegisterValidation("servicetype", func(fl validator.FieldLevel) bool {
		serviceType := model.TenderServiceType(fl.Field().String())
		switch serviceType {
		case model.TenderServiceTypeConstruction, model.TenderServiceTypeDelivery, model.TenderServiceTypeManufacture:
			return true
		default:
			return false
		}
	})
}

func ValidateStruct(s interface{}) error {
	err := ValidatorInstance.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				return fmt.Errorf("field '%s' validation failed on the '%s' tag", fieldError.Field(), fieldError.Tag())
			}
		}
		return err
	}
	return nil
}
