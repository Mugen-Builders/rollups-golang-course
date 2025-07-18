package main

import (
	"dapp/rollups"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	state string
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

func HandleAdvance(data *rollups.AdvanceResponse) error {
	dataMarshal, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("HandleAdvance: failed marshaling json: %w", err)
	}
	infolog.Println("Received advance request data", string(dataMarshal))

	decodedPayload, err := hex.DecodeString(data.Payload[2:])
	if err != nil {
		return fmt.Errorf("Handler: Error decoding payload: %w", err)
	}
	state = strings.ToUpper(string(decodedPayload))
	infolog.Printf("%s - to upper: %s\n", decodedPayload, state)
	return nil
}

func HandleInspect(data *rollups.InspectResponse) error {
	dataMarshal, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("HandleInspect: failed marshaling json: %w", err)
	}
	infolog.Println("Received inspect request data", string(dataMarshal))
	rollups.SendReport(&rollups.ReportRequest{
		Payload: rollups.Str2Hex(state),
	})
	return nil
}

func Handler(response *rollups.FinishResponse) error {
	var err error

	switch response.Type {
	case "advance_state":
		data := new(rollups.AdvanceResponse)
		if err = json.Unmarshal(response.Data, data); err != nil {
			return fmt.Errorf("Handler: Error unmarshaling advance: %w", err)
		}
		err = HandleAdvance(data)
	case "inspect_state":
		data := new(rollups.InspectResponse)
		if err = json.Unmarshal(response.Data, data); err != nil {
			return fmt.Errorf("Handler: Error unmarshaling inspect: %w", err)
		}
		err = HandleInspect(data)
	}
	return err
}

func main() {
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
			err = Handler(&response)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}
