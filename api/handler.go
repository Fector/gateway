package api

import (
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"net/http"
)

type MessageHandler struct {
	ingress *chan model.Message
	memory  *memory.Memory
	error   *chan error
}

func NewMessageHandler(ingress *chan model.Message, memory *memory.Memory, error *chan error) *MessageHandler {
	return &MessageHandler{
		ingress: ingress,
		memory:  memory,
		error:   error,
	}
}

func (h *MessageHandler) get(w http.ResponseWriter, r *http.Request) {

}

func (h *MessageHandler) put(w http.ResponseWriter, r *http.Request) {

}

func (h *MessageHandler) post(w http.ResponseWriter, r *http.Request) {

}

func (h *MessageHandler) patch(w http.ResponseWriter, r *http.Request) {

}

func (h *MessageHandler) delete(w http.ResponseWriter, r *http.Request) {

}

func (h *MessageHandler) options(w http.ResponseWriter, r *http.Request) {

}
