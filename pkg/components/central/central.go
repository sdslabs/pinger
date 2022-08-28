package central

import (
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

func Run(ctx *appcontext.Context) error {
	fmt.Println("Central server is running!")
	// run GRPC server and expose API for managing checks

	// DEBUG: getAllAgents
	agents, err := getAllAgents(ctx)
	if err != nil {
		return fmt.Errorf("cannot list agents: %w", err)
	}

	for _, agent := range agents {
		fmt.Printf("DEBUG: getAllAgents: %s\n", agent)
	}

	return nil
}

func getAllAgents(ctx *appcontext.Context) ([]string, error) {
	redisServerAddr := "localhost:6379" // replace by reading from central.yml
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisServerAddr,
		Password: "",
		DB:       0,
	})

	zName := "agent_nodes" // replace by reading from config
	return rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   zName,
		Start: "0",
		Stop:  "-1",
	}).Result()
}
