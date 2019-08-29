package api

import (
	"github.com/DeathHand/gateway/memory"
	"github.com/DeathHand/gateway/model"
	"net/http"
)

type MessageHandler struct {
	ingress *chan model.Message
	memory  *memory.Memory
}

func NewMessageHandler(ingress *chan model.Message, memory *memory.Memory) *MessageHandler {
	return &MessageHandler{
		ingress: ingress,
		memory:  memory,
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
