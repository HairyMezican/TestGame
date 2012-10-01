package wrapper

import (
	"../../definitions/game"
	"../../gamedata"
	"../../model"
	"../../server"
)

type Game struct {
	Base game.Game
}

func (this Game) ModeInfo(pmi *gamedata.PreModeInfo) map[string]gamedata.ModeInfo {
	modes := make(map[string]gamedata.ModeInfo)
	for modename, mode := range this.Base.Modes() {
		playable, name := mode.Playable(pmi.User)
		if playable {
			var modeinfo gamedata.ModeInfo
			modeinfo.Name = name
			groups := make(map[string]gamedata.GroupDescription)
			modeinfo.Groups = &groups
			for groupname, group := range mode.Groups() {
				playable, name := group.Playable(pmi.User)
				if playable {
					var groupinfo gamedata.GroupDescription
					groupinfo.Name = name

					(*modeinfo.Groups)[groupname] = groupinfo
				}
			}

			modes[modename] = modeinfo
		}
	}
	return modes
}

func (this Game) JoinInfo(pji *gamedata.PreJoinInfo) map[string]gamedata.JoinInfo {
	joinMethods := make(map[string]gamedata.JoinInfo)
	mode := this.Base.Modes()[pji.Mode] //TODO: Make sure mode is playable by user
	group := mode.Groups()[pji.Group]   //TODO: Make sure group is playable by user
	joins := group.JoinMethods(pji.User)
	for joinname, join := range joins {
		var joininfo gamedata.JoinInfo
		joininfo.Name = join.Name()
		joinMethods[joinname] = joininfo
	}
	return joinMethods
}

func (this Game) PlayerCounts(ppci *gamedata.PrePlayerCountInfo) gamedata.ModePlayerCounts {
	var playerCounts gamedata.ModePlayerCounts
	mode := this.Base.Modes()[ppci.Mode]
	playerCounts.Players = make(map[string]int)
	for groupname, group := range mode.Groups() {
		playerCounts.Players[groupname] = group.Count()
	}
	return playerCounts
}

func (this Game) Start(info *gamedata.PreStartInfo) gamedata.StartInfo {
	var result gamedata.StartInfo

	clash := model.NewClash(info.Mode, info.ResultUrl)
	for groupname, group := range info.Groups {
		for name, join := range group.Players {
			clash.AddPlayer(name, groupname, join.Join)
		}
	}
	clash.Serialize()

	result.Url = server.GetServerFor(clash)
	result.Players = make(map[string]gamedata.PlayerStartInfo)
	for _, player := range clash.Players {
		result.Players[player.ID] = gamedata.PlayerStartInfo{UrlValues: player.UrlVals()}
	}

	return result
}