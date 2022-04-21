package main

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func Smartctl(args ...string) gjson.Result {
	args = append([]string{"--json=c"}, args...)
	cmd := exec.Command("smartctl", args...)
	output, _ := cmd.Output()
	return gjson.ParseBytes(output)
}

type Device struct {
	Name     string
	Type     string
	Protocol string
}

func GetDevices() []*Device {
	devices := []*Device{}
	Smartctl("--scan-open").Get("devices").ForEach(func(k, v gjson.Result) bool {
		device := &Device{}
		json.Unmarshal([]byte(v.Raw), device)
		devices = append(devices, device)
		return true
	})
	return devices
}

type Result struct {
	ModelName       string `json:"model_name"`
	SerialNumber    string `json:"serial_number"`
	FirmwareVersion string `json:"firmware_version"`
	Passed          bool
	Attributes      map[string]float64
}

func GetAll(dev *Device) *Result {
	j := Smartctl("--xall", dev.Name)

	r := &Result{}
	json.Unmarshal([]byte(j.Raw), r)

	r.Passed = j.Get("smart_status.passed").Bool()

	r.Attributes = map[string]float64{}
	r.Attributes["temperature"] = j.Get("temperature.current").Float()
	r.Attributes["power_on_hours"] = j.Get("power_on_time.hours").Float()

	switch dev.Type {
	case "nvme":
		j.Get("nvme_smart_health_information_log").ForEach(func(k, v gjson.Result) bool {
			if v.Type == gjson.JSON {
				return true
			}
			r.Attributes[k.Str] = v.Float()
			return true
		})
	case "sat":
		j.Get("ata_smart_attributes.table").ForEach(func(k, v gjson.Result) bool {
			name := strings.ReplaceAll(strings.ToLower(v.Get("name").Str), "-", "_")
			rawStr := v.Get("raw.string").Str

			rawFloat, err := strconv.ParseFloat(rawStr, 64)
			if err != nil {
				rawStr = strings.TrimSpace(rawStr[:strings.IndexByte(rawStr, '(')])
				if rawFloat, err = strconv.ParseFloat(rawStr, 64); err != nil {
					return true
				}
			}

			r.Attributes[name] = rawFloat

			if name == "total_lbas_written" {
				r.Attributes["data_units_written"] = rawFloat
			}

			return true
		})
	}

	return r
}
