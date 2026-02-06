package controller

import (
	"errors"
	"strings"
	"unicode"
)

func validatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("パスワードは8文字以上にしてください")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
		hasSymbol bool
	)

	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune("!?-_", char):
			hasSymbol = true
		}
	}

	if !hasUpper {
		return errors.New("大文字を含めてください")
	}
	if !hasLower {
		return errors.New("小文字を含めてください")
	}
	if !hasNumber {
		return errors.New("数字を含めてください")
	}
	if !hasSymbol {
		return errors.New("記号(!?-_)を1文字以上含めてください")
	}

	return nil
}
