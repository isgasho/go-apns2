package apns2

import "encoding/json"

// Payload For each notification, compose a JSON dictionary object (as defined by RFC 4627).
// This dictionary must contain another dictionary identified by the aps key. The aps
// dictionary can contain one or more properties that specify the following user
// notification types: An alert message to display to the user A number
// to badge the app icon with A sound to play
type Payload struct {
	// If this property is included, the system displays a standard alert or a banner, based on the user’s setting.
	Alert Alert

	// The number to display as the badge of the app icon. If this
	// property is absent, the badge is not changed. To remove
	// the badge, set the value of this property to 0.
	Badge uint

	// The name of a sound file in the app bundle or in the Library/Sounds folder of the app’s data container. The sound
	// in this file is played as an alert. If the sound file doesn’t exist or default is specified as the value,
	// the default alert sound is played.The audio must be in one of the audio data formats that are
	// compatible with system sounds.
	Sound string

	// Provide this key with a value of 1 to indicate that new content is available. Including
	// this key and value means that when your app is launched in the background or resumed,
	// application:didReceiveRemoteNotification:fetchCompletionHandler: is called.
	ContentAvailable bool

	// Provide this key with a string value that represents the identifier
	// property of the UIMutableUserNotificationCategory object
	// you created to define custom actions
	Category string
}

// Map returns a valid payload
func (p *Payload) Map() map[string]interface{} {
	payload := make(map[string]interface{}, 4)

	if !p.Alert.isValid() {
		if p.Alert.isSimpleForm() {
			payload["alert"] = p.Alert.Body
		} else {
			payload["alert"] = p.Alert
		}
	}

	if p.Sound != "" {
		payload["sound"] = p.Sound
	}

	if p.ContentAvailable {
		payload["content-available"] = 1
	}

	if p.Category != "" {
		payload["category"] = p.Category
	}

	return map[string]interface{}{"aps": payload}
}

// MarshalJSON returns []byte
func (p Payload) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Map())
}
