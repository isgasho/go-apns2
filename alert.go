package apns2

type Alert struct {
	// A short string describing the purpose of the notification.
	// Apple Watch displays this string as part of the notification interface.
	// This string is displayed only briefly and should be crafted so that it can be understood quickly.
	// This key was added in iOS 8.2.
	Title string `json:"title,omitempty"`

	// The text of the alert message.
	Body string `json:"body,omitempty"`

	// The key to a title string in the Localizable.strings file
	// for the current localization. The key string can be formatted
	// with %@ and %n$@ specifiers to take the variables specified in the
	// title-loc-args array. See Localized Formatted Strings for more information.
	// This key was added in iOS 8.2.
	TitleLocKey string `json:"title-loc-key,omitempty"`

	// Variable string values to appear in place of the format specifiers in title-loc-key.
	// See Localized Formatted Strings for more information.
	// This key was added in iOS 8.2.
	TitleLocArgs []string `json:"title-loc-args,omitempty"`

	// If a string is specified, the system displays an alert that includes the Close and View buttons.
	// The string is used as a key to get a localized string in the current localization to use
	// for the right button’s title instead of “View”. See Localized Formatted Strings for more information.
	ActionLocKey string `json:"action-loc-key,omitempty"`

	// A key to an alert-message string in a Localizable.strings file for the current
	// localization (which is set by the user’s language preference).
	// The key string can be formatted with %@ and %n$@ specifiers to take the variables specified in the loc-args array. See Localized Formatted Strings for more information.
	LocKey string `json:"loc-key,omitempty"`

	// Variable string values to appear in place of the format specifiers in loc-key. See Localized Formatted Strings for more information.
	LocArgs []string `json:"loc-args,omitempty"`

	// The filename of an image file in the app bundle; it may include the extension or omit it.
	// The image is used as the launch image when users tap the action button or move the action slider.
	// If this property is not specified, the system either uses the previous snapshot,uses the image
	// identified by the UILaunchImageFile key in the app’s Info.plist file, or falls back to Default.png.
	// This property was added in iOS 4.0.
	LaunchImage string `json:"launch-image,omitempty"`
}

func (a *Alert) isSimple() bool {
	return len(a.Title) == 0 && len(a.TitleLocKey) == 0 && len(a.TitleLocArgs) == 0 && len(a.LocKey) == 0 && len(a.LocArgs) == 0 && len(a.ActionLocKey) == 0 && len(a.LaunchImage) == 0
}

func (a *Alert) isZero() bool {
	return len(a.Body) == 0 && a.isSimple()
}
