package api

import (
	"cardgame/api/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	battleHandler := handlers.NewBattleHandler()

	// Battle endpoints
	router.HandleFunc("/battle/create", battleHandler.CreateBattle).Methods("POST")
	router.HandleFunc("/battle/{battle_id}/play-card", battleHandler.PlayCard).Methods("POST")
	router.HandleFunc("/battle/{battle_id}", battleHandler.GetBattleState).Methods("GET")

	return router
}
