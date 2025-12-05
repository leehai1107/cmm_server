package dbfx

import (
	"github.com/leehai1107/cmm_server/pkg/infra"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// const defaultTimeout = 10 * time.Second

var Module = fx.Provide(
	providePostgresDB,

// provideMongoDBClient,
// provideMongoDBDatabase,
// provideMongoDBService,
)

func providePostgresDB() *gorm.DB {
	return infra.GetDB()
}

// func provideMongoDBClient(lifecycle fx.Lifecycle) (*mongo.Client, error) {
// 	dbCfg := config.DBConfig()
// 	client, err := infra.InitMongoDBClient(dbCfg.MongoURI)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lifecycle.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
// 			defer cancel()
// 			client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbCfg.MongoURI))
// 			if err != nil {
// 				logger.Errorw("failed to connect mongo DB", "error", err)
// 				return err
// 			}
// 			if err := client.Ping(ctx, nil); err != nil {
// 				logger.Errorw("failed to ping mongo DB", "error", err)
// 				return err
// 			}
// 			logger.Info("connect mongo DB successful!!!")
// 			return nil
// 		},
// 		OnStop: func(ctx context.Context) error {
// 			return client.Disconnect(ctx)
// 		},
// 	})
// 	return client, nil
// }

// func provideMongoDBDatabase(client *mongo.Client) (*mongo.Database, error) {
// 	dbCfg := config.DBConfig()
// 	if dbCfg.MongoURI == "" {
// 		return nil, errors.New("mongo URI is empty")
// 	}
// 	connStr, err := connstring.Parse(dbCfg.MongoURI)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if connStr.Database == "" {
// 		return nil, errors.New("database name is empty, just configure it in mongo uri")
// 	}
// 	return client.Database(connStr.Database), nil
// }

// func provideMongoDBService(client *mongo.Client, db *mongo.Database) mgo.MongoDBService {
// 	return mgo.NewMongoDBService(client, db)
// }
