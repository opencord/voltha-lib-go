package meters

import (
	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

var logger log.CLogger

func init() {
	// Setup this package so that it's log level can be modified at run time
	var err error
	logger, err = log.RegisterPackage(log.JSON, log.ErrorLevel, log.Fields{})
	if err != nil {
		panic(err)
	}
}
