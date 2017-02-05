package fishfinger

import (
	"testing"

	"github.com/facebookgo/ensure"
)

// -----------------------------------------------------------------------------

func TestCompose_NewCompose(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
	}{
		{"regular", "test/regular/docker-compose.yaml", ensure.Nil},
		{"missingFile", "test/badFormat/missing-file.yaml", ensure.NotNil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCompose(tc.path)
			tc.expectedErrorFunc(t, err)
		})
	}
}

func TestCompose_Start(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		serviceName       string
		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
	}{
		{
			"regular",
			"test/regular/docker-compose.yaml",
			"consul-01",
			ensure.Nil,
		},
		{
			"badFormat",
			"test/badFormat/docker-compose.yaml",
			"consul-01",
			ensure.NotNil,
		},
		{
			"missingService",
			"test/regular/docker-compose.yaml",
			"consulito",
			ensure.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCompose(tc.path)
			ensure.Nil(t, err)
			ensure.NotNil(t, c)
			tc.expectedErrorFunc(t, c.Start(tc.serviceName))

			if s, err := c.Status(tc.serviceName); err == nil && s == true {
				ensure.Nil(t, c.Stop(tc.serviceName))
			}
		})
	}
}

func TestCompose_Stop(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		stopAgain bool
	}{
		{"regular", "test/regular/docker-compose.yaml", false},
		{"alreadyStopped", "test/regular/docker-compose.yaml", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCompose(tc.path)
			ensure.Nil(t, err)
			ensure.NotNil(t, c)
			ensure.Nil(t, c.Start("consul-01"))
			ensure.Nil(t, c.Stop("consul-01"))

			if tc.stopAgain == true {
				ensure.Nil(t, c.Stop("consul-01"))
			}
		})
	}
}

func TestCompose_Status(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		serviceName       string
		activateFunc      func(c *Compose)
		expectedCheckFunc func(ensure.Fataler, bool, ...interface{})
		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
	}{
		{
			"regularUp",
			"test/regular/docker-compose.yaml",
			"consul-01",
			func(c *Compose) { _ = c.Start("consul-01") },
			ensure.True,
			ensure.Nil,
		},
		{
			"regularUpAndDown",
			"test/regular/docker-compose.yaml",
			"consul-01",
			func(c *Compose) {
				_ = c.Start("consul-01")
				_ = c.Stop("consul-01")
			},
			ensure.False,
			ensure.Nil,
		},
		{
			"regularNoUp",
			"test/regular/docker-compose.yaml",
			"consul-01",
			func(c *Compose) { _ = c.Stop("consul-01") },
			ensure.False,
			ensure.Nil,
		},
		{
			"missingService",
			"test/regular/docker-compose.yaml",
			"consulito",
			func(c *Compose) {
				_ = c.Start("consulito")
			},
			ensure.False,
			ensure.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCompose(tc.path)
			ensure.Nil(t, err)
			ensure.NotNil(t, c)

			tc.activateFunc(c)
			s, err := c.Status(tc.serviceName)
			tc.expectedCheckFunc(t, s)
			tc.expectedErrorFunc(t, err)
			_ = c.Stop(tc.serviceName)
		})
	}
}

func TestCompose_Port(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		port              string
		expected          string
		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
	}{
		{
			"regular",
			"test/regular/docker-compose.yaml",
			"8500/tcp",
			"0.0.0.0:5015",
			ensure.Nil,
		},
		{
			"missingPort",
			"test/regular/docker-compose.yaml",
			"5555/tcp",
			"",
			ensure.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCompose(tc.path)
			ensure.Nil(t, err)
			ensure.Nil(t, c.Start("consul-01"))

			addr, err := c.Port("consul-01", tc.port)
			tc.expectedErrorFunc(t, err)
			ensure.SameElements(t, addr, tc.expected)
			_ = c.Stop("consul-01")
		})
	}
}

func TestCompose_Env(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		varName           string
		expected          string
		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
	}{
		{
			"regular",
			"test/regular/docker-compose.yaml",
			"DEBUG",
			"working",
			ensure.Nil,
		},
		{
			"missingVar",
			"test/regular/docker-compose.yaml",
			"MISSING_VAR",
			"",
			ensure.NotNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCompose(tc.path)
			ensure.Nil(t, err)
			ensure.Nil(t, c.Start("consul-01"))

			v, err := c.Env("consul-01", tc.varName)
			tc.expectedErrorFunc(t, err)
			ensure.StringContains(t, v, tc.expected)
			_ = c.Stop("consul-01")
		})
	}
}

func TestCompose_Info(t *testing.T) {
	testCases := []struct {
		name              string
		path              string
		startService      bool
		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
		expectedCheckFunc func(ensure.Fataler, bool, ...interface{})
	}{
		{
			"regular",
			"test/regular/docker-compose.yaml",
			true,
			ensure.Nil,
			ensure.True,
		},
		{
			"serviceNotStarted",
			"test/badFormat/docker-compose.yaml",
			false,
			ensure.NotNil,
			ensure.False,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewCompose(tc.path)
			ensure.Nil(t, err)
			ensure.NotNil(t, c)
			if tc.startService {
				ensure.Nil(t, c.Start("consul-01"))
			}
			info, err := c.Info("consul-01")
			tc.expectedErrorFunc(t, err)
			tc.expectedCheckFunc(t, len(info) > 0)
			_ = c.Stop("consul-01")
		})
	}
}

// func TestCompose_Startbackoff(t *testing.T) {
// 	testCases := []struct {
// 		name              string
// 		path              string
// 		startService      bool
// 		backOff           func(*Compose, string) error
// 		expectedErrorFunc func(ensure.Fataler, interface{}, ...interface{})
// 	}{
// 		{
// 			"regular",
// 			"test/regular/docker-compose.yaml",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			c, err := NewCompose(tc.path)
// 			ensure.Nil(t, err)
// 			ensure.NotNil(t, c)

// 		})
// 	}
// }
