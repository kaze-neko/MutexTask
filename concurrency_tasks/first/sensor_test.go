package main

import (
	"first/sensor"
	"testing"
	"time"
)

func TestSensor(t *testing.T) {
	var sensor = sensor.Sensor{}
	var expected = "Датчик выключен!"
	sensor.On()
	time.Sleep(time.Millisecond * 10)
	sensor.Off()
	var log = sensor.GetLog()
	if len(log) == 0 || log[len(log)-1] != expected {
		t.Errorf("Последнее сообщение в логе \"%s\", ожидалось \"%s\"\n", log[len(log)-1], expected)
	}
}
