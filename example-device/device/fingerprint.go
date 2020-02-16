package device

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/gousb/usbid"
	"github.com/hashicorp/nomad/plugins/device"
	"github.com/google/gousb"
)

// doFingerprint is the long-running goroutine that detects device changes
func (d *SkeletonDevicePlugin) doFingerprint(ctx context.Context, devices chan *device.FingerprintResponse) {
	defer close(devices)

	// Create a timer that will fire immediately for the first detection
	ticker := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(d.fingerprintPeriod)
		}

		d.writeFingerprintToChannel(devices)
	}
}

// fingerprintedDevice is what we "discover" and transform into device.Device objects.
//
// plugin implementations will likely have a native struct provided by the corresonding SDK
type fingerprintedDevice struct {
	Vendor     string
	Product    string
	Address    int
	PCIBusID   string
}

// writeFingerprintToChannel collects fingerprint info, partitions devices into
// device groups, and sends the data over the provided channel.
func (d *SkeletonDevicePlugin) writeFingerprintToChannel(devices chan<- *device.FingerprintResponse) {

	usbCtx := gousb.NewContext()
	defer usbCtx.Close()

	// OpenDevices is used to find the devices to open.
	d.logger.Info("fingerprinting USB devices")
	foundDevs, err := usbCtx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// The usbid package can be used to print out human readable information.
		d.logger.Info("examining device", "vendor_id", desc.Vendor, "product_id", desc.Product)
		return uint16(desc.Vendor) == d.targetVendorId && uint16(desc.Product) == d.targetProductId
	})

	if err != nil {
		d.logger.Warn("error opening usb devices", "err", err)
		devices <- &device.FingerprintResponse{
			Devices: nil,
			Error:   err,
		}
		return
	}

	d.deviceLock.Lock()
	defer d.deviceLock.Unlock()

	discoveredDevices := make([]*fingerprintedDevice, 0, len(foundDevs))
	for _, d := range foundDevs {
		v, p := vendorAndProduct(d.Desc)
		discoveredDevices = append(discoveredDevices, &fingerprintedDevice{
			Vendor: v,
			Product: p,
			Address:    d.Desc.Address,
			PCIBusID: strconv.Itoa(d.Desc.Bus),
		})
	}

	// short-cut for our diff
	if d.devices != nil && len(discoveredDevices) == len(d.devices) {
		d.logger.Info("no changes to devices")
		return
	}

	// during fingerprinting, devices are grouped by "device group" in
	// order to facilitate scheduling
	// devices in the same device group should have the same
	// Vendor, Type, and Name ("Model")
	// Build Fingerprint response with computed groups and send it over the channel
	d.devices = map[string]string{}
	deviceListByDeviceName := make(map[string][]*fingerprintedDevice)
	for _, device := range discoveredDevices {
		deviceName := device.Product
		deviceListByDeviceName[deviceName] = append(deviceListByDeviceName[deviceName], device)
		d.devices[strconv.Itoa(device.Address)] = deviceName
	}

	// Build Fingerprint response with computed groups and send it over the channel
	deviceGroups := make([]*device.DeviceGroup, 0, len(deviceListByDeviceName))
	for groupName, devices := range deviceListByDeviceName {
		deviceGroups = append(deviceGroups, deviceGroupFromFingerprintData(groupName, devices))
	}
	d.logger.Info("sending device fingerprint", "num_devices", len(discoveredDevices))
	devices <- device.NewFingerprint(deviceGroups...)
}

// deviceGroupFromFingerprintData composes deviceGroup from a slice of detected devicers
func deviceGroupFromFingerprintData(groupName string, deviceList []*fingerprintedDevice) *device.DeviceGroup {
	// deviceGroup without devices makes no sense -> return nil when no devices are provided
	if len(deviceList) == 0 {
		return nil
	}

	devices := make([]*device.Device, 0, len(deviceList))
	for _, dev := range deviceList {
		devices = append(devices, &device.Device{
			ID:      strconv.Itoa(dev.Address),
			Healthy: true,
			HwLocality: &device.DeviceLocality{
				PciBusID: dev.PCIBusID,
			},
		})
	}

	deviceGroup := &device.DeviceGroup{
		Vendor: deviceList[0].Vendor,
		Type:    "usb",
		Name:    groupName,
		Devices: devices,
		Attributes: nil,
	}
	return deviceGroup
}

func vendorAndProduct(desc *gousb.DeviceDesc) (vendor, product string) {
	vendor = fmt.Sprintf("unknown(%v)", desc.Vendor)
	product = fmt.Sprintf("unknown(%v)", desc.Product)
	if v, ok := usbid.Vendors[desc.Vendor]; ok {
		vendor = v.Name
		if d, ok := v.Product[desc.Product]; ok {
			product = d.Name
		}
	}
	return
}
