//go:build !windows

package logger

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (p *Logger) WatchReset() {
	for {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)
		<-sigs
		if err := p.Reset(); err != nil {
			log.Println(err)
		}
	}
}
