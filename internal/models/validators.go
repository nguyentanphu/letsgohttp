package models

import "github.com/go-playground/validator/v10"

func (userModel *UserModel) UniqueEmailValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	exists, _ := userModel.EmailExists(value)
	return !exists
}
