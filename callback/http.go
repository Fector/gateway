package callback

import (
	"bytes"
	"encoding/json"
	"github.com/DeathHand/gateway/model"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpCallback struct {
	Callback
	*http.Client
	notify *chan model.Message
	error  *chan error
}

func (c *HttpCallback) closeRespBody(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		*c.error <- err
		return
	}
}

func (c *HttpCallback) send(message *model.Message) {
	data, err := json.Marshal(message)
	if err != nil {
		*c.error <- err
		return
	}
	req, err := http.NewRequest("POST", message.Callback, bytes.NewBuffer(data))
	if err != nil {
		*c.error <- err
		return
	}
	resp, err := c.Do(req)
	defer c.closeRespBody(resp)
	if err != nil {
		*c.error <- err
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		*c.error <- err
		return
	}
	log.Println(body)
}

func (c *HttpCallback) Add(message *model.Message) {
	*c.notify <- *message
}

func (c *HttpCallback) Run() {
	for {
		message := <-*c.notify
		c.send(&message)
	}
}

func NewHttpCallback(error *chan error) *HttpCallback {
	return &HttpCallback{
		Client: &http.Client{},
		error:  error,
	}
}
