package signals

import (
	"os"
	"os/signal"
)

// ListenSignal will listen the given os.Signal until the signal received
// and then execute fc when
func ListenSignal(sg os.Signal, fc func()) {
	scChan := make(chan os.Signal, 1)
	signal.Notify(scChan, sg)
	go func() {
		<-scChan
		signal.Stop(scChan)
		fc()
	}()
}
