package wizards

import (
	"github.com/malikhisyam/user-graph-service/config"
	"github.com/malikhisyam/user-graph-service/shared/util"

	relationHttp "github.com/malikhisyam/user-graph-service/domains/relations/handlers/http"
	relationRepo "github.com/malikhisyam/user-graph-service/domains/relations/repositories"
	relationUc "github.com/malikhisyam/user-graph-service/domains/relations/usecases"
	"github.com/malikhisyam/user-graph-service/infrastructures"
)

var (
	Config             = config.GetConfig()
	PostgresDatabase   = infrastructures.NewPostgresDatabase(Config)
	RedisClient        = infrastructures.InitRedis()
	LoggerInstance, _ = util.NewLogger();
	RelationRepository = relationRepo.NewRelationRepository(PostgresDatabase, RedisClient, LoggerInstance)
	RelationUseCase = relationUc.NewRelationUseCase(RelationRepository)
	RelationHttp = relationHttp.NewRelationHttp(RelationUseCase)
)