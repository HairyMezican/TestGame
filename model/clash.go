//the model is used by the dispatcher to create clashes and save them to a database in a form which can then be later picked up by the game server
package model

import (
	"../clashmessager"
	"../global"
	hashfile "../hash"
	"encoding/json"
	"github.com/HairyMezican/Middleware/sessioner"
	"github.com/HairyMezican/SimpleRedis/redis"
	"net/url"
	"strings"
	"time"
)

type Clash struct {
	Mode        string
	ResponseUrl clashmessager.Messager
	Hash        hashfile.Hash
	Players     map[string]*Player
}

type Player struct {
	Hash  hashfile.Hash
	ID    string
	Group string
	Join  string
}

var Clashes map[string]map[hashfile.Hash]*Clash

func init() {
	Clashes = make(map[string]map[hashfile.Hash]*Clash)
}

func clashKey(mode string, hash hashfile.Hash) redis.String {
	return clashPrefix(mode, hash).String("ClashData")
}

func clashPrefix(mode string, hash hashfile.Hash) redis.Prefix {
	return global.Redis.Prefix("Games:" + mode + ":" + hash.Str() + ":")
}

func NewClash(mode string, url string) *Clash {
	clash := new(Clash)
	clash.Mode = mode
	clash.ResponseUrl = clashmessager.Messager(url)

	h := hashfile.Gethash(mode + time.Now().Format(time.RFC1123))
	clash.Hash = h
	clash.Players = make(map[string]*Player)

	return clash
}

func GetClashAndPlayerFromUrl(u *url.URL) (*Clash, *Player) {
	sPath := u.Path
	for sPath[len(sPath)-1] == "/"[0] {
		sPath = sPath[:len(sPath)-1]
	}
	sections := strings.Split(sPath, "/")
	count := len(sections) - 1
	if count < 1 {
		return nil, nil
	}
	mode := sections[len(sections)-2]
	hash, valid := hashfile.Rehash(sections[len(sections)-1])
	if !valid {
		return nil, nil
	}

	c := loadClashFromHash(mode, hash)

	val := u.Query().Get("player")
	hash, valid = hashfile.Rehash(val)
	if !valid {
		return c, nil
	}
	p := c.LoadPlayerFromHash(hash)
	return c, p
}

type V map[string]interface{}

func (vars V) SetClashAndPlayerToCookie(c *Clash, p *Player) {
	if c != nil {
		(sessioner.V)(vars).Set("Mode", c.Mode)
		(sessioner.V)(vars).Set("Clash", uint32(c.Hash))
		if p != nil {
			(sessioner.V)(vars).Set("Player", uint32(p.Hash))
		} else {
			(sessioner.V)(vars).Clear("Player")
		}
	} else {
		(sessioner.V)(vars).Clear("Clash")
		(sessioner.V)(vars).Clear("Player")
	}
}

func (vars V) GetClashFromCookie() *Clash {
	mode_i := (sessioner.V)(vars).Get("Mode")
	if mode_i == nil {
		return nil
	}
	mode, ok := mode_i.(string)
	if !ok {
		return nil
	}

	clash_i := (sessioner.V)(vars).Get("Clash")
	if clash_i == nil {
		return nil
	}
	clash, ok := clash_i.(uint32)
	if !ok {
		return nil
	}

	return loadClashFromHash(mode, hashfile.Hash(clash))
}

func (vars V) GetPlayerFromCookie(c *Clash) *Player {
	player_i := (sessioner.V)(vars).Get("Player")
	if player_i == nil {
		return nil
	}
	player, ok := player_i.(uint32)
	if !ok {
		return nil
	}

	return c.LoadPlayerFromHash(hashfile.Hash(player))
}

func loadClashFromHash(mode string, hash hashfile.Hash) *Clash {
	if Clashes[mode] == nil {
		Clashes[mode] = make(map[hashfile.Hash]*Clash)
	}
	if Clashes[mode][hash] == nil {
		Clashes[mode][hash] = unserialize(mode, hash)
	}
	return Clashes[mode][hash]
}

func unserialize(mode string, hash hashfile.Hash) *Clash {
	this := new(Clash)
	err := json.Unmarshal([]byte(<-clashKey(mode, hash).Get()), this)
	if err != nil {
		return nil
	}
	return this
}

func (this *Clash) Serialize() {
	clashData, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	key := clashKey(this.Mode, this.Hash)
	key.Set(string(clashData))
}

func (this *Clash) Url() string {
	return "/games/" + this.Mode + "/" + this.Hash.Str()
}

func (this *Clash) AddPlayer(id string, group string, join string) *Player {
	p := new(Player)

	p.ID = id
	p.Group = group
	p.Join = join

	hash := hashfile.Gethash(this.Mode + this.Hash.Str() + id)
	p.Hash = hash

	this.Players[hash.Str()] = p
	return p
}

func (this *Clash) LoadPlayerFromHash(hash hashfile.Hash) *Player {
	if this == nil {
		return nil
	}
	return this.Players[hash.Str()]
}

func (this *Clash) Prefix() redis.Prefix {
	return clashPrefix(this.Mode, this.Hash)
}

func (this *Player) UrlVals() map[string]string {
	return map[string]string{"player": this.Hash.Str()}
}
