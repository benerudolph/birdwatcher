package endpoints

import (
	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Status(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Status()
}