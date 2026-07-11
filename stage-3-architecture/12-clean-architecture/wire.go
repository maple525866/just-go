//go:build wireinject

package main

import (
	"github.com/google/wire"

	"just-go/stage-3-architecture/12-clean-architecture/infrastructure/memory"
	"just-go/stage-3-architecture/12-clean-architecture/interface/httpapi"
	"just-go/stage-3-architecture/12-clean-architecture/usecase"
)

func initializeHandler() *httpapi.Handler {
	wire.Build(
		memory.NewArticleRepository,
		wire.Bind(new(usecase.ArticleRepository), new(*memory.ArticleRepository)),
		usecase.NewSystemClock,
		usecase.NewSequentialIDGenerator,
		usecase.NewArticleService,
		httpapi.NewHandler,
	)
	return nil
}
