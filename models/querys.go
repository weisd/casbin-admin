package models

import (
	"fmt"
	"regexp"
)

// SQLEx SQLEx
type SQLEx struct {
	Order string
	Limit int
}

const (
	ConditionEQ         = "="
	ConditionGT         = ">"
	ConditionGTE        = ">="
	ConditionLT         = "<"
	ConditionLTE        = "<="
	ConditionNOTEQ      = "!="
	ConditionLIKE       = "LIKE"
	ConditionNOTLIKE    = "NOT LIKE"
	ConditionIN         = "IN"
	ConditionNOTIN      = "NOT IN"
	ConditionBETWEEN    = "BETWEEN"
	ConditionNOTBETWEEN = "NOT BETWEEN"
	ConditionISNULL     = "IS NULL"
	ConditionISNOTNULL  = "IS NOT NULL"
)

// Conditions Conditions
var Conditions = []string{
	ConditionEQ,
	ConditionGT,
	ConditionGTE,
	ConditionLT,
	ConditionLTE,
	ConditionNOTEQ,
	ConditionLIKE,
	ConditionNOTLIKE,
	ConditionIN,
	ConditionNOTIN,
	ConditionBETWEEN,
	ConditionNOTBETWEEN,
	ConditionISNULL,
	ConditionISNOTNULL,
}

// MakeCondition MakeCondition
func MakeCondition(cond string) (string, int) {
	switch cond {
	case ConditionIN, ConditionNOTIN:
		return fmt.Sprintf("%s(?)", cond), 1
	case ConditionBETWEEN, ConditionNOTBETWEEN:
		return fmt.Sprintf("%s ? AND ?", cond), 2
	default:
		return fmt.Sprintf(" %s ?", cond), 1
	}
}

// InConditions InConditions
func InConditions(cond string) bool {
	for i := range Conditions {
		if cond == Conditions[i] {
			return true
		}
	}

	return false
}

// ValidField ValidField
func ValidField(field string) bool {
	patten := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return patten.MatchString(field)
}
