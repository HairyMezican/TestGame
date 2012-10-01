
package list

import (
	"../definitions/player"
	"../hash"
)

type PlayerList struct {
	groups map[player.Group]player.Players
	playersByID map[string]*player.Player
	playersByHash map[hash.Hash]*player.Player
	allPlayers Players
}

func newPlayerList(players Players) PlayerList {
	c.allPlayers = players
	
	c.groups = make(map[string]Players)
	c.playersByID(map[string]*Player)
	c.playersByHash(map[string]*Player)
	
	for player,_ := range(players) {
		group,ok := c.groups[player.Group]
		if !ok {
			group = make(Players)
			c.groups[player.Group] = group
		}
		group.Add(player)

		c.playersByID[player.ID] = p
		c.playersByHash[player.Hash] = p
	}
}

func (this PlayerList) GetPlayerByName(ID string) *Player {
	return this.playersByID[ID]
}

func (this PlayerList) GetPlayerByHash(hash hash.Hash) *Player {
	return this.playersByHash[hash]
}

func (this PlayerList) GetGroup(group string) Players {
	return this.groups[group].Copy()
}

func (this PlayerList) AllPlayers() Players {
	return allPlayers.Copy()
}

type Players map[*Player]struct{}


func (this Players) AsGroup() Players {
	return this
}

func (this Players) Message(m string) {
	for player,_ := range(this) {
		player.Message(m)
	}
}

func (this Players) Add(player *Player) Players {
	this[player] = nothing{}
	return this
}

func (this Players) Plus(player *Player) Players {
	return this.Copy().Add(player)
}

func (this Players) Remove(player *Player) Players {
	delete(this,player)
	return this
}

func (this Players) Minus(player *Player) Players {
	return this.Copy().Remove(player)
}

func (this Players) Contains(player *Player) bool {
	_,ok := this[player]
	return ok
}


type Grouper interface {
	AsGroup() Players
}

func (this Players) AndIs(other Grouper) Players {
	result := make(Players)
	others := Grouper.Group()
	for player,_ := range(this) {
		if others.Contains(player) {
			result.Add(player)
		}
	}
	return result
}

func (this Players) OrIs(other Grouper) Players {
	result := make(Players)
	others := Grouper.Group()
	for player,_ := range(this) {
		result.Add(player)
	}
	for player,_ := range(others) {
		result.Add(player)
	}
	return result
}

func (this Players) Except(other Grouper) Players {
	result := make(Players)
	others := Grouper.Group()
	for player,_ := range(this) {
		if !others.Contains(player) {
			result.Add(player)
		}
	}
	return result
}

func (this Players) Copy() Players {
	result := make(Players)
	for player,_ := range(this) {
		result.Add(player)
	}
	return result
}

func (this Players) List() []clash.Players {
	result := make([]clash.Players,len(this))
	i := 0
	for player,_ := range(this) {
		result[i] = player
		i++
	}
	return result
}