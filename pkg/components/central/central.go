package central

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"github.com/sdslabs/pinger/pkg/config"
	"github.com/sdslabs/pinger/pkg/config/configfile"
	"github.com/sdslabs/pinger/pkg/util/appcontext"

	agentProto "github.com/sdslabs/pinger/pkg/components/agent/proto"
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
	err = AddCheck(ctx, &config.Check{
		ID:       "http-get-google",
		Name:     "HTTP Get Google",
		Interval: time.Duration(10e9),
		Timeout:  time.Duration(5e9),
		Input: config.Component{
			Type: "HTTP",
		},
		Output: config.Component{
			Type: "TIMEOUT",
		},
		Target: config.Component{
			Type:  "URL",
			Value: "http://google.com",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func RunApply(ctx *appcontext.Context, checkdiff *configfile.CheckDiff) error {
	allErrors := ""

	for i := range checkdiff.Additions {
		err := AddCheck(ctx, &checkdiff.Additions[i])
		if err != nil {
			allErrors = fmt.Sprintf("%serror in applying addition: %s\n", allErrors, err.Error())
		}
	}

	for i := range checkdiff.Removals {
		err := RemoveCheck(ctx, checkdiff.Removals[i])
		if err != nil {
			allErrors = fmt.Sprintf("%serror in applying removal: %s\n", allErrors, err.Error())
		}
	}

	if allErrors != "" {
		return errors.New(allErrors)
	}

	return nil
}

func AddCheck(ctx *appcontext.Context, check *config.Check) error {
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

	client := agentProto.NewAgentClient(conn)

	protoCheck := config.CheckToProto(check)
	res, err := client.PushCheck(context.Background(), &protoCheck) // is using context.Background here correct?
	if err != nil {
		return err
	}

	if res.GetError() != "" {
		return fmt.Errorf("could not add check: %s", res.GetError())
	}

	// Also make use of res.GetSuccessful()

	return nil
}

func RemoveCheck(ctx *appcontext.Context, checkID string) error {
	agents, err := getAllAgents(ctx)
	if err != nil {
		return err
	}

	for _, agent := range agents {
		err := removeCheckFromAgent(ctx, checkID, agent)
		if err != nil {
			ctx.Logger().
				WithError(err).
				Errorln("could not broadcast remove check message to an agent")
		}
	}

	return nil
}

func removeCheckFromAgent(ctx *appcontext.Context, checkID string, agent string) error {
	conn, err := grpc.Dial(agent, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("cannot connect to agent with tcp address %s: %w", agent, err)
	}

	defer func() {
		_err := conn.Close()
		if _err != nil {
			ctx.Logger().
				WithError(_err).
				Errorln("could not close connection to agent's grpc server")
		}
	}()

	client := agentProto.NewAgentClient(conn)

	res, err := client.RemoveCheck(context.Background(), &agentProto.CheckID{
		ID: checkID,
	}) // is using context.Background here correct?
	if err != nil {
		return err
	}

	// Also make use of res.GetSuccessful()

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
