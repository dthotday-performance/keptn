// Code generated by go-swagger; DO NOT EDIT.

package logs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
	validate "github.com/go-openapi/validate"
)

// GetLogsByEventIDHandlerFunc turns a function with the right signature into a get logs by event Id handler
type GetLogsByEventIDHandlerFunc func(GetLogsByEventIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetLogsByEventIDHandlerFunc) Handle(params GetLogsByEventIDParams) middleware.Responder {
	return fn(params)
}

// GetLogsByEventIDHandler interface for that can handle valid get logs by event Id params
type GetLogsByEventIDHandler interface {
	Handle(GetLogsByEventIDParams) middleware.Responder
}

// NewGetLogsByEventID creates a new http.Handler for the get logs by event Id operation
func NewGetLogsByEventID(ctx *middleware.Context, handler GetLogsByEventIDHandler) *GetLogsByEventID {
	return &GetLogsByEventID{Context: ctx, Handler: handler}
}

/*GetLogsByEventID swagger:route GET /logs/eventId/{eventId} logs getLogsByEventId

Get logs by eventId

*/
type GetLogsByEventID struct {
	Context *middleware.Context
	Handler GetLogsByEventIDHandler
}

func (o *GetLogsByEventID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetLogsByEventIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetLogsByEventIDDefaultBody get logs by event ID default body
// swagger:model GetLogsByEventIDDefaultBody
type GetLogsByEventIDDefaultBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// fields
	Fields string `json:"fields,omitempty"`

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this get logs by event ID default body
func (o *GetLogsByEventIDDefaultBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetLogsByEventIDDefaultBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("getLogsByEventId default"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetLogsByEventIDDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetLogsByEventIDDefaultBody) UnmarshalBinary(b []byte) error {
	var res GetLogsByEventIDDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetLogsByEventIDOKBodyItems0 get logs by event ID o k body items0
// swagger:model GetLogsByEventIDOKBodyItems0
type GetLogsByEventIDOKBodyItems0 struct {

	// event Id
	EventID string `json:"eventId,omitempty"`

	// keptn context
	KeptnContext string `json:"keptnContext,omitempty"`

	// keptn service
	KeptnService string `json:"keptnService,omitempty"`

	// log level
	LogLevel string `json:"logLevel,omitempty"`

	// message
	Message string `json:"message,omitempty"`

	// timestamp
	// Format: date-time
	Timestamp strfmt.DateTime `json:"timestamp,omitempty"`
}

// Validate validates this get logs by event ID o k body items0
func (o *GetLogsByEventIDOKBodyItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetLogsByEventIDOKBodyItems0) validateTimestamp(formats strfmt.Registry) error {

	if swag.IsZero(o.Timestamp) { // not required
		return nil
	}

	if err := validate.FormatOf("timestamp", "body", "date-time", o.Timestamp.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetLogsByEventIDOKBodyItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetLogsByEventIDOKBodyItems0) UnmarshalBinary(b []byte) error {
	var res GetLogsByEventIDOKBodyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
