/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/Elemeants/ModbusTester/utils"
	"github.com/fatih/color"
	"github.com/jacobsa/go-serial/serial"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "sensor",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen).SprintFunc()

		options := serial.OpenOptions{
			PortName:               devicePort,
			BaudRate:               baudRate,
			DataBits:               8,
			StopBits:               1,
			MinimumReadSize:        4,
			InterCharacterTimeout:  2000,
			Rs485Enable:            true,
			Rs485RtsHighDuringSend: true,
		}

		color.Yellow("======================= TEST MODBUS =======================")
		fmt.Printf("    COM:         %s\n", green(devicePort))
		fmt.Printf("    BAUDRATE:    %s\n", green(baudRate))
		fmt.Printf("    SENSOR ID:   %s\n", green(sensorId))
		fmt.Printf("    SENSOR TYPE: ")

		switch sensorType {
		case string(utils.SensorPh):
			color.Green("PH")
		case string(utils.SensorLeaf):
			color.Green("LEAF")
		case string(utils.SensorSoil):
			color.Green("SOIL")
		default:
			color.Red("UNKNOWN")
			log.Fatal("Unknown sensor type")
		}

		fmt.Println(" Opening serial port...")

		port, err := serial.Open(options)
		if err != nil {
			log.Fatalf("serial.Open: %v", err)
		} else {
			// Make sure to close it later.
			defer port.Close()
		}

		sensor := utils.NewModbusSensor(sensorId, utils.ModbusSensorType(sensorType))
		sensor.Read(&port)

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&devicePort, "port", "p", "", "Port where the test will run")
	testCmd.Flags().UintVarP(&baudRate, "baudrate", "b", 9600, "SerialCOM baudrate")

	testCmd.Flags().UintVarP(&sensorId, "id", "i", 1, "Modbus Sensor ID")
	testCmd.Flags().StringVarP(&sensorType, "type", "t", "PH", "Modbus Sensor Type ['pH' : 'Leaf' : 'Soil']")

	testCmd.MarkFlagRequired("port")
}
