package network

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"vpn-orc/persistence"
)

type ApiGateway struct {
	orchestrator OrchestratorInterface
}

type ApiGatewayInterface interface {
	ListenAndServe()
}

func NewApiGateway(orchestrator OrchestratorInterface) ApiGatewayInterface {
	return &ApiGateway{
		orchestrator: orchestrator,
	}
}

func (a *ApiGateway) registerPeer(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var o persistence.OnboardingRequest
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if o.TenantId == 0 || o.PublicKey == "" || o.Id == "" || o.Addr == "" {
		http.Error(w, "Invalid / Empty fields in request", http.StatusBadRequest)
		return
	}

	// run orchestrator OnboardPeer
	res, err := a.orchestrator.OnboardPeer(o.TenantId, o.Id, []byte(o.PublicKey), o.Addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *ApiGateway) heartbeat(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var h persistence.HeartbeatRequest
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.TenantId == 0 || h.Id == "" {
		http.Error(w, "Invalid / Empty fields in request", http.StatusBadRequest)
		return
	}
	a.orchestrator.UpdatePeer(h.TenantId, h.Id)
}

func status(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "application is running.")
}

func (a *ApiGateway) ListenAndServe() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", status).Methods("GET")
	router.HandleFunc("/connect", a.registerPeer).Methods("POST")
	router.HandleFunc("/heartbeat", a.heartbeat).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
