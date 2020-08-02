package main

import (
	"gopkg.in/go-playground/validator.v9"
)

type APIValidator struct {
	validator *validator.Validate
}

func (cv *APIValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (s *Server) Validator() {
	s.server.Validator = &APIValidator{validator: validator.New()}
}
