package fishfinger

import (
	"context"
	"fmt"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

// -----------------------------------------------------------------------------

// Compose is a structure that implements `Origin`. It uses `libcompose` to
// manage a Docker environment from a compose-file.
type Compose struct {
	*project.Project
}

// NewCompose returns a new `Compose` or an `error` if something occurs during
// setup. It uses `composeFile` as the path to the compose-file to up.
func NewCompose(composeFile string) (*Compose, error) {
	p, err := docker.NewProject(
		&ctx.Context{
			Context: project.Context{ComposeFiles: []string{composeFile}},
		}, nil)
	if err != nil {
		return nil, err
	}
	return &Compose{Project: p.(*project.Project)}, nil
}

// -----------------------------------------------------------------------------

// Start launches all `services` specified. If `services` is empty, all services
// described in the compose-file are launched. If something occurs during any
// `services` start, it returns an `error`.
func (c *Compose) Start(services ...string) error {
	return c.Up(context.Background(), options.Up{}, services...)
}

// StartBackoff launches all `services` specified and returns only when
// `backoffFunc` completes. If `services` is empty, all services described in
// the compose-file are launched and `backoffFunc` is used for each of them.
//
// NOTE: This project offers default implementations of `backoffFunc` but user
// should provide an implementation according to the Docker images used.
func (c *Compose) StartBackoff(backoffFunc func(*Compose, string) error, services ...string) error {
	for _, service := range services {
		if err := c.Up(context.Background(), options.Up{}, service); err != nil {
			return err
		}
		if err := backoffFunc(c, service); err != nil {
			return err
		}
	}
	return nil
}

// Stop removes all `services` specified. If `services` is empty, all services
// described in the compose-file are removed. If something occurs during any
// `services` stop, it returns an `error`.
func (c *Compose) Stop(services ...string) error {
	if err := c.Down(context.Background(), options.Down{}, services...); err != nil {
		return err
	}
	return c.Delete(context.Background(), options.Delete{}, services...)
}

// Status returns `true` if `service` is in running state, `false` otherwise.
// If `service` does not exist, it returns an `error`.
func (c *Compose) Status(service string) (bool, error) {
	serviceObj, err := c.CreateService(service)
	if err != nil {
		return false, err
	}
	ctnrs, err := serviceObj.Containers(context.Background())
	if err != nil || len(ctnrs) == 0 {
		return false, err
	}
	return ctnrs[0].IsRunning(context.Background())
}

// Port returns the host address where `service` is bound to `port`. If
// `service` does not exist or `port` is not bound, it returns an `error`.
//
// NOTE: `port` MUST be of the following form `<portNb>/<protocol>` such as
// `80/tcp`.
func (c *Compose) Port(service, port string) (string, error) {
	s, err := c.CreateService(service)
	if err != nil {
		return "", err
	}
	ctnr, err := s.Containers(context.Background())
	if err != nil {
		return "", err
	} else if len(ctnr) == 0 {
		return "", fmt.Errorf("no container available")
	}
	addr, err := ctnr[0].Port(context.Background(), port)
	if err != nil {
		return "", err
	} else if addr == "" {
		return "", fmt.Errorf("port not bound")
	}
	return addr, nil
}

// Env returns the value of the environment variable `varName` set in container
// `service` by the compose-file. If `service` or `varName` does not exist, It
// returns an `error`.
func (c *Compose) Env(service, varName string) (string, error) {
	s, err := c.CreateService(service)
	if err != nil {
		return "", err
	}
	v, ok := s.Config().Environment.ToMap()[varName]
	if !ok {
		return "", fmt.Errorf("environment var not found")
	}
	return v, nil
}

// Info returns a `project.InfoSet` about a specified `service`.
func (c *Compose) Info(service string) (project.InfoSet, error) {
	serviceObj, err := c.CreateService(service)
	if err != nil {
		return nil, err
	}
	info, err := serviceObj.Info(context.Background())
	if err != nil {
		return nil, err
	}
	return info, nil
}
