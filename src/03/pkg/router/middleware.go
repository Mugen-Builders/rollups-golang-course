package router

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rollmelette/rollmelette"
)

type Middleware func(interface{}) interface{}

func LoggingMiddleware(handler interface{}) interface{} {
	switch h := handler.(type) {
	case AdvanceHandlerFunc:
		return AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
			res, err := json.Marshal(metadata)
			if err != nil {
				return fmt.Errorf("failed to marshal metadata: %w", err)
			}
			log.Printf("Advance request - metadata: %s", string(res))
			return h(env, metadata, deposit, payload)
		})
	case InspectHandlerFunc:
		return InspectHandlerFunc(func(env rollmelette.EnvInspector, payload []byte) error {
			log.Printf("Inspect request - payload: %s", string(payload))
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
