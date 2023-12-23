package helpers

import (
	"io"

	"github.com/antsrp/gdb_ex/internal/interfaces/logger"
)

/*func HandleCloser(logger *zap.Logger, resource string, closer io.Closer) {
	if err := closer.Close(); err != nil {
		logger.Sugar().Infof("Can't close %q: %s", resource, err)
	}
}*/

func HandleCloser(log logger.Logger, resource string, closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Info("Can't close %q: %s", resource, err)
	}
}
