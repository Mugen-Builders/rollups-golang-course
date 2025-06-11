package router

import (
	"fmt"
	"log"

	"github.com/rollmelette/rollmelette"
)

type Middleware func(interface{}) interface{}

func LoggingMiddleware(handler interface{}) interface{} {
	switch h := handler.(type) {
	case AdvanceHandlerFunc:
		return AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
			log.Printf("Advance request - Sender: %s, Deposit: %s, Payload: %s", metadata.MsgSender.String(), deposit.String(), string(payload))
			return h(env, metadata, deposit, payload)
		})
	case InspectHandlerFunc:
		return InspectHandlerFunc(func(env rollmelette.EnvInspector, payload []byte) error {
			log.Printf("Inspect request - Payload: %s", string(payload))
			return h(env, payload)
		})
	default:
		return handler
	}
}

func ErrorHandlingMiddleware(handler interface{}) interface{} {
	switch h := handler.(type) {
	case AdvanceHandlerFunc:
		return AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
			err := h(env, metadata, deposit, payload)
			if err != nil {
				env.Report([]byte(fmt.Sprintf("Error: %v", err)))
			}
			return err
		})
	case InspectHandlerFunc:
		return InspectHandlerFunc(func(env rollmelette.EnvInspector, payload []byte) error {
			err := h(env, payload)
			if err != nil {
				env.Report([]byte(fmt.Sprintf("Error: %v", err)))
			}
			return err
		})
	default:
		return handler
	}
}
