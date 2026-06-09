package devices

import (
	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/openapi"
)

// Device is the domain struct for a registered device.
type Device struct {
	ID          string
	Name        string
	UDID        string
	DeviceClass DeviceClass
	Model       string
	Platform    Platform
	Status      DeviceStatus
}

// DeviceClass is an alias for the generated openapi.DeviceClass type.
type DeviceClass = openapi.DeviceClass

// Platform is an alias for the generated openapi.BundleIdPlatform type.
type Platform = openapi.BundleIdPlatform

// DeviceStatus is an alias for the generated openapi.DeviceStatus type.
type DeviceStatus = openapi.DeviceStatus

func fromResponse(d openapi.Device) *Device {
	if d.Attributes == nil {
		return &Device{ID: d.Id}
	}
	return &Device{
		ID:          d.Id,
		Name:        ptr.Deref(d.Attributes.Name),
		UDID:        ptr.Deref(d.Attributes.Udid),
		DeviceClass: ptr.Deref(d.Attributes.DeviceClass),
		Model:       ptr.Deref(d.Attributes.Model),
		Platform:    ptr.Deref(d.Attributes.Platform),
		Status:      ptr.Deref(d.Attributes.Status),
	}
}
