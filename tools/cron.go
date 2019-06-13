package tools

import "github.com/robfig/cron"

type Crontab struct {
	Spec 		string

	execFunc 	func()
	tic 		*cron.Cron
}

func NewCrontab(spec string, fn func()) *Crontab {
	return &Crontab{
		Spec: spec,
		execFunc: fn,
	}
}

func (c *Crontab) Start() error {
	tic := cron.New()
	if err := tic.AddFunc(c.Spec, c.execFunc); err != nil {
		return err
	}
	c.tic = tic
	tic.Start()

	return nil
}

func (c *Crontab) Stop() {
	if c.tic == nil {
		return
	}
	c.tic.Stop()
}
