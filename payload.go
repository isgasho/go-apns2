package apns2

import "encoding/json"

type Payload struct {
	Alert Alert
}

func (p *Payload) Map() map[string]interface{} {
	payload := make(map[string]interface{}, 4)

	if !p.Alert.isZero() {
		if p.Alert.isSimple() {
			payload["alert"] = p.Alert.Body
		} else {
			payload["alert"] = p.Alert
		}
	}

	return map[string]interface{}{"aps": payload}
}

func (a Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Map())
}
