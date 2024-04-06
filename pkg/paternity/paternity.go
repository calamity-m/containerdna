package paternity

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"runtime/debug"
)

func Paternity(parent string, child string) {
	log.Debug().Msgf("Working with parent: %s, child: %s", parent, child)

	info, _ := debug.ReadBuildInfo()

	fmt.Println(info)

}
