package clash

import (
	"../../definitions/clash"
	"../../hash"
	"../../model"
	"../../singleton"
	"../player"
)

type Clash struct {
	id hash.Hash
	clash.Clash
	players map[hash.Hash]*player.Player
}

var allClashes map[hash.Hash]*Clash

func init() {
	allClashes = make(map[hash.Hash]*Clash)
}

func Get(clash *model.Clash) *Clash {
	if c, ok := allClashes[clash.Hash]; ok && c != nil {
		return c
	}

	c := newClash(clash)
	c.id = clash.Hash
	allClashes[clash.Hash] = c
	return c
}

func newClash(c *model.Clash) *Clash {
	game := singleton.Get()
	mode := game.Modes()[c.Mode]

	players := make([]clash.Player, len(c.Players))
	hashmap := make(map[hash.Hash]*player.Player)
	for _, pl := range c.Players {
		p := player.New(mode, pl)

		players[len(hashmap)] = p
		hashmap[pl.Hash] = p
	}

	this := new(Clash)
	this.Clash = mode.Clash(players, c.ResponseUrl)
	this.players = hashmap
	for _, pl := range hashmap {
		pl.SetParent(this.Clash)
	}
	return this
}

func (this Clash) Get(p *model.Player) *player.Player {
	return this.players[p.Hash]
}

func (this Clash) ID() hash.Hash {
	return this.id
}
