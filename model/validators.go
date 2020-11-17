package model

import (
	"errors"
	"reflect"
	"regexp"

	com "github.com/securechat/driver"

	"gopkg.in/validator.v2"
)

func Nonemail(v interface{}, param string) error {
	var myRegex, err = regexp.Compile(".+?\\@.+?\\..{2,5}")
	if err != nil {
		com.Log.Println("Error while compiling regex :" + err.Error())
	}

	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}

	if !myRegex.Match([]byte(st.String())) {
		return errors.New("not in email format")
	}
	return nil
}

func Nonint(v interface{}, param string) error {

	if reflect.ValueOf(v).Kind() != reflect.Int {
		return validator.ErrUnsupported
	}
	return nil
}
