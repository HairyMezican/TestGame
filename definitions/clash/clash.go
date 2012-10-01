package clash

import (
	"../../hash"
	"../player"
	"code.google.com/p/go.net/websocket"
)

type Clash interface {
	PlayerAction(Player, string) //this gets called whenever a player performs an action
	PlayerConnected(Player)      //lets us know when a player connects
	PlayerDisconnected(Player)   //lets us know when a player has disconnected
}

type Player interface {
	Message(string) //messages the player
	ID() hash.Hash
	Name() string                //gets the players name or ID
	JoinData() player.JoinMethod //gets the join information of the player
	Group() player.Group         //gets the group that the player joined
	IsSpectator() bool
	IsActive() bool
	RemoveConnection(conn *websocket.Conn)
}
