package handler

import (
	"strings"
)

const START_RUNE = '{'
const END_RUNE = '}'


func formatTemplate(template string, params map[string]string) string {
	if !strings.ContainsRune(template, START_RUNE) {
		return template
	}
	result := ""
	varDetect := false
	variableName := ""
	for _, char := range template {
		if char == START_RUNE {
			varDetect = true
			continue
		}
		if char == END_RUNE {
			varDetect = false
			result += params[variableName]
			variableName = ""
			continue
		}
		if varDetect {
			variableName += string(char)
		} else {
			result += string(char)
		}
	}
	return result
}
