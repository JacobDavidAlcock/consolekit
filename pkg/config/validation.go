package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Validator defines validation rules for configuration values
type Validator struct {
	rules map[string]ValidationRule
}

// ValidationRule represents a validation rule for a config key
type ValidationRule struct {
	Type     string                           // string, int, bool, email, etc.
	Required bool                             // whether the field is required
	Pattern  *regexp.Regexp                   // regex pattern for string validation
	Min      *int                             // minimum value for integers
	Max      *int                             // maximum value for integers
	Custom   func(interface{}) error          // custom validation function
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string]ValidationRule),
	}
}

// AddRule adds a validation rule for a config key
func (v *Validator) AddRule(key string, rule ValidationRule) {
	v.rules[key] = rule
}

// AddStringRule adds a string validation rule
func (v *Validator) AddStringRule(key string, required bool, pattern string) error {
	var regex *regexp.Regexp
	var err error
	
	if pattern != "" {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern for %s: %w", key, err)
		}
	}
	
	v.rules[key] = ValidationRule{
		Type:     "string",
		Required: required,
		Pattern:  regex,
	}
	return nil
}

// AddIntRule adds an integer validation rule
func (v *Validator) AddIntRule(key string, required bool, min, max *int) {
	v.rules[key] = ValidationRule{
		Type:     "int",
		Required: required,
		Min:      min,
		Max:      max,
	}
}

// AddEmailRule adds an email validation rule
func (v *Validator) AddEmailRule(key string, required bool) {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex, _ := regexp.Compile(emailPattern)
	
	v.rules[key] = ValidationRule{
		Type:     "email",
		Required: required,
		Pattern:  regex,
	}
}

// Validate validates a configuration map against defined rules
func (v *Validator) Validate(config map[string]interface{}) error {
	for key, rule := range v.rules {
		value, exists := config[key]
		
		// Check if required field is missing
		if rule.Required && !exists {
			return fmt.Errorf("required field missing: %s", key)
		}
		
		// Skip validation if field is not present and not required
		if !exists {
			continue
		}
		
		// Validate based on type
		if err := v.validateValue(key, value, rule); err != nil {
			return err
		}
	}
	
	return nil
}

// validateValue validates a single value against a rule
func (v *Validator) validateValue(key string, value interface{}, rule ValidationRule) error {
	switch rule.Type {
	case "string":
		return v.validateString(key, value, rule)
	case "int":
		return v.validateInt(key, value, rule)
	case "bool":
		return v.validateBool(key, value, rule)
	case "email":
		return v.validateEmail(key, value, rule)
	default:
		// Use custom validation if available
		if rule.Custom != nil {
			return rule.Custom(value)
		}
	}
	
	return nil
}

// validateString validates a string value
func (v *Validator) validateString(key string, value interface{}, rule ValidationRule) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field %s must be a string", key)
	}
	
	if rule.Pattern != nil && !rule.Pattern.MatchString(str) {
		return fmt.Errorf("field %s does not match required pattern", key)
	}
	
	return nil
}

// validateInt validates an integer value
func (v *Validator) validateInt(key string, value interface{}, rule ValidationRule) error {
	var num int
	var err error
	
	switch v := value.(type) {
	case int:
		num = v
	case string:
		num, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("field %s must be an integer", key)
		}
	default:
		return fmt.Errorf("field %s must be an integer", key)
	}
	
	if rule.Min != nil && num < *rule.Min {
		return fmt.Errorf("field %s must be at least %d", key, *rule.Min)
	}
	
	if rule.Max != nil && num > *rule.Max {
		return fmt.Errorf("field %s must be at most %d", key, *rule.Max)
	}
	
	return nil
}

// validateBool validates a boolean value
func (v *Validator) validateBool(key string, value interface{}, rule ValidationRule) error {
	switch v := value.(type) {
	case bool:
		return nil
	case string:
		if strings.ToLower(v) == "true" || strings.ToLower(v) == "false" {
			return nil
		}
		return fmt.Errorf("field %s must be a boolean (true/false)", key)
	default:
		return fmt.Errorf("field %s must be a boolean", key)
	}
}

// validateEmail validates an email value
func (v *Validator) validateEmail(key string, value interface{}, rule ValidationRule) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field %s must be a string", key)
	}
	
	if rule.Pattern != nil && !rule.Pattern.MatchString(str) {
		return fmt.Errorf("field %s must be a valid email address", key)
	}
	
	return nil
}