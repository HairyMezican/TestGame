package clashmessager

import (
	"../definitions/clash"
	"../messenger"
)

type Messager string

func (this Messager) send(message interface{}, response interface{}) {
	messenger.JSONmessage(message, string(this), "PUT", response)
}

func (this Messager) Results(ranks [][]clash.Player) {
	var results struct {
		Results [][]string `json:"results"`
	}
	results.Results = make([][]string, len(ranks))
	for i, rank := range ranks {
		results.Results[i] = make([]string, len(rank))
		for j, player := range rank {
			results.Results[i][j] = player.Name()
		}
	}
	this.send(results, nil)
}

func (this Messager) RemovePlayer(player clash.Player) {
	var removedPlayer struct {
		Remove string `json:"remove"`
	}
	removedPlayer.Remove = player.Name()
	this.send(removedPlayer, nil)
}
