package kanggo

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Validator 验证器接口
type Validator interface {
	Validate(value interface{}) error
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
	Tag     string
	Value   interface{}
}

// Error 实现 error 接口
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors 多个验证错误
type ValidationErrors []*ValidationError

// Error 实现 error 接口
func (e ValidationErrors) Error() string {
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// StructValidator 结构体验证器
type StructValidator struct {
	tagName string
}

// NewValidator 创建新的验证器
func NewValidator() *StructValidator {
	return &StructValidator{
		tagName: "validate",
	}
}

// Validate 验证结构体
func (v *StructValidator) Validate(s interface{}) error {
	val := reflect.ValueOf(s)

	// 如果是指针，获取其指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 只处理结构体
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validator only accepts struct types")
	}

	var errors ValidationErrors

	// 遍历所有字段
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 获取验证标签
		tag := field.Tag.Get(v.tagName)
		if tag == "" || tag == "-" {
			continue
		}

		// 解析验证规则
		rules := parseRules(tag)

		// 执行验证
		for _, rule := range rules {
			if err := v.validateField(field.Name, fieldValue, rule); err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateField 验证单个字段
func (v *StructValidator) validateField(fieldName string, value reflect.Value, rule Rule) *ValidationError {
	switch rule.Name {
	case "required":
		if !validateRequired(value) {
			return &ValidationError{
				Field:   fieldName,
				Message: "field is required",
				Tag:     "required",
				Value:   value.Interface(),
			}
		}

	case "email":
		if !validateEmail(value) {
			return &ValidationError{
				Field:   fieldName,
				Message: "invalid email format",
				Tag:     "email",
				Value:   value.Interface(),
			}
		}

	case "min":
		if !validateMin(value, rule.Param) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("value must be at least %s", rule.Param),
				Tag:     "min",
				Value:   value.Interface(),
			}
		}

	case "max":
		if !validateMax(value, rule.Param) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("value must be at most %s", rule.Param),
				Tag:     "max",
				Value:   value.Interface(),
			}
		}

	case "len":
		if !validateLen(value, rule.Param) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be %s", rule.Param),
				Tag:     "len",
				Value:   value.Interface(),
			}
		}

	case "minlen":
		if !validateMinLen(value, rule.Param) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be at least %s", rule.Param),
				Tag:     "minlen",
				Value:   value.Interface(),
			}
		}

	case "maxlen":
		if !validateMaxLen(value, rule.Param) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be at most %s", rule.Param),
				Tag:     "maxlen",
				Value:   value.Interface(),
			}
		}

	case "alpha":
		if !validateAlpha(value) {
			return &ValidationError{
				Field:   fieldName,
				Message: "field must contain only letters",
				Tag:     "alpha",
				Value:   value.Interface(),
			}
		}

	case "alphanum":
		if !validateAlphaNum(value) {
			return &ValidationError{
				Field:   fieldName,
				Message: "field must contain only letters and numbers",
				Tag:     "alphanum",
				Value:   value.Interface(),
			}
		}

	case "numeric":
		if !validateNumeric(value) {
			return &ValidationError{
				Field:   fieldName,
				Message: "field must be numeric",
				Tag:     "numeric",
				Value:   value.Interface(),
			}
		}

	case "url":
		if !validateURL(value) {
			return &ValidationError{
				Field:   fieldName,
				Message: "invalid URL format",
				Tag:     "url",
				Value:   value.Interface(),
			}
		}

	case "regex":
		if !validateRegex(value, rule.Param) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("field does not match pattern %s", rule.Param),
				Tag:     "regex",
				Value:   value.Interface(),
			}
		}
	}

	return nil
}

// Rule 验证规则
type Rule struct {
	Name  string
	Param string
}

// parseRules 解析验证规则
func parseRules(tag string) []Rule {
	var rules []Rule

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 解析规则和参数
		if idx := strings.Index(part, "="); idx != -1 {
			rules = append(rules, Rule{
				Name:  strings.TrimSpace(part[:idx]),
				Param: strings.TrimSpace(part[idx+1:]),
			})
		} else {
			rules = append(rules, Rule{
				Name: part,
			})
		}
	}

	return rules
}

// 验证函数实现

func validateRequired(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() != ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return value.Len() > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return value.Float() != 0
	case reflect.Bool:
		return true // bool 类型总是有值
	case reflect.Ptr, reflect.Interface:
		return !value.IsNil()
	}
	return false
}

func validateEmail(value reflect.Value) bool {
	if value.Kind() != reflect.String {
		return false
	}
	email := value.String()
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validateMin(value reflect.Value, param string) bool {
	min, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return false
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()) >= min
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()) >= min
	case reflect.Float32, reflect.Float64:
		return value.Float() >= min
	}
	return false
}

func validateMax(value reflect.Value, param string) bool {
	max, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return false
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()) <= max
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()) <= max
	case reflect.Float32, reflect.Float64:
		return value.Float() <= max
	}
	return false
}

func validateLen(value reflect.Value, param string) bool {
	length, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	switch value.Kind() {
	case reflect.String:
		return len(value.String()) == length
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() == length
	}
	return false
}

func validateMinLen(value reflect.Value, param string) bool {
	minLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	switch value.Kind() {
	case reflect.String:
		return len(value.String()) >= minLen
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() >= minLen
	}
	return false
}

func validateMaxLen(value reflect.Value, param string) bool {
	maxLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	switch value.Kind() {
	case reflect.String:
		return len(value.String()) <= maxLen
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() <= maxLen
	}
	return false
}

func validateAlpha(value reflect.Value) bool {
	if value.Kind() != reflect.String {
		return false
	}
	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return alphaRegex.MatchString(value.String())
}

func validateAlphaNum(value reflect.Value) bool {
	if value.Kind() != reflect.String {
		return false
	}
	alphaNumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphaNumRegex.MatchString(value.String())
}

func validateNumeric(value reflect.Value) bool {
	if value.Kind() != reflect.String {
		return false
	}
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	return numericRegex.MatchString(value.String())
}

func validateURL(value reflect.Value) bool {
	if value.Kind() != reflect.String {
		return false
	}
	url := value.String()
	urlRegex := regexp.MustCompile(`^https?://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=%]+$`)
	return urlRegex.MatchString(url)
}

func validateRegex(value reflect.Value, pattern string) bool {
	if value.Kind() != reflect.String {
		return false
	}
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(value.String())
}

// ValidateStruct 便捷函数：验证结构体
func ValidateStruct(s interface{}) error {
	validator := NewValidator()
	return validator.Validate(s)
}

// Context 扩展方法：验证请求体
func (c *Context) ValidateJSON(v interface{}) error {
	// 先绑定 JSON
	if err := c.BindJSON(v); err != nil {
		return err
	}

	// 然后验证
	return ValidateStruct(v)
}
