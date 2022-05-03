// Package era provides a simple custom error supporting addition of an error code, friendly error message,
// and additional key/value data.
//
// Usage:
//
//		// Creating an error
//		if err := doSomething(x, y); err != nil {
//			return era.New(fmt.Errorf("doing something: %w", err),
//				era.WithCode(EInternalError),
//				era.WithMessage("An internal error occured."),
//				era.WithFields(era.F{ "x": x, "y": y }),
//			)
//		}
//
//		// Retrieving the custom data
//		code := era.Code(err)
//		msg := era.Message(err)
//		fields := era.Fields(err)
//
// When an error is wrapped multiple times with era, Code() and Message() will return the outermost code and message,
// so you can overwrite this data further up the call stack where you may have more context.
//
//		// In your service layer
//		return era.Error(
//			fmt.Errorf("checking password: %w", err),
//			era.WithCode(EInvalidPassword),
//		)
//
//		// In your handler (assuming you return an error to a wrapping error handler that then
//		// extracts and returns the message)
//		if err := svc.Login(user, pw); err != nil {
//			err = fmt.Errorf("logging in: %w", err)
//			fldOpt := era.WithFields(era.F{ "user": user })
//			if era.Code(err) == EInvalidPassword {
//				return era.New(err, era.WithMessage("Invalid login credentials."), fldOpt)
//			}
//			return era.New(err, era.WithMessage("Internal error."), fldOpt)
//		}
//
// Field data from multiple wrapper era errors is combined, with data from outermost errors taking precedence
// if the same key exists more than once.

package era

import (
	"errors"
)

// Option is an era error option.
type Option interface {
	apply(*eraError)
}

// Option represents multiple era error options using for creating default sets of options.
//
// Usage:
//
//		opts := era.Options{era.WithCode("abc"), era.WithMessage("def")}
//		err := era.New(errors.New("my error"), opts)
type Options []Option

func (o Options) apply(e *eraError) {
	for _, opt := range o {
		opt.apply(e)
	}
}

// F is a map typpe used to store error field data.
type F map[string]interface{}

type eraError struct {
	err     error
	code    string
	message string
	fields  F
}

// Error returns the error string of the wrapped error.
func (e *eraError) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error.
func (e *eraError) Unwrap() error {
	return e.err
}

// New creates a new error with the specified options, wrapping the passed error.
func New(e error, opts ...Option) error {
	err := &eraError{
		err: e,
	}

	for _, opt := range opts {
		opt.apply(err)
	}

	return err
}

type codeOption string

func (o codeOption) apply(e *eraError) {
	e.code = string(o)
}

// WithCode is an option used to specify an error code for the error.
func WithCode(code string) Option {
	return codeOption(code)
}

type messageOption string

func (o messageOption) apply(e *eraError) {
	e.message = string(o)
}

// WithMessage is an option used to specify an friendly error message for the error.
func WithMessage(msg string) Option {
	return messageOption(msg)
}

type fieldsOption F

func (o fieldsOption) apply(e *eraError) {
	e.fields = F(o)
}

// WithFields is an option used to specify key/value data for the error.
func WithFields(fields F) Option {
	return fieldsOption(fields)
}

func (e *eraError) errorCode() string {
	return e.code
}

func (e *eraError) errorMessage() string {
	return e.message
}

func (e *eraError) errorFields() F {
	return e.fields
}

// Code retrieves the error code of the error, or an empty string if no code is present.
// If error codes are defined on multiple wrapped errors, the outermost code will be returned.
func Code(e error) string {
	for e != nil {
		if me, ok := e.(interface{ errorCode() string }); ok && me.errorCode() != "" {
			return me.errorCode()
		}
		e = errors.Unwrap(e)
	}
	return ""
}

// Message retrieves the friendly message of the error, or an empty string if no message is present.
// If messages are defined on multiple wrapped errors, the outermost message will be returned.
func Message(e error) string {
	for e != nil {
		if me, ok := e.(interface{ errorMessage() string }); ok && me.errorMessage() != "" {
			return me.errorMessage()
		}
		e = errors.Unwrap(e)
	}
	return ""
}

// Fields retrieves the field key/value data of the error, or an empty F{} value if no field data is present.
// If field data is defined on multiple wrapped errors, all field data will be returned.
// If the same key exists in multiple wrapped errors, the value of the outermost error will be used.
func Fields(e error) F {
	fields := F{}
	for e != nil {
		if f, ok := e.(interface{ errorFields() F }); ok {
			addFields := f.errorFields()
			for k, v := range addFields {
				// If the same key already exists, don't replace it
				if _, ok := fields[k]; ok {
					continue
				}
				fields[k] = v
			}
		}
		e = errors.Unwrap(e)
	}
	return fields
}
