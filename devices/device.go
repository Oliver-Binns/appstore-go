package devices

import "github.com/oliver-binns/appstore-go/openapi"

// Device is the domain struct for a registered device.
type Device struct {
	ID          string
	Name        string
	UDID        string
	DeviceClass DeviceClass
	Model       string
	OS          string
	Platform    Platform
	Status      DeviceStatus
}

// DeviceClass is an alias for the generated openapi.DeviceClass type.
type DeviceClass = openapi.DeviceClass

// DeviceClass constants re-exported with idiomatic Go names.
const (
	AppleTV        = openapi.APPLETV
	AppleVisionPro = openapi.APPLEVISIONPRO
	AppleWatch     = openapi.APPLEWATCH
	IPad           = openapi.IPAD
	IPhone         = openapi.IPHONE
	IPod           = openapi.IPOD
	Mac            = openapi.MAC
)

// Platform is an alias for the generated openapi.BundleIdPlatform type.
type Platform = openapi.BundleIdPlatform

// Platform constants re-exported with idiomatic Go names.
const (
	IOS   = openapi.IOS
	MacOS = openapi.MACOS
)

// DeviceStatus is an alias for the generated openapi.DeviceStatus type.
type DeviceStatus = openapi.DeviceStatus

// DeviceStatus constants re-exported with idiomatic Go names.
const (
	Enabled  = openapi.ENABLED
	Disabled = openapi.DISABLED
)

func fromResponse(d openapi.Device) *Device {
	return &Device{
		ID:          d.Id,
		Name:        derefString(d.Attributes.Name),
		UDID:        derefString(d.Attributes.Udid),
		DeviceClass: derefDeviceClass(d.Attributes.DeviceClass),
		Model:       derefString(d.Attributes.Model),
		OS:          derefString(d.Attributes.Os),
		Platform:    derefPlatform(d.Attributes.Platform),
		Status:      derefDeviceStatus(d.Attributes.Status),
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefDeviceClass(c *openapi.DeviceClass) openapi.DeviceClass {
	if c == nil {
		return ""
	}
	return *c
}

func derefPlatform(p *openapi.BundleIdPlatform) openapi.BundleIdPlatform {
	if p == nil {
		return ""
	}
	return *p
}

func derefDeviceStatus(s *openapi.DeviceStatus) openapi.DeviceStatus {
	if s == nil {
		return ""
	}
	return *s
}
