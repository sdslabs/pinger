// Package utils contains various utilities that are required across the code.
package utils

import log "github.com/sirupsen/logrus"

func init() {
	var err error
	Config, err = GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
}
