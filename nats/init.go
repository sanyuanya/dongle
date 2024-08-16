package nats

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/nats-io/nats.go"
)

func init() {
	opt, err := nats.NkeyOptionFromSeed("./3.nk")
	if err != nil {
		log.Errorf("nkey option from seed error: %v", err)
	}

	nc, err := nats.Connect("nats://81.70.154.116:4222", opt)
	if err != nil {
		log.Errorf("nats connect error: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Errorf("jet stream error: %v", err)
	}

	_, err = js.StreamInfo("dongle")

	if err != nil {
		if err == nats.ErrStreamNotFound {
			_, err := js.AddStream(&nats.StreamConfig{
				Name:      "dongle",
				Subjects:  []string{"dongle.*"},
				Retention: nats.WorkQueuePolicy,
				MaxMsgs:   1000000,
				MaxBytes:  1024 * 1024 * 1024,
				Discard:   nats.DiscardOld,
				Storage:   nats.FileStorage,
				MaxAge:    2 * 60,
			})

			if err != nil {
				log.Errorf("add stream error: %v", err)
			}
		} else {
			log.Errorf("stream info error: %v", err)
		}
	}
}
