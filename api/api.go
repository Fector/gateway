package api

import (
	"github.com/DeathHand/gateway/service"
	goji "goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

type Api struct {
	*goji.Mux
}

func NewApi(e *service.Service) *Api {
	mux := goji.NewMux()
	m := NewMessageHandler(e.Ingress, &e.Memory)
	mux.HandleFunc(pat.Get("/message/:id"), m.get)
	mux.HandleFunc(pat.Post("/message"), m.post)
	mux.HandleFunc(pat.Put("/message/:id"), m.put)
	mux.HandleFunc(pat.Delete("/message/:id"), m.delete)
	mux.HandleFunc(pat.Options("/message"), m.options)
	return &Api{Mux: mux}
}

func (a *Api) Serve() {
	log.Fatal(http.ListenAndServe("localhost:8000", a.Mux))
}
