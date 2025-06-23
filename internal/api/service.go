package api

import (
	"time"

	"github.com/kardianos/service"
)

var (
	_ service.Interface = (*Shutit)(nil)

	cfg = &service.Config{
		Name:        "shutit",
		DisplayName: "shutit",
		Description: "Service that forces your computer off during specified hours. Go to bed.",
		Arguments:   []string{},
	}
)

type Shutit struct {
	stop chan struct{}

	svc service.Service
}

func NewShutit() (*Shutit, error) {
	s := &Shutit{
		stop: make(chan struct{}),
	}

	svc, err := service.New(s, cfg)
	if err != nil {
		return nil, err
	}

	s.svc = svc

	return s, nil
}

func (s *Shutit) Exec() error {
	return s.svc.Run()
}

func (s *Shutit) Start(svc service.Service) error {
	s.stop = make(chan struct{})

	go s.runner()

	return nil
}

func (s *Shutit) Stop(svc service.Service) error {
	close(s.stop)

	return nil
}

func (s *Shutit) runner() {
	defer func() {
		if service.Interactive() {
			s.Stop(s.svc)
		} else {
			s.svc.Stop()
		}
	}()

	tick := time.NewTicker(30 * time.Second)
	for {
		select {
		case t := <-tick.C:
			if t.Local().Hour() < 22 && t.Local().Hour() > 7 {
				continue
			}

			Shutdown()
		case <-s.stop:
			tick.Stop()

			return
		}
	}
}

func Install() error {
	svc, err := service.New(&Shutit{}, cfg)
	if err != nil {
		return err
	}

	return svc.Install()
}

func Uninstall() error {
	svc, err := service.New(&Shutit{}, cfg)
	if err != nil {
		return err
	}

	return svc.Uninstall()
}
