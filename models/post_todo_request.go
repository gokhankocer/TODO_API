// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PostTodoRequest post todo request
//
// swagger:model PostTodoRequest
type PostTodoRequest struct {
	

	// description
	Description string `json:"description,omitempty"`

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this post todo request
func (m *PostTodoRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post todo request based on context it is used
func (m *PostTodoRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PostTodoRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostTodoRequest) UnmarshalBinary(b []byte) error {
	var res PostTodoRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
