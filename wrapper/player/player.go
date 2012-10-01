package player

import (
	"../../definitions/clash"
	"../../definitions/game"
	"../../definitions/player"
	"../../hash"
	"../../model"
	"code.google.com/p/go.net/websocket"
)

type nothing struct{}

type Player struct {
	hash_id     hash.Hash
	id          string
	group       player.Group
	join        player.JoinMethod
	connections map[*websocket.Conn]nothing
	parent      clash.Clash
}

type Spectator websocket.Conn

func New(mode game.Mode, player *model.Player) *Player {
	p := new(Player)
	groups := mode.Groups()
	p.group = groups[player.Group]
	p.id = player.ID
	p.hash_id = player.Hash
	joins := p.Group().JoinMethods(player.ID)
	p.join = joins[player.Join]
	p.connections = make(map[*websocket.Conn]nothing)
	return p
}

func (this *Player) SetParent(parent clash.Clash) {
	this.parent = parent
}

func (this *Player) AddConnection(conn *websocket.Conn) {
	this.connections[conn] = nothing{}
	if len(this.connections) == 1 {
		this.parent.PlayerConnected(this)
	}
}

func (this *Player) RemoveConnection(conn *websocket.Conn) {
	delete(this.connections, conn)
	if len(this.connections) == 0 {
		this.parent.PlayerDisconnected(this)
	}
}

func (this Player) Message(m string) {
	for conn, _ := range this.connections {
		websocket.Message.Send(conn, m)
	}
}

func (this Player) Name() string {
	return this.id
}

func (this Player) Group() player.Group {
	return this.group
}

func (this Player) JoinData() player.JoinMethod {
	return this.join
}

func (this Player) IsSpectator() bool {
	return false
}

func (this Player) IsActive() bool {
	return len(this.connections) > 0
}

func (this Player) ID() hash.Hash {
	return this.hash_id
}

func (this Spectator) Name() string {
	return ""
}

func (this Spectator) Group() player.Group {
	return nil
}

func (this Spectator) JoinData() player.JoinMethod {
	return nil
}

func (this Spectator) IsSpectator() bool {
	return true
}

func (this Spectator) IsActive() bool {
	return false
}

func (this Spectator) ID() hash.Hash {
	panic("Spectators do not have IDs")
}

func (this *Spectator) Message(s string) {
	websocket.Message.Send((*websocket.Conn)(this), s)
}

func (this *Spectator) RemoveConnection(conn *websocket.Conn) {
}
