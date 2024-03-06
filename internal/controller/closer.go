package controller

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Closer struct {
	mux   sync.Mutex
	funcs []Func
}

func NewCloser() *Closer {
	return &Closer{
		mux:   sync.Mutex{},
		funcs: make([]Func, 0),
	}
}

func (c *Closer) Add(f Func) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	messages := make([]string, 0, len(c.funcs))
	shutdownDone := make(chan struct{}, 1)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				messages = append(messages, err.Error())
			}
		}

		shutdownDone <- struct{}{}
	}()

	select {
	case <-shutdownDone:
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %w", ctx.Err())
	}

	if len(messages) > 0 {
		return fmt.Errorf("shutdown failed:\n%s", strings.Join(messages, "\n"))
	}

	return nil
}

type Func func(ctx context.Context) error
