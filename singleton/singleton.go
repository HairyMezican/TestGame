package singleton

import (
	def "../definitions/game"
)

func Set(g def.Game) {
	singleton = g
}

func Get() def.Game {
	return singleton
}

var singleton def.Game
