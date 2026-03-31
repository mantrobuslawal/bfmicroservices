package main

import (
	//"log"
	//"os"
	repo "github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/adapters/repository"
	"github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/adapters/grpc"
	"github.com/mantrobuslawal/bfmicroservices/catalog.git/internal/application/core/api"
	"github.com/mantrobuslawal/bfmicroservices/catalog.git/config"
)

func main() {
	repoAdapter, _ := repo.NewAdapter(repo.SliceCatalog)
	application := api.NewApplication(repoAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}


