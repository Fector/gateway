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
	errChan *chan error
}

func (c *HttpCallback) closeRespBody(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		*c.errChan <- err
		return
	}
}

func (c *HttpCallback) Send(message *model.Message) {
	data, err := json.Marshal(message)
	if err != nil {
		*c.errChan <- err
		return
	}
	req, err := http.NewRequest("POST", message.Callback, bytes.NewBuffer(data))
	if err != nil {
		*c.errChan <- err
		return
	}
	resp, err := c.Do(req)
	defer c.closeRespBody(resp)
	if err != nil {
		*c.errChan <- err
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		*c.errChan <- err
		return
	}
	log.Println(body)
}

func (c *HttpCallback) Error() *chan error {
	return c.errChan
}

func NewHttpCallback() *HttpCallback {
	return &HttpCallback{
		Client: &http.Client{},
	}
}
