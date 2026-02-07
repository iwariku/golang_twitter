package controller

import (
	"errors"
	"regexp"
)

func validatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("パスワードは8文字以上にしてください")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(p) {
		return errors.New("大文字を含めてください")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(p) {
		return errors.New("数字を含めてください")
	}
	if !regexp.MustCompile(`[!?_\-]`).MatchString(p) {
		return errors.New("記号(!?-_)を1文字以上含めてください")
	}

	return nil
}
