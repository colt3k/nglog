package errexit

import (
	"os"
	"strings"

	log "github.com/colt3k/nglog/ng"
)

// InValidExit
func InValidExit(val interface{}, msg string, exit bool) {

	switch val.(type) {
	case bool:
	case string:
		if len(strings.TrimSpace(val.(string))) <= 0 {
			log.Logln(log.ERROR, msg)
			if exit {
				os.Exit(-1)
			}
		}
	}

}
