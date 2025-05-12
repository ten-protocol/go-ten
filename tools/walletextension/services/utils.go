package services

import (
	"fmt"
)

func audit(services *Services, msg string, params ...any) {
	//if services.Config.VerboseFlag { // TODO: Ziga - fix this
	services.logger.Info(fmt.Sprintf(msg, params...))
	//}
}
