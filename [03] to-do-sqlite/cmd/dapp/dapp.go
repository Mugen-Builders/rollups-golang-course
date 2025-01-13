package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Mugen-Builders/to-do-sqlite/configs"
	"github.com/Mugen-Builders/to-do-sqlite/pkg/rollups"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

func DappStrategy(response *rollups.FinishResponse, router *rollups.Router, ih *InspectHandlers) error {
	switch response.Type {
	case "advance_state":
		var data rollups.AdvanceResponse
		if err := json.Unmarshal(response.Data, &data); err != nil {
			return err
		}
		return router.Advance(&data)
	case "inspect_state":
		return ih.TodoInspectHandlers.FindAllTodosHandler()
	}
	return nil
}

func main() {
	// Database setup (SQLite)
	db, err := configs.SetupSQlite("database.db")
	if err != nil {
		errlog.Panicln("Error: could not setup database: ", err)
	}
	infolog.Println("Database setup successful")

	sqlDB, err := db.DB()
	if err != nil {
		errlog.Panicln("Error: could not get database connection: ", err)
	}
	defer sqlDB.Close()

	// Dependency injection with Wire ( It could be done with others libraries like dig, fx, etc )
	ah, err := NewAdvanceHandlers(db)
	if err != nil {
		errlog.Panicln("Failed to initialize advance handlers", "error", err)
	}
	infolog.Println("Advance handlers initialized")

	ih, err := NewInspectHandlers(db)
	if err != nil {
		errlog.Panicln("Failed to initialize inspect handlers", "error", err)
	}
	infolog.Println("Inspect handlers initialized")

	// Router setup and handlers registration
	r := rollups.NewRouter()
	r.HandleAdvance("create_todo", ah.TodoAdvanceHandlers.CreateTodoHandler)
	r.HandleAdvance("update_todo", ah.TodoAdvanceHandlers.UpdateTodoHandler)
	r.HandleAdvance("delete_todo", ah.TodoAdvanceHandlers.DeleteTodoHandler)
	infolog.Println("Router setup successful")

	// Polling loop ( Is there something new to process? )
	finish := rollups.FinishRequest{Status: "accept"}
	for {
		infolog.Println("Sending finish")
		res, err := rollups.SendFinish(&finish)
		if err != nil {
			errlog.Panicln("Error: error making http request: ", err)
		}
		infolog.Println("Received finish status ", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			infolog.Println("No pending rollup request, trying again")
		} else {

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				errlog.Panicln("Error: could not read response body: ", err)
			}

			var response rollups.FinishResponse
			err = json.Unmarshal(resBody, &response)
			if err != nil {
				errlog.Panicln("Error: unmarshaling body:", err)
			}
			finish.Status = "accept"

			// Strategy pattern to handle different types of requests (advance or inspect ?)
			err = DappStrategy(&response, r, ih)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}
