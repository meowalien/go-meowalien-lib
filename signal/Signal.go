package signal

import (
	"os"
	"os/signal"
)

func ListenSignal(sg os.Signal, fc func()) {
	scChan := make(chan os.Signal, 1)
	signal.Notify(scChan, sg)
	go func() {
		<-scChan
		signal.Stop(scChan)
		fc()
	}()
}
