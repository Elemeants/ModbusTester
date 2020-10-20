package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/fatih/color"
)

type ModbusSensorType string

const (
	SensorPh   ModbusSensorType = "PH"
	SensorLeaf ModbusSensorType = "LEAF"
	SensorSoil ModbusSensorType = "SOIL"
)

type ModbusSensor struct {
	id         uint
	sensorType ModbusSensorType

	rx_buffer []byte
}

func toUint16(buff []byte) uint16 {
	return uint16(buff[1])<<8 | uint16(buff[0])
}

func (m *ModbusSensor) Id() uint                     { return m.id }
func (m *ModbusSensor) SensorType() ModbusSensorType { return m.sensorType }

func NewModbusSensor(id uint, sensorType ModbusSensorType) *ModbusSensor {
	return &ModbusSensor{id, sensorType, make([]byte, 0)}
}

func (m *ModbusSensor) PrintMeasure() {
	fmt.Print("   SENSOR ID: 	   ")
	color.Green(strconv.FormatUint(uint64(m.rx_buffer[0]), 10))
	fmt.Print("   SENSOR TYPE:    ")
	switch m.rx_buffer[2] {
	case 1:
		color.Green(string(SensorPh))
	case 2:
		color.Green(string(SensorLeaf))
	case 4:
		color.Green(string(SensorSoil))
	default:
		color.Red("Unknown")
	}

	fmt.Println("   SENSOR MEASURE:")

	switch m.sensorType {
	case SensorPh:
		ph := binary.BigEndian.Uint16(m.rx_buffer[3:5])
		fmt.Printf("          PH: %f", float32(ph)/10)
	case SensorLeaf:
		leaf_humidity := binary.BigEndian.Uint16(m.rx_buffer[3:5])
		leaf_temp := binary.BigEndian.Uint16(m.rx_buffer[5:7])
		fmt.Printf("        LEAF: %f°C %f%%", (float32(leaf_temp)/100)-20.0, float32(leaf_humidity)/10)
	case SensorSoil:
		soil_humidity := float32(binary.BigEndian.Uint16(m.rx_buffer[3:5]))
		soil_temp := float32(binary.BigEndian.Uint16(m.rx_buffer[5:7]))
		soil_ec := float32(binary.BigEndian.Uint16(m.rx_buffer[7:9]))
		soil_salinity := float32(binary.BigEndian.Uint16(m.rx_buffer[9:11]))
		fmt.Printf("        SOIL: %f°C %f%% %fuS/cm %fppm", soil_humidity/10, soil_temp/10, soil_ec, soil_salinity)
	}
	fmt.Println()

	len_rx_buffer := len(m.rx_buffer)
	calculatedCRC := CalculateCRC16(m.rx_buffer[0 : len_rx_buffer-2])
	receivedCRC := toUint16(m.rx_buffer[len_rx_buffer-2 : len_rx_buffer])

	fmt.Print("   CRC:            ")
	if calculatedCRC == receivedCRC {
		color.Green("OK")
	} else {
		color.Red("Error expected 0x%04X was 0x%04X", calculatedCRC, receivedCRC)
	}
}

func (m *ModbusSensor) CalcResponseLenght() int {
	const RX_METADATA_LENGHT = 3
	const RX_CRC_LENGHT = 2
	responseLength := RX_METADATA_LENGHT + RX_CRC_LENGHT

	switch m.sensorType {
	case SensorPh:
		responseLength += (0x01 * 2)
	case SensorLeaf:
		responseLength += (0x02 * 2)
	case SensorSoil:
		responseLength += (0x04 * 2)
	}

	return responseLength
}

func (m *ModbusSensor) Read(ref_port *io.ReadWriteCloser) {
	port := *ref_port
	cmd := CreateCommand(m.id, string(m.sensorType))

	fmt.Print(" Sending command ... ")
	PrintBuffer(cmd)

	port.Write(cmd)

	responseLength := m.CalcResponseLenght()

	m.rx_buffer = make([]byte, responseLength)
	n, err := port.Read(m.rx_buffer)

	if err != nil {
		if err == io.EOF {
			color.Red("Error reading data from serial port...")
		}
	} else {
		fmt.Print(" Response ... ")
		PrintBuffer(m.rx_buffer)

		if n != responseLength {
			color.Red("Unexpected data lenght, expected %d, got %d...", responseLength, n)
		} else {
			m.PrintMeasure()
		}
	}
}

func CreateCommand(deviceId uint, deviceType string) []byte {
	base := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	base[0] = byte(deviceId) // Sensor id
	base[1] = 0x03           // Read byte command
	// No of variable to read each variable have 2 bytes so 16 bits
	switch deviceType {
	case string(SensorPh):
		base[5] = 0x01
	case string(SensorLeaf):
		base[5] = 0x02
	case string(SensorSoil):
		base[5] = 0x04
	}

	crc := CalculateCRC16(base[0:6])
	base[6] = byte(crc & 0xFF)
	base[7] = byte(crc >> 8)

	return base
}
