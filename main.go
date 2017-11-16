package main

import (
	"fmt"

	"github.com/johnsudaar/go-bytearena/agent"
	"github.com/johnsudaar/go-bytearena/models"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	log.Info("Starting...")
	a, err := agent.FromEnv("clear_beta")
	if err != nil {
		panic(err)
	}

	c, err := a.Start()
	if err != nil {
		panic(err)
	}

	log.Info("Sending handshake...")
	err = a.Handshake()
	if err != nil {
		panic(err)
	}

	for {
		raw := <-c
		switch v := raw.(type) {
		case models.ErrorEvent:
			log.Panic("Error received: ", v.Error.Error())
		case models.Welcome:
			log.Info("Handshake response received")
		case models.Perception:
			log.Info("Perception received")
			acts := models.Actions{}
			acts = acts.Steer(models.Vector2{0, 1})
			log.Info("Sending actions")
			err := a.Do(acts)
			if err != nil {
				log.WithError(err).Panic("fail to send actions")
			}
		case models.RawEvent:
			log.Debug("Raw event received")
		default:
			log.Warn(fmt.Sprintf("Unexpected event type: %T", raw))
		}
	}
}
