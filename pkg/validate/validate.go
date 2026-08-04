package validate

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mstgnz/cronjob/pkg/config"
	"github.com/robfig/cron/v3"
)

func Validate(structure any) error {
	validate := config.App().Validator
	var errStr string
	var errSlc []error
	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(structure)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}
		for _, err := range err.(validator.ValidationErrors) {
			errStr = fmt.Sprintf("%s %s %s %s", err.Tag(), err.Param(), err.Field(), err.Type().String())
			errSlc = append(errSlc, errors.New(errStr))
			errStr = ""
		}
		// from here you can create your own error messages in whatever language you wish
		return errors.Join(errSlc...)
	}
	return nil
}

// custom validates are called in main
func CustomValidate() {
	CustomCronValidate()
	CustomNoEmptyValidate()
}

// The Go Playground Validator has a “cron” validation mechanism, but it does not work correctly.
// So we will validate with “robfig/cron parser”.
// The expression is only parsed here; it must never be registered on the running
// scheduler, otherwise validation would plant an entry with a nil job that panics on its first tick.
func CustomCronValidate() {
	_ = config.App().Validator.RegisterValidation("cron", func(fl validator.FieldLevel) bool {
		_, err := cron.ParseStandard(fl.Field().String())
		return err == nil
	})
}

// The Go Playground Validator package does not have a validation tag that directly checks whether slices are empty.
// In the case of slices, this tag checks if the slice itself exists, but does not check if the contents of the slice are empty.
// We have written a special validation function to check if slices are empty.
func CustomNoEmptyValidate() {
	_ = config.App().Validator.RegisterValidation("nonempty", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		// Ensure the field is a slice or array
		if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
			return false
		}
		return field.Len() > 0
	})
}
