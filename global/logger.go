package global

import (
	"github.com/zm50/gte/core"
	"github.com/zm50/gte/trait"
)

var logger trait.Logger

func SetLogger(filename string, maxSize int, maxBackups int, maxAge int, compress bool) {
	logger = core.NewLogger(filename, maxSize, maxBackups, maxAge, compress)
}

func Logger() trait.Logger {
	return logger
}
