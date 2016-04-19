package tesla

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CommandResponse struct {
	Response struct {
		Reason string `json:"reason"`
		Result bool   `json:"result"`
	} `json:"response"`
}

func (v Vehicle) AutoparkForward() error {
	return v.autoPark("start_forward")
}

func (v Vehicle) AutoparkReverse() error {
	return v.autoPark("start_reverse")
}

func (v Vehicle) autoPark(action string) error {
	driveState, _ := v.DriveState()
	data := url.Values{}
	data.Set("vehicle_id", strconv.Itoa(v.VehicleID))
	data.Add("lat", strconv.FormatFloat(driveState.Latitude, 'f', 6, 64))
	data.Add("lon", strconv.FormatFloat(driveState.Longitude, 'f', 6, 64))
	data.Add("action", action)

	u, _ := url.ParseRequestURI(BaseURL)
	u.Path = "/api/1/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/autopark_request"
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+ActiveClient.Token.AccessToken)
	req.Header.Set("Accept", "application/json")

	_, err = client.Do(req)
	return nil
}

func (v Vehicle) TriggerHomelink() error {
	// apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/trigger_homelink"
	return nil
}

func (v Vehicle) Wakeup() (*Vehicle, error) {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/wake_up"
	body, err := sendCommand(apiUrl, nil)
	if err != nil {
		return nil, err
	}
	vehicleResponse := &VehicleResponse{}
	err = json.Unmarshal(body, vehicleResponse)
	if err != nil {
		return nil, err
	}
	return vehicleResponse.Response, nil
}

func (v Vehicle) OpenChargePort() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_port_door_open"
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v Vehicle) SetChargeLimitStandard() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_standard"
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v Vehicle) SetChargeLimitMax() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_max_range"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// func (v Vehicle) SetChargeLimit(limit int) error {
// 	url := BaseURL + "/vehicles/" + strconv.Itoa(v.VehicleID) + "/command/set_charge_limit?=" + strconv.Itoa(limit)
// 	_, err := v.Client.postURLEncoded(url, nil)
// 	return err
// }

func (v Vehicle) StartCharging() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_start"
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v Vehicle) StopCharging() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_stop"
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v Vehicle) FlashLights() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/flash_lights"
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v *Vehicle) HonkHorn() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/honk_horn"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// func (v Vehicle) UnlockDoors() error {
// 	apiUrl := BaseURL + "/vehicles/" + strconv.Itoa(v.VehicleID) + "/command/unlock_doors"
// 	_, err := v.Client.postURLEncoded(apiUrl, nil)
// 	return err
// }

func (v Vehicle) LockDoors() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/door_lock"
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v Vehicle) SetTemprature(driver float64, passenger float64) error {
	driveTemp := strconv.FormatFloat(driver, 'f', -1, 32)
	passengerTemp := strconv.FormatFloat(passenger, 'f', -1, 32)
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/set_temps?driver_temp=" + driveTemp + "&passenger_temp=" + passengerTemp
	_, err := ActiveClient.post(apiUrl, nil)
	return err
}

func (v Vehicle) StartAirConditioning() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/auto_conditioning_start"
	_, err := sendCommand(url, nil)
	return err
}

func (v Vehicle) StopAirConditioning() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/auto_conditioning_stop"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// func (v Vehicle) MovePanoRoof(state string, percent int) error {
// 	apiUrl := BaseURL + "/vehicles/" + strconv.Itoa(v.VehicleID) + "/command/sun_roof_control?"
// 	_, err := v.Client.postURLEncoded(apiUrl, nil)
// 	return err
// }

func (v Vehicle) Start(password string) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/remote_start_drive?password=" + password
	_, err := sendCommand(apiUrl, nil)
	return err
}

func (v Vehicle) OpenTrunk(trunk string) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/trunk_open" // ?which_trunk=" + trunk
	theJson := `{"which_trunk": "` + trunk + `"}`
	_, err := ActiveClient.post(apiUrl, []byte(theJson))
	return err
}

func sendCommand(url string, reqBody []byte) ([]byte, error) {
	body, err := ActiveClient.post(url, reqBody)
	if err != nil {
		return nil, err
	}
	response := &CommandResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if response.Response.Result != true && response.Response.Reason != "" {
		return nil, errors.New(response.Response.Reason)
	}
	return body, nil
}
