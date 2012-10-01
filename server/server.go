package server

import (
	"../global"
	"../model"
	"encoding/json"
	"errors"
	"github.com/HairyMezican/SimpleRedis/redis"
)

func availableServersKey() redis.Set {
	return global.Redis.Set("AvailableServers")
}

func serversInUseKey() redis.Set {
	return global.Redis.Set("InUseServers")
}

func serverCapacity(servername string) redis.Integer {
	return global.Redis.Integer("Servers:" + servername + ":Capacity")
}

func gameServerKey(c *model.Clash) redis.String {
	return c.Prefix().String("Server")
}

func defaultServerKey() redis.String {
	return global.Redis.String("Default-Server")
}

type Server struct {
	Name  string
	ID    int
	clash *model.Clash
}

func (this Server) Data() string {
	result, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(result)
}

func (this *Server) AddClash(c *model.Clash) {
	this.clash = c
	gameServerKey(c).Set(this.Data())
}

func (this Server) Url() string {
	return this.Name + this.clash.Url()
}

func getClashServer(c *model.Clash) *Server {
	s := recreateServer(<-gameServerKey(c).Get())
	s.clash = c
	return s
}

func recreateServer(data string) *Server {
	s := new(Server)
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		panic(err)
	}
	return s
}

func getDefaultServer() *Server {
	server := <-defaultServerKey().Get()
	if server == "" {
		panic(errors.New("No Default Server"))
	}
	return recreateServer(server)
}

func getAvailableServer() *Server {
	server, any := <-availableServersKey().Pop()
	if !any {
		return nil
	}
	serversInUseKey().Add(server)

	return recreateServer(server)
}

func (this Server) Release() {
	if this.ID == -1 {
		return //this is a default server, don't need to store it again
	}
	j := this.Data()

	if <-serversInUseKey().Remove(j) {
		availableServersKey().Add(j)
	}
}

//TODO: Turn the following 4 functions into an interface

func GetServerFor(c *model.Clash) string {
	s := getAvailableServer()
	if s == nil {
		//TODO: log that there are no more available servers
		s = getDefaultServer()
	}
	s.AddClash(c)
	return s.Url()
}

func ReleaseServerFor(c *model.Clash) {
	s := getClashServer(c)
	s.Release()
}

func AddServerCapacity(servername string, newcapacity int) {
	total := <-serverCapacity(servername).IncrementBy(newcapacity)
	for i := total - newcapacity; i < total; i++ {
		newServer := Server{
			Name: servername,
			ID:   i,
		}

		availableServersKey().Add(newServer.Data())
	}
}

func RemoveServerCapacity(servername string, removedcapacity int) {
	total := <-serverCapacity(servername).DecrementBy(removedcapacity)
	for i := total; i < total+removedcapacity; i++ {
		removedServer := Server{
			Name: servername,
			ID:   i,
		}

		data := removedServer.Data()

		availableServersKey().Remove(data)
		serversInUseKey().Remove(data)
	}
}

func SetDefaultServer(servername string) {
	defaultServer := Server{
		Name: servername,
		ID:   -1,
	}
	defaultServerKey().Set(defaultServer.Data())
}

func UnsetDefaultServer() {
	defaultServerKey().Delete()
}
