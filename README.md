# Era
Simple Go custom errors supporting error codes, friendly messages, and key/value data.

# Usage
```go
// Creating an error
if err := doSomething(x, y); err != nil {
  return era.New(fmt.Errorf("doing something: %w", err),
    era.WithCode(EInternalError),
    era.WithMessage("An internal error occured."),
    era.WithFields(era.F{ "x": x, "y": y }),
  )
}

// Retrieving the custom data
code := era.Code(err) // EInternalError
msg := era.Message(err) // "An internal error occured."
fields := era.Fields(err) // { "x": x, "y": y }
```
