package game

import (
	"../../clashmessager"
	"../clash"
	"../player"
)

type Game interface {
	CommUrl() string
	Playable(user string) bool
	Modes() map[string]Mode
}

type Mode interface {
	Playable(user string) (bool, string)                      //returns whether the mode is playable to the user, and how it appears to this user
	Groups() map[string]player.Group                          //returns the all of the groups that fill up an instance of this game
	Clash([]clash.Player, clashmessager.Messager) clash.Clash //starts a clash with the listed players, will callback the function when done with a list of the tiers telling how the players did
}
