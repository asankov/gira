// +build integration_tests

package integrationtests

import (
	// to register PostreSQL driver
	"testing"
	"time"

	"github.com/asankov/gira/cmd/api/server"
	"github.com/asankov/gira/internal/auth"
	"github.com/asankov/gira/pkg/client"
	"github.com/asankov/gira/pkg/models/postgres"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) *client.Client {
	db, err := server.NewDB(&server.DBOptions{
		Host:   "localhost",
		Port:   21665,
		User:   "gira",
		DBName: "gira",
		DBPass: "password",
		UseSSL: false,
	})
	require.NoError(t, err)

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	s, err := server.New(&server.Options{
		Log:            log,
		Authenticator:  auth.NewAutheniticator("s3cR37"),
		GameModel:      &postgres.GameModel{DB: db},
		UserModel:      &postgres.UserModel{DB: db},
		UserGamesModel: &postgres.UserGamesModel{DB: db},
		FranchiseModel: &postgres.FranchiseModel{DB: db},
	})
	require.NoError(t, err)

	go func() {
		_ = s.Start(21666)
	}()
	time.Sleep(5 * time.Second)

	cl, err := client.New("http://localhost:21666")
	require.NoError(t, err)

	return cl
}
