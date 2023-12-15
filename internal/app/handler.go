package app

import (
	"context"
	"distapp/models"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"os"
)

func (i *Instance) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, gen := i.storage.GetValue("id")
	fmt.Fprintf(w, "%v:%v", val, gen)
}

func (i *Instance) Set(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if id == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i.storage.SetValue(string(bodyBytes), id)
	//ctx, _ := context.WithTimeout(r.Context(), time.Second*2)
	ctx := context.Background()
	if name, err := os.Hostname(); err == nil {
		ctx = context.WithValue(ctx, "name", name)
	}
	go i.notifyOthers(ctx, i.storage, id)
	fmt.Fprintf(w, "%v", string(bodyBytes))
}

func (i *Instance) Notify(w http.ResponseWriter, r *http.Request) {
	var req models.NotifyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request body given"))
		return
	}

	if changed := i.storage.NotifyValue(req.Value, req.Key, req.Gen); changed {
		log.Printf(
			"NewVal: %v Gen: %v Notifier: %v",
			req.Value,
			req.Gen,
			r.URL.Query().Get("notifier"))
	}
	w.WriteHeader(http.StatusOK)
}
