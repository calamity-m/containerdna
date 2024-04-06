package paternity

import (
	"github.com/rs/zerolog/log"
)

func Paternity(parent string, child string) {
	log.Debug().Msgf("Working with parent: %s, child: %s", parent, child)
}
