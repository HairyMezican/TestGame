package main

import (
	"fmt"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/TheTemplater/templater"
	"github.com/HairyMezican/Middleware/interceptor"
	"github.com/HairyMezican/Middleware/parser"
	"github.com/HairyMezican/Middleware/websocketer"
	"github.com/HairyMezican/Middleware/renderer"
//	"github.com/HairyMezican/Middleware/sessioner"
	clashwrapper "./wrapper/clash"
	"./wrapper/player"
	wrapper "./wrapper/game"
	gamedef "./definitions/game"
	playerdef "./definitions/clash"
	"./gamedata"
	"./model"
	_ "./game"
	"./singleton"
	"encoding/json"
	"./server"
	"./staticer"
	"math/rand"
)

const (
	gameClash = "gameClash"
	gamePlayer = "gamePlayer"
	clashID = "ClashID"
	playerID = "PlayerID"
	modeIndex = "modeIndex"
	isDispatcher = true
	isGameServer = true
	debugMode = true
)

type V map[string]interface{}

func (this V) Set(c *clashwrapper.Clash,p playerdef.Player) {
	if c != nil {
		this[gameClash] = c
		this[clashID] = c.ID()
	}else{
		delete(this,gameClash)
		delete(this,clashID)
	}
	if p != nil {
		this[gamePlayer] = p
		if !p.IsSpectator() {
			this[playerID] = p.ID()
		}
	}else{
		delete(this,gamePlayer)
		delete(this,playerID)
	}
}

func (this V) Get() (*clashwrapper.Clash,playerdef.Player) {
	c := this[gameClash].(*clashwrapper.Clash)
	p := this[gamePlayer].(*player.Player)
	return c,p
}

type GameQuery struct {
	wrapper.Game
}

type message struct {
	ModeQuery *gamedata.PreModeInfo `json:"modeinfo"`
	JoinQuery *gamedata.PreJoinInfo `json:"joininfo"`
	Start *gamedata.PreStartInfo `json:"start"`
	PlayerCountQuery *gamedata.PrePlayerCountInfo `json:"playercount"`
}

func getErrorString(err interface{}) string {
	if e,ok := err.(error);ok {
		return e.Error()
	}
	if s,ok := err.(string);ok {
		return s
	}
	return "Unknown Error"
}

func (this GameQuery) GetResponse(m message) (response interface{}){
	
	defer func(){
//		rec := recover()
//		if rec != nil {
//			s := getErrorString(rec)
//			response = gamedata.Error{Problem:s}
//		}
	}()
	
	if m.ModeQuery != nil {
		return this.ModeInfo(m.ModeQuery)
	} else if m.JoinQuery != nil {
		return this.JoinInfo(m.JoinQuery)
	} else if m.Start != nil {
		return this.Start(m.Start)
	} else if m.PlayerCountQuery != nil {
		return this.PlayerCounts(m.PlayerCountQuery)
	}
	
	panic("unknown message type")
	return nil
}

type randomer struct {}
func (this randomer) Run(vars map[string]interface{}, next func()) {
	vars["Rand1"] = rand.Int()
	vars["Rand2"] = rand.Int()
	next()
}

func (this GameQuery) Run(vars map[string]interface{}, next func()) {
	v := (httper.V)(vars)
	r := v.GetRequest()

	var m message
	json.NewDecoder(r.Body).Decode(&m)
	
	response,err := json.Marshal(this.GetResponse(m))
	if err != nil {
		return
	}
	v.SetMessage(response)
	v.AddHeader("content-type","application/json")
}

func PlayerGetter(game gamedef.Game) rack.Func {
	return func(vars map[string]interface{},next func()) {
		r := (httper.V)(vars).GetRequest()
		c,p := model.GetClashAndPlayerFromUrl(r.URL)
		if c == nil {
			panic("Can't find clash - "+r.URL.Path)
		}
		if p == nil {
			p = (model.V)(vars).GetPlayerFromCookie(c)
		}
				
		clash := clashwrapper.Get(c)
		player := clash.Get(p)
		V(vars).Set(clash,player)
		vars[modeIndex] = c.Mode
		next()
	}
}

var Renderer rack.Func = func(vars map[string]interface{},next func()) {
	renderer.V(vars).Render(vars[modeIndex].(string))
}

func GameAction(game gamedef.Game) rack.Middleware {	
	ws := websocketer.New()
	ws.OnOpen(rack.Func(func(vars map[string]interface{},next func()){
		c,p := V(vars).Get()
		conn := websocketer.V(vars).GetSocket()
		
		if p != nil {
			p.(*player.Player).AddConnection(conn)
		} else {
			V(vars).Set(c,(*player.Spectator)(conn))
		}
	}))
	ws.OnClose(rack.Func(func(vars map[string]interface{},next func()){
		_,p := V(vars).Get()
		conn := websocketer.V(vars).GetSocket()
		
		p.RemoveConnection(conn)
	}))
	ws.OnMessage(rack.Func(func(vars map[string]interface{},next func()){
		c,p := V(vars).Get()
		m := (websocketer.V)(vars).GetMessage()
		
		c.PlayerAction(p,*m.(*string))
	}))
	return ws
	
}

var Debugger rack.Func = func(vars map[string]interface{},next func()) {
	r := (httper.V)(vars).GetRequest()
	fmt.Println("Request Received for "+r.URL.Path)
	next()
}

func main() {
	templater.LoadFromFiles("views",nil)
	g := singleton.Get()
	
	rackup := rack.New()
	rackup.Add(Debugger)
	if isDispatcher {
		rackup.Add(parser.Form)
		cept := interceptor.New()
		cept.Intercept(g.CommUrl(),GameQuery{wrapper.Game{g}})
		rackup.Add(cept)
	}
	if isGameServer {
		rackup.Add(staticer.New("/", "static"))
//		rackup.Add(sessioner.Middleware)
		rackup.Add(randomer{})
		rackup.Add(PlayerGetter(g))
		rackup.Add(GameAction(g))
		rackup.Add(Renderer)
	}
	
	
	var base string
	var port string
	if debugMode {
		base = "localhost"	
		port = ":3001"
	} else {
		base = "108.166.123.13"
		port = ":80"
	}
	
	if isGameServer {
		server.AddServerCapacity("http://"+base+port,100)	//room for 100 games
	}
	
	conn := httper.HttpConnection(port)
	fmt.Println("Starting Server - use "+base+port+g.CommUrl())
	err := conn.Go(rackup)
	if err != nil {
		fmt.Println("Error:"+err.Error())
	}
}