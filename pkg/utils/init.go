// Package utils contains various utilities that are required across the code.
package utils

func init() {
	var err error
	StatusConf, err = GetConfig()
	if err != nil {
		panic(err)
	}
}
