package era_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jez321/era"
)

func TestEra(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err         error
		wantError   string
		wantCode    string
		wantMessage string
		wantFields  era.F
	}{
		"standard error": {
			err:        era.New(errors.New("an error")),
			wantError:  "an error",
			wantFields: era.F{},
		},
		"with message": {
			err:         era.New(errors.New("an error"), era.WithMessage("a message")),
			wantError:   "an error",
			wantMessage: "a message",
			wantFields:  era.F{},
		},
		"with code": {
			err:        era.New(errors.New("an error"), era.WithCode("a code")),
			wantError:  "an error",
			wantCode:   "a code",
			wantFields: era.F{},
		},
		"with fields": {
			err:        era.New(errors.New("an error"), era.WithFields(era.F{"key": "value", "key2": 25})),
			wantError:  "an error",
			wantFields: era.F{"key": "value", "key2": 25},
		},
		"wrapped": {
			err: era.New(
				fmt.Errorf("wrapped error: %w", era.New(
					errors.New("inner error"),
					era.WithCode("inner code"),
					era.WithMessage("inner message"),
					era.WithFields(era.F{"key": "inner val", "inner": "inner val"}),
				)),
				era.WithCode("outer code"),
				era.WithMessage("outer message"),
				era.WithFields(era.F{"key": "outer val", "outer": "outer val"}),
			),
			wantError:   "wrapped error: inner error",
			wantCode:    "outer code",
			wantMessage: "outer message",
			wantFields:  era.F{"key": "outer val", "inner": "inner val", "outer": "outer val"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.err.Error() != tc.wantError {
				t.Errorf("expected error: %v, got: %v", tc.wantError, tc.err.Error())
			}

			if era.Code(tc.err) != tc.wantCode {
				t.Errorf("expected code: %v, got: %v", tc.wantCode, era.Code(tc.err))
			}

			if era.Message(tc.err) != tc.wantMessage {
				t.Errorf("expected message: %v, got: %v", tc.wantMessage, era.Message(tc.err))
			}

			fldDiff := cmp.Diff(era.Fields(tc.err), tc.wantFields)
			if fldDiff != "" {
				t.Errorf("field data doesn't match: %v", fldDiff)
			}
		})
	}
}

func TestEraMultipleOptions(t *testing.T) {
	opts := era.Options{era.WithCode("abc"), era.WithMessage("def")}
	err := era.New(errors.New("my error"), opts)

	if err.Error() != "my error" {
		t.Errorf("expected error: my error, got: %v", err.Error())
	}

	if era.Code(err) != "abc" {
		t.Errorf("expected code: abc, got: %v", era.Code(err))
	}

	if era.Message(err) != "def" {
		t.Errorf("expected message: def, got: %v", era.Message(err))
	}
}
