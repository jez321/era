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

When an error is wrapped multiple times with era, `Code()` and `Message()` will return the outermost code and message,
so you can overwrite this data further up the call stack where you may have more context.

```go
// In your service layer
return era.Error(
  fmt.Errorf("checking password: %w", err),
  era.WithCode(EInvalidPassword),
)

// In your handler (assuming you return an error to a wrapping error handler
// that then extracts and returns the message)
if err := svc.Login(user, pw); err != nil {
  err = fmt.Errorf("logging in: %w", err)
  opts := era.Options{era.WithFields(era.F{ "user": user })}
  if era.Code(err) == EInvalidPassword {
      opts = append(opts, era.WithMessage("Invalid login credentials."))
  }
  return era.New(err, opts)
}
```

Field data from multiple wrapper era errors is combined, with values from outermost errors taking precedence
if the same key exists more than once.

Note that codes, messages, and field data isn't included in the output of `Error()`.
If you'd like to print them out you can retrieve the values from `Code()`, `Message()`, and `Fields()`.
