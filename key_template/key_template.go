package key_template

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
)

const START_RUNE = '{'
const END_RUNE = '}'

func FormatTemplate(template string, params map[string]string) string {
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
			result += extractExpressionFromParams(variableName, params)
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

func jwt(s, name string) string {
	splited := strings.Split(s, ".")
	if len(splited) != 3 {
		return ""
	}
	body := splited[1]
	bodyB, err := b64.StdEncoding.DecodeString(body)
	if err != nil {
		return ""
	}
	claims := map[string]string{}
	json.Unmarshal(bodyB, &claims)
	return claims[name]
}

func hash(s string) string {
	md5Hash := md5.Sum([]byte(s))
	return hex.EncodeToString(md5Hash[:])
}

var functions = map[string]func(s string) string{
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
	"title": strings.ToTitle,
	"trim":  strings.TrimSpace,
	"hash":  hash,
}
var functionsParam = map[string]func(s, param string) string{
	"jwt":     jwt,
	"trim":    func(s, param string) string { return strings.Trim(s, param) },
	"replace": func(s, param string) string { return strings.Replace(s, param, "", -1) },
}

func extractExpressionFromParams(expression string, params map[string]string) string {
	exp := strings.SplitN(expression, "|", 2)
	val := params[exp[0]]
	if len(exp) > 1 {
		fNameS := strings.SplitN(exp[1], ":", 2)
		if len(fNameS) > 1 {
			if f, ok := functionsParam[fNameS[0]]; ok {
				return f(val, fNameS[1])
			}
		}
		if f, ok := functions[exp[1]]; ok {
			return f(val)
		}
	}
	return val
}
