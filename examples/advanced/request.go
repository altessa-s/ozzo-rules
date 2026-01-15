package main

import (
	ozzo_rules "github.com/altessa-s/ozzo-rules"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// API request validation
type ListRequest struct {
	Limit  int
	Offset int
	Status string
}

func (r ListRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Limit,
			ozzo_rules.ListLimit()),
		validation.Field(&r.Offset,
			ozzo_rules.ListOffset()),
		validation.Field(&r.Status,
			ozzo_rules.OneOf("active", "inactive", "pending")),
	)
}
