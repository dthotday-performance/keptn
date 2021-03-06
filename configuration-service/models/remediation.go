// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Remediation remediation
//
// swagger:model Remediation
type Remediation struct {

	// Executed action
	Action string `json:"action,omitempty"`

	// ID of the event
	EventID string `json:"eventId,omitempty"`

	// Keptn Context ID of the event
	KeptnContext string `json:"keptnContext,omitempty"`

	// Time of the event
	Time string `json:"time,omitempty"`

	// Type of the event
	Type string `json:"type,omitempty"`
}

// Validate validates this remediation
func (m *Remediation) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Remediation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Remediation) UnmarshalBinary(b []byte) error {
	var res Remediation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
