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
	d   string  // Output Date yyyymmdd 20060102
	t   string  // Time hh:mm time 14:00 15:04
	v1  int     // Energy Generation  1 watt hours 10000
	v2  int     // Power Generation number watts 2000
	v3  int     // Energy Consumption  number watt hours 10000
	v4  int     // Power Consumption  watts 2000
	v5  float32 // Temperature  // decimal celsius 23.4
	v6  float32 // Voltage decimal volts 239.2
	c1  int     // Cumulative Flag 1 or 0 number
	n   int     // Net Flag 1 or 0
	v7  int     // user defined
	v8  int     // user defined
	v9  int     // user defined
	v10 int     // user defined
	v11 int     // user defined
	v12 int     // user defined
	m1  string  // user defined
}

type API struct {
	APIKey     string
	SystemId   string
	LastUpdate int
}

// 60 requests per hour.		  1 minute
// 300 per hour in Donation mode  5 s

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
		d: now.Format("20060102"),
		t: now.Format("15:04"),
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
		"d":  {frame.d},
		"t":  {frame.t},
		"v2": {fmt.Sprintf("%d", frame.v2)},
		"v4": {fmt.Sprintf("%d", frame.v4)},
	}

	req, err := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Pvoutput-Apikey", api.APIKey)
	req.Header.Add("X-Pvoutput-SystemId", api.SystemId)

	resp, err := client.Do(req)
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
