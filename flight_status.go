package lufthansa

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type OperatorCarrier struct {
	AirlineID    string `json:"AirlineID"`
	FlightNumber int    `json:"FlightNumber"`
}

type Status struct {
	Code       string `json:"code"`
	Definition string `json:"definition"`
}

type LHTime time.Time

func (lht *LHTime) UnmarshalJSON(data []byte) error {
	ts := string(data)
	if len(ts) < 10 {
		return errors.New(fmt.Sprintf("Invalid time string: %s\n", ts))
	}

	ts = strings.Replace(ts, "\"", "", -1)
	t, err := time.Parse("2006-01-02T15:04Z", ts)
	if err != nil {
		return err
	}

	*lht = LHTime(t)

	return nil
}

type Datetime struct {
	Datetime LHTime `json:"Datetime"`
}

type FlightLeg struct {
	AirportCode      string   `json:"AirportCode"`
	ScheduledTimeUTC Datetime `json:"ScheduledTimeUTC"`
	ActualTimeUTC    Datetime `json:"ActualTimeUTC"`
	TimeStatus       Status   `json:"TimeStatus"`
}

type Flight struct {
	Departure       FlightLeg       `json:"Departure"`
	Arrival         FlightLeg       `json:"Arrival"`
	OperatorCarrier OperatorCarrier `json:"OperatorCarrier"`
	FlightStatus    Status          `json:"FlightStatus"`
}

type Flights struct {
	Flight Flight `json:"Flight"`
}

type FlightStatus struct {
	Flights Flights `json:"Flights"`
}

func (lh *Lufthansa) FlightStatusGet(iata string) (FlightStatus, error) {

	now := time.Now().UTC()

	url := fmt.Sprintf("/operations/flightstatus/%s/%d-%d-%d",
		iata,
		now.Year(),
		now.Month(),
		now.Day())

	lh.sh.SetBearerAuth(lh.token)
	resp, err := lh.sh.Get(url)
	if err != nil {
		return FlightStatus{}, err
	}

	d := struct {
		FlightStatus FlightStatus `json:"FlightStatusResource"`
	}{}

	if err := json.Unmarshal([]byte(resp), &d); err != nil {
		return FlightStatus{}, err
	}

	return d.FlightStatus, nil

}
