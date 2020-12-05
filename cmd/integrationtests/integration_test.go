// +build integration_tests

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/asankov/gira/cmd/api/database"

	"github.com/asankov/gira/pkg/client"
	"github.com/erikh/duct"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/stretchr/testify/require"

	// to register PostgreSQL driver
	_ "github.com/lib/pq"
)

var (
	apiImage        = flag.String("api-image", "ghcr.io/asankov/gira/api", "the image of the api to be used")
	apiVersion      = flag.String("api-version", "latest", "the version of the api image to be used")
	postgresVersion = flag.String("postgres-version", "latest", "the version of the postgres image to be used")

	// testsToRun keeps all the tests that will be executed in the suite.
	// when adding new test, you should add it to this slice
	testsToRun = []func(*testing.T, *client.Client){
		testCreateAndGetAll,
		testUserLifecycle,
	}
)

func TestIntegration(t *testing.T) {
	flag.Parse()

	c := duct.New(duct.Manifest{
		{
			Name:         "gira-api",
			Image:        fmt.Sprintf("%s:%s", *apiImage, *apiVersion),
			PortForwards: map[int]int{21666: 21666},
			Env: []string{
				"PORT=21666",
				"DB_HOST=gira-postgres",
				"DB_PORT=5432",
				"DB_USER=gira",
				"DB_PASS=password",
				"DB_NAME=gira",
				"SESSION_SECRET=sessi0nsecRE7",
				"LOG_LEVEL=debug",
			},
			BootWait: time.Second * 30,
			AliveFunc: func(ctx context.Context, d *docker.Client, s string) error {
				for {
					conn, err := net.Dial("tcp", "localhost:21666")
					if err != nil {
						log.Printf("Error while dialing container: %v", err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					conn.Close()
					return nil
				}
			},
		},
		{
			Name:         "gira-postgres",
			Image:        fmt.Sprintf("postgres:%s", *postgresVersion),
			PortForwards: map[int]int{21665: 5432},
			Env: []string{
				"POSTGRES_USER=gira",
				"POSTGRES_PASSWORD=password",
			},
			BootWait: time.Second * 30,
		},
	}, duct.WithNewNetwork("gira-integration-tests-network"))

	c.HandleSignals(true)

	t.Cleanup(func() {
		if err := c.Teardown(context.Background()); err != nil {
			t.Fatal("error while tearing down services", err)
		}
	})

	if err := c.Launch(context.Background()); err != nil {
		t.Fatal("error while spinning up services:", err)
	}

	if err := database.Init(&database.DBOptions{
		Host:   "localhost",
		Port:   21665,
		User:   "gira",
		DBName: "gira",
		DBPass: "password",
		UseSSL: false,
	}); err != nil {
		t.Fatal("error while initializing db:", err)
	}

	cl, err := client.New("http://localhost:21666")
	require.NoError(t, err)

	for _, test := range testsToRun {
		t.Run(funcName(test), func(t *testing.T) {
			test(t, cl)
		})
	}
}

func funcName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
