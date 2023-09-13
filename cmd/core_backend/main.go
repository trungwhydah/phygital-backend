package main

import (
	"log"

	"github.com/go-playground/validator/v10"

	config "backend-service/config/core_backend"

	"backend-service/internal/core_backend/infrastructure/callers"
	"backend-service/internal/core_backend/infrastructure/connections"
	"backend-service/internal/core_backend/infrastructure/firebase"
	"backend-service/internal/core_backend/infrastructure/router"
	"backend-service/internal/core_backend/infrastructure/storage"
	"backend-service/internal/core_backend/registry"
)

func main() {
	config.LoadConfig()

	mongo := connections.NewMongo()
	v := validator.New()
	c := callers.NewCaller()
	fb, err := firebase.NewFirebaseClient(config.C.Firebase.FirebaseProjectID)
	if err != nil {
		log.Fatalln("Failed to Initialize Firebase")
	}
	gs := storage.NewGCPClient()
	rg := registry.NewInteractor(mongo, v, c, fb, gs)
	mdw := rg.NewMiddlewareServices()
	h := rg.NewAppHandler()

	// nftService := rg.NewNFTGlobalService()
	// go nftService.DeployContract()
	// go nftService.ListenEvent()

	router.Initialize(h, mdw)
}
