package lintpack

import (
	"log"
	"strings"
)

type parameters map[string]interface{}

func (p parameters) Int(key string, defaultValue int) int {
	k := strings.ToLower(key)
	if value, ok := p[k]; ok {
		if value, ok := value.(int); ok {
			return value
		}
		log.Printf("incorrect value for `%s`, want int", key)
	}
	return defaultValue
}

func (p parameters) String(key, defaultValue string) string {
	k := strings.ToLower(key)
	if value, ok := p[k]; ok {
		if value, ok := value.(string); ok {
			return value
		}
		log.Printf("incorrect value for `%s`, want int", key)
	}
	return defaultValue
}

func (p parameters) Bool(key string, defaultValue bool) bool {
	k := strings.ToLower(key)
	if value, ok := p[k]; ok {
		if value, ok := value.(bool); ok {
			return value
		}
		log.Printf("incorrect value for `%s`, want bool", key)
	}
	return defaultValue
}
