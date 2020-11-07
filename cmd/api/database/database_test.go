package database

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

func TestOptionsValidate(t *testing.T) {
	opts := &DBOptions{
		Host:   "localhost",
		Port:   5432,
		User:   "test",
		DBName: "test",
		DBPass: "73С7",
		UseSSL: false,
	}

	require.NoError(t, opts.Validate())
}

func TestConnectionString(t *testing.T) {
	testCases := []struct {
		Name           string
		Options        *DBOptions
		ExpectedString string
	}{
		{
			Name: "SSL disabled",
			Options: &DBOptions{
				Host:   "localhost",
				Port:   5432,
				User:   "test",
				DBName: "test",
				DBPass: "73C7",
				UseSSL: false,
			},
			ExpectedString: "host=localhost port=5432 user=test dbname=test password=73C7 sslmode=disable",
		},
		{
			Name: "SSL enabled",
			Options: &DBOptions{
				Host:   "localhost",
				Port:   5432,
				User:   "test",
				DBName: "test",
				DBPass: "73C7",
				UseSSL: true,
			},
			ExpectedString: "host=localhost port=5432 user=test dbname=test password=73C7",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			conn, err := testCase.Options.ConnectionString()

			require.NoError(t, err)
			require.Equal(t, testCase.ExpectedString, conn)
		})
	}
}

func TestOptionsValidateAndConnectionStringMissingFields(t *testing.T) {
	testCases := []struct {
		Name           string
		Options        *DBOptions
		ExpectedErrors []error
	}{
		{
			Name: "Missing host",
			Options: &DBOptions{
				Port:   5432,
				User:   "test",
				DBName: "test",
				DBPass: "pass",
				UseSSL: false,
			},
			ExpectedErrors: []error{ErrDBHostMandatory},
		},
		{
			Name: "Missing port",
			Options: &DBOptions{
				Host:   "localhost",
				User:   "test",
				DBName: "test",
				DBPass: "73С7",
				UseSSL: false,
			},
			ExpectedErrors: []error{ErrDBPortMandatory},
		},
		{
			Name: "Missing user",
			Options: &DBOptions{
				Host:   "localhost",
				Port:   5432,
				DBName: "test",
				DBPass: "73С7",
				UseSSL: false,
			},
			ExpectedErrors: []error{ErrDBUserMandatory},
		},
		{
			Name: "Missing name",
			Options: &DBOptions{
				Host:   "localhost",
				Port:   5432,
				User:   "test",
				DBPass: "73С7",
				UseSSL: false,
			},
			ExpectedErrors: []error{ErrDBNameMandatory},
		},
		{
			Name: "Missing password",
			Options: &DBOptions{
				Host:   "localhost",
				Port:   5432,
				User:   "test",
				DBName: "test",
				UseSSL: false,
			},
			ExpectedErrors: []error{ErrDBPasswordMandatory},
		},
		{
			Name: "Missing password and user",
			Options: &DBOptions{
				Host:   "localhost",
				Port:   5432,
				DBName: "test",
				UseSSL: false,
			},
			ExpectedErrors: []error{ErrDBPasswordMandatory, ErrDBUserMandatory},
		},
		{
			Name: "Missing host and port",
			Options: &DBOptions{
				DBName: "test",
				User:   "user",
				DBPass: "password",
				UseSSL: true,
			},
			ExpectedErrors: []error{ErrDBHostMandatory, ErrDBPortMandatory},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var assertErrors = func(err error) {
				require.Error(t, err)

				var multiErr *multierror.Error
				require.True(t, errors.As(err, &multiErr))
				require.Equal(t, len(testCase.ExpectedErrors), len(multiErr.Errors))
				for _, err := range multiErr.Errors {
					require.Contains(t, testCase.ExpectedErrors, err)
				}
			}

			err := testCase.Options.Validate()
			assertErrors(err)

			conn, err := testCase.Options.ConnectionString()
			require.Empty(t, conn)
			assertErrors(err)
		})
	}
}
