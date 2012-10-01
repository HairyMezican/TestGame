package game

import (
	"../clashmessager"
	"../definitions/clash"
	"../definitions/game"
	"../definitions/player"
	"../singleton"
	"fmt"
)

func getUserJoinMethods(user string) map[string]player.JoinMethod {
	return map[string]player.JoinMethod{
		"normal": NormalMethod{},
	}
}

type MainGame struct{}

func (MainGame) CommUrl() string {
	return "/setup.json"
}

func (MainGame) Playable(user string) bool {
	return true
}

func (MainGame) Modes() map[string]game.Mode {
	return map[string]game.Mode{
		"normal": NormalMode{},
		"racist": RacistMode{},
	}
}

type NormalMode struct{}

func (NormalMode) Playable(user string) (bool, string) {
	return true, "Play a normal game!"
}

func (NormalMode) Groups() map[string]player.Group {
	return map[string]player.Group{
		"normal": NormalGroup{},
	}
}

type RacistMode struct{}

func (RacistMode) Playable(user string) (bool, string) {
	return true, "Play a racist game!"
}

func (RacistMode) Groups() map[string]player.Group {
	return map[string]player.Group{
		"black": BlackGroup{},
		"white": WhiteGroup{},
	}
}

func (RacistMode) Clash(players []clash.Player, messager clashmessager.Messager) clash.Clash {
	return nil
}

type NormalGroup struct{}

func (NormalGroup) Playable(user string) (bool, string) {
	return true, "It'll be fun!"
}

func (NormalGroup) Count() int {
	return 2
}

func (NormalGroup) JoinMethods(user string) map[string]player.JoinMethod {
	return getUserJoinMethods(user)
}

type BlackGroup struct{}

func (BlackGroup) Playable(user string) (bool, string) {
	return true, "As a black person!"
}

func (BlackGroup) Count() int {
	return 1
}

func (BlackGroup) JoinMethods(user string) map[string]player.JoinMethod {
	return getUserJoinMethods(user)
}

type WhiteGroup struct{}

func (WhiteGroup) Playable(user string) (bool, string) {
	return true, "As a white person!"
}

func (WhiteGroup) Count() int {
	return 1
}

func (WhiteGroup) JoinMethods(user string) map[string]player.JoinMethod {
	return getUserJoinMethods(user)
}

type NormalMethod struct{}

func (NormalMethod) Name() string {
	return "Normal"
}

type PlayerInfo struct {
	active bool
}

func (this *PlayerInfo) deactivate() {
	this.active = false
}

type NormalClash struct {
	playersLeft     int
	activePlayers   map[clash.Player]struct{}
	finishedPlayers map[clash.Player]int
	players         map[clash.Player]*PlayerInfo
	results         [][]clash.Player
	m               clashmessager.Messager
}

func (NormalMode) Clash(players []clash.Player, messager clashmessager.Messager) clash.Clash {
	result := &NormalClash{
		playersLeft:     len(players),
		activePlayers:   map[clash.Player]struct{}{},
		finishedPlayers: map[clash.Player]int{},
		players:         map[clash.Player]*PlayerInfo{},
		results:         [][]clash.Player{},
		m:               messager,
	}
	for _, p := range players {
		result.players[p] = &PlayerInfo{active: true}
		result.activePlayers[p] = struct{}{}
	}
	return result
}

func (this *NormalClash) PlayerDone(p clash.Player) {
	this.m.RemovePlayer(p)

	tier := make([]clash.Player, 1)
	tier[0] = p
	this.results = append(this.results, tier)

	this.players[p].deactivate()
	this.finishedPlayers[p] = len(this.finishedPlayers) + 1
	delete(this.activePlayers, p)

	this.MessagePosition(p)
	if len(this.activePlayers) == 1 {
		for new_p := range this.activePlayers {
			this.PlayerDone(new_p)
		}
		this.m.Results(this.results)
	}
}

func (this *NormalClash) MessagePosition(p clash.Player) {
	place := this.finishedPlayers[p]
	if place == 1 {
		p.Message("Winner")
	} else {
		p.Message("Loser")
	}
	p.Message(fmt.Sprint("Place ", place, []string{"th", "st", "nd", "rd", "th", "th", "th", "th", "th", "th"}[place%10]))
}

func (this *NormalClash) PlayerAction(p clash.Player, action string) {
	switch action {
	case "I'm Here":
		if p.IsSpectator() {
			p.Message("Spectator")
		} else if _, ok := this.activePlayers[p]; ok {
			p.Message("Player")
		} else {
			this.MessagePosition(p)
		}
	case "Win":
		if !p.IsSpectator() && this.players[p].active {
			this.PlayerDone(p)
		}

	}
}

func (this *NormalClash) PlayerDisconnected(p clash.Player) {
}

func (this *NormalClash) PlayerConnected(p clash.Player) {

}

func init() {
	singleton.Set(new(MainGame))
}
