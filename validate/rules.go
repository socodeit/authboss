package validate

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"

	"gopkg.in/authboss.v0"
)

// Rules defines a ruleset by which a string can be validated.
type Rules struct {
	// Field is the name of the field this is intended to validate.
	Field string
	// MatchError describes the MustMatch regexp to a user.
	MatchError           string
	MustMatch            *regexp.Regexp
	MinLength, MaxLength int
	MinLetters           int
	MinNumeric           int
	MinSymbols           int
	AllowWhitespace      bool
}

// Errors returns an array of errors for each validation error that
// is present in the given string. Returns nil if there are no errors.
func (r Rules) Errors(toValidate string) authboss.ErrorList {
	errs := make(authboss.ErrorList, 0)

	ln := len(toValidate)
	if ln == 0 {
		errs = append(errs, authboss.FieldError{r.Field, errors.New("Cannot be blank")})
		return err
	}

	if r.MustMatch != nil {
		if !r.MustMatch.MatchString(toValidate) {
			errs = append(errs, authboss.FieldError{r.Field, errors.New(r.MatchError)})
		}
	}

	if (r.MinLength > 0 && ln < r.MinLength) || (r.MaxLength > 0 && ln > r.MaxLength) {
		errs = append(errs, authboss.FieldError{r.Field, errors.New(r.lengthErr())})
	}

	chars, numeric, symbols, whitespace := tallyCharacters(toValidate)
	if chars < r.MinLetters {
		errs = append(errs, authboss.FieldError{r.Field, errors.New(r.charErr())})
	}
	if numeric < r.MinNumeric {
		errs = append(errs, authboss.FieldError{r.Field, errors.New(r.numericErr())})
	}
	if symbols < r.MinSymbols {
		errs = append(errs, authboss.FieldError{r.Field, errors.New(r.symbolErr())})
	}
	if !r.AllowWhitespace && whitespace > 0 {
		errs = append(errs, authboss.FieldError{r.Field, errors.New("No whitespace permitted")})
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

// IsValid checks toValidate to make sure it's valid according to the rules.
func (r Rules) IsValid(toValidate string) bool {
	return nil == r.Errors(toValidate)
}

// Rules returns an array of strings describing the rules.
func (r Rules) Rules() []string {
	rules := make([]string, 0)

	if r.MustMatch != nil {
		rules = append(rules, r.MatchError)
	}

	if e := r.lengthErr(); len(e) > 0 {
		rules = append(rules, e)
	}
	if e := r.charErr(); len(e) > 0 {
		rules = append(rules, e)
	}
	if e := r.numericErr(); len(e) > 0 {
		rules = append(rules, e)
	}
	if e := r.symbolErr(); len(e) > 0 {
		rules = append(rules, e)
	}

	return rules
}

func (r Rules) lengthErr() (err string) {
	switch {
	case r.MinLength > 0 && r.MaxLength > 0:
		err = fmt.Sprintf("Must be between %d and %d characters", r.MinLength, r.MaxLength)
	case r.MinLength > 0:
		err = fmt.Sprintf("Must be at least %d characters", r.MinLength)
	case r.MaxLength > 0:
		err = fmt.Sprintf("Must be at most %d characters", r.MaxLength)
	}

	return err
}

func (r Rules) charErr() (err string) {
	if r.MinLetters > 0 {
		err = fmt.Sprintf("Must contain at least %d letters", r.MinLetters)
	}
	return err
}

func (r Rules) numericErr() (err string) {
	if r.MinNumeric > 0 {
		err = fmt.Sprintf("Must contain at least %d numbers", r.MinNumeric)
	}
	return err
}

func (r Rules) symbolErr() (err string) {
	if r.MinSymbols > 0 {
		err = fmt.Sprintf("Must contain at least %d symbols", r.MinSymbols)
	}
	return err
}

func tallyCharacters(s string) (chars, numeric, symbols, whitespace int) {
	for _, c := range s {
		switch {
		case unicode.IsLetter(c):
			chars++
		case unicode.IsDigit(c):
			numeric++
		case unicode.IsSpace(c):
			whitespace++
		default:
			symbols++
		}
	}

	return chars, numeric, symbols, whitespace
}
