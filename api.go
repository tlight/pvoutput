package pvoutput

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Frame struct {
	Date string  `json:"d"`   // Output Date yyyymmdd 20060102
	Time string  `json:"t"`   // Time hh:mm time 14:00 15:04
	V1   int64   `json:"v1"`  // Energy Generation  1 watt hours 10000
	V2   int64   `json:"v2"`  // Power Generation number watts 2000
	V3   int64   `json:"v3"`  // Energy Consumption  number watt hours 10000
	V4   int64   `json:"v4"`  // Power Consumption  watts 2000
	V5   float32 `json:"v5"`  // Temperature  // decimal celsius 23.4
	V6   float32 `json:"v6"`  // Voltage decimal volts 239.2
	C1   int64   `json:"c1"`  // Cumulative Flag 1 or 0 number
	N    int64   `json:"n"`   // Net Flag 1 or 0
	V7   int64   `json:"v7"`  // user defined
	V8   int64   `json:"v8"`  // user defined
	V9   int64   `json:"v9"`  // user defined
	V10  int64   `json:"v10"` // user defined
	V11  int64   `json:"v11"` // user defined
	V12  int64   `json:"v12"` // user defined
	M1   string  `json:"m1"`  // user defined
}

// 60 requests per hour.		  1 minute
// 300 per hour in Donation mode  5 s
type API struct {
	APIKey     string
	SystemId   string
	LastUpdate int64
}

func NewAPI(APIKey, SystemId string) *API {
	api := &API{
		APIKey:   APIKey,
		SystemId: SystemId,
	}
	return api
}

func NewFrame() *Frame {
	now := time.Now()
	frame := &Frame{
		Date: now.Format("20060102"),
		Time: now.Format("15:04"),
	}
	return frame
}

// curl -d "d=20111201" -d "t=10:00" -d "v1=1000" -d "v2=150" -H "X-Pvoutput-Apikey: Your-API-Key" -H "X-Pvoutput-SystemId: Your-System-Id" https://pvoutput.org/service/r2/addstatus.jsp
func (api *API) AddStatus(frame *Frame) (string, error) {
	uri := "https://pvoutput.org/service/r2/addstatus.jsp"
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	form := url.Values{
		"d":  {frame.Date},
		"t":  {frame.Time},
		"v2": {fmt.Sprintf("%d", frame.V2)},
		"v4": {fmt.Sprintf("%d", frame.V4)},
		"v6": {fmt.Sprintf("%f", frame.V6)},
	}

	req, err := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Pvoutput-Apikey", api.APIKey)
	req.Header.Add("X-Pvoutput-SystemId", api.SystemId)

	resp, err := client.Do(req)
	api.LastUpdate = time.Now().Unix()

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil

	// read X-Rate-Limit-Remaining - remaining for hour
	// X-Rate-Limit-Limit  total requests for hour
	// X-Rate-Limit-Reset - unix time when reset
}
