package central

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/go-redis/redis/v8"

	"github.com/sdslabs/pinger/pkg/components/agent/proto"
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

	// DEBUG: getAgentWithLowestLoad
	lowestLoadAgent, err := getAgentWithLowestLoad(ctx)
	if err != nil {
		return fmt.Errorf("cannot get agent with lowest load: %w", err)
	}
	fmt.Printf("DEBUG: getAgentWithLowestLoad: %s\n", lowestLoadAgent)

	// DEBUG: AddCheck
	err = AddCheck(ctx)
	if err != nil {
		return err
	}

	return nil
}

func AddCheck(ctx *appcontext.Context) error {
	lowestLoadAgent, err := getAgentWithLowestLoad(ctx)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(lowestLoadAgent, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("cannot connect to agent with tcp address %s: %w", lowestLoadAgent, err)
	}

	defer func() {
		_err := conn.Close()
		if _err != nil {
			ctx.Logger().
				WithError(_err).
				Errorln("could not close connection to agent's grpc server")
		}
	}()

	client := proto.NewAgentClient(conn)

	res, err := client.PushCheck(context.Background(), &proto.Check{
		ID:       "http-get-google",
		Name:     "HTTP Get Google",
		Interval: 10e9,
		Timeout:  5e9,
		Input: &proto.Component{
			Type: "HTTP",
		},
		Output: &proto.Component{
			Type: "TIMEOUT",
		},
		Target: &proto.Component{
			Type:  "URL",
			Value: "http://google.com",
		},
	})
	if err != nil {
		return err
	}

	if res.GetError() != "" {
		return fmt.Errorf("could not add check: %s", res.GetError())
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

	defer func() {
		err := rdb.Close()
		if err != nil {
			ctx.Logger().
				WithError(err).
				Errorln("could not close connection to redis server")
		}
	}()

	zName := "agent_nodes" // replace by reading from config
	return rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   zName,
		Start: "0",
		Stop:  "-1",
	}).Result()
}

func getAgentWithLowestLoad(ctx *appcontext.Context) (string, error) {
	redisServerAddr := "localhost:6379" // replace by reading from central.yml
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisServerAddr,
		Password: "",
		DB:       0,
	})

	defer func() {
		err := rdb.Close()
		if err != nil {
			ctx.Logger().
				WithError(err).
				Errorln("could not close connection to redis server")
		}
	}()

	zName := "agent_nodes" // replace by reading from config
	res, err := rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   zName,
		Start: "0",
		Stop:  "0",
	}).Result()

	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", nil
	}

	return res[0], nil
}
