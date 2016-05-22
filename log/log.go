package log

import (
	"os"

	log "github.com/inconshreveable/log15"
)

func init() {
	root := log.Root()
	root.SetHandler(
		log.CallerFileHandler(
			log.StreamHandler(os.Stdout, log.TerminalFormat())))

}
