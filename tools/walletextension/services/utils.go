package services

import (
	"fmt"
)

func audit(services *Services, msg string, params ...any) {
	if services.Config.VerboseFlag {
		services.logger.Info(fmt.Sprintf(msg, params...))
	}
}
