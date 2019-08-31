package main

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/controller"
	"github.com/sdslabs/status/pkg/probes"
)

func testFunc(ctx context.Context, a string) error {
	fmt.Printf("Inside the test function, %s\n", a)
	prober := probes.NewHTTPProber()

	resp, err := prober.GetProbe("https://google.com/", nil, nil, time.Second*10)
	if err != nil {
		return fmt.Errorf("Error while probing: %s", err)
	}

	fmt.Println("Response From probe is with status code: %d", resp.StatusCode)

	return nil
}

func main() {
	manager := controller.NewManager()

	cf, _ := controller.NewControllerFunction(testFunc, "Sample Argument")
	ctrl := controller.ControllerInternal{
		DoFunc:      cf,
		RunInterval: time.Second * 5,
	}

	manager.UpdateController("test-controller", ctrl)

	time.Sleep(time.Second * 5000)
}

func init() {
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}
