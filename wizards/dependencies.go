package wizards

import (
	"github.com/malikhisyam/user-graph-service/config"

	relationHttp "github.com/malikhisyam/user-graph-service/domains/relations/handlers/http"
	relationRepo "github.com/malikhisyam/user-graph-service/domains/relations/repositories"
	relationUc "github.com/malikhisyam/user-graph-service/domains/relations/usecases"
	"github.com/malikhisyam/user-graph-service/infrastructures"
)

var (
	Config             = config.GetConfig()
	PostgresDatabase   = infrastructures.NewPostgresDatabase(Config)
	RelationRepository = relationRepo.NewRelationRepository(PostgresDatabase)
	RelationUseCase = relationUc.NewRelationUseCase(RelationRepository)
	RelationHttp = relationHttp.NewRelationHttp(RelationUseCase)
)