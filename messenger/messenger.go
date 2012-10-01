package messenger

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type Codec struct {
	Mime   string
	Encode func(interface{}) (io.Reader, error)
	Decode func(io.Reader, interface{}) error
}

type codecMap map[string]*Codec

type message struct {
	message io.Reader
	mime    string
}

var mimeTypes codecMap = make(codecMap)

func RegisterCodec(c *Codec) {
	mimeTypes[c.Mime] = c
}

type messageHandler func(m message) error

func (this Codec) CreateMessage(content interface{}) (*message, error) {
	r, err := this.Encode(content)
	if err != nil {
		return nil, err
	}
	result := new(message)
	result.message = r
	result.mime = this.Mime
	return result, nil
}

func (this codecMap) DecodeMessage(m message, result interface{}) error {
	codecName := strings.Split(m.mime, ";")[0]
	codec := this[codecName]
	if codec == nil {
		return errors.New("Undefined Codec for: " + m.mime)
	}
	return codec.Decode(m.message, result)
}

func (this message) contactSite(url, method string, mh messageHandler) error {
	req, err := http.NewRequest(method, url, this.message)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", this.mime)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return mh(message{
		mime:    res.Header.Get("content-type"),
		message: res.Body,
	})
}

func (this message) SendTo(url, method string, mh messageHandler) error {
	result := make(chan error, 1)

	t := time.NewTimer(time.Second)
	defer t.Stop()

	go func() {
		result <- this.contactSite(url, method, mh)
	}()

	select {
	case <-t.C:
		return errors.New("Timed Out")
	case err := <-result:
		return err
	}
	panic("unreachable!")
}

func (this message) GetResponse(url, method string, response interface{}) error {
	return this.SendTo(url, method, func(m message) error {
		return mimeTypes.DecodeMessage(m, response)
	})
}
