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
	"time"

	"github.com/Elemeants/ModbusTester/utils"
	"github.com/fatih/color"
	"github.com/jacobsa/go-serial/serial"
	"github.com/spf13/cobra"
)

var mustRepeat bool
var asTable bool
var showAll bool

var pomasSensors []utils.ModbusSensor

func printSensors() {
	fmt.Printf("	|%4s|%10s|\n", "ID", "TYPE")
	fmt.Printf("	|---------------|\n")

	for _, sensor := range pomasSensors {
		fmt.Printf("	|%4d|%10s|\n", sensor.Id(), string(sensor.SensorType()))
	}

	fmt.Printf("	-----------------\n")
}

// pomasCmd represents the test command
var pomasCmd = &cobra.Command{
	Use:   "test",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if showAll {
			printSensors()
			return
		}

		if devicePort == "" {
			log.Fatal("Port must be specified")
		}

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

		color.Yellow("======================= POMAS MODBUS =======================")
		fmt.Printf("    COM:         %s\n", green(devicePort))
		fmt.Printf("    BAUDRATE:    %s\n", green(baudRate))

		fmt.Println(" Opening serial port...")

		port, err := serial.Open(options)
		if err != nil {
			log.Fatalf("serial.Open: %v", err)
		} else {
			// Make sure to close it later.
			defer port.Close()
		}

		for {

			for _, sensor := range pomasSensors {
				if asTable {
					sensor.ReadAsTable(&port)
				} else {
					sensor.Read(&port)
				}
				time.Sleep(time.Millisecond * 500)
			}
			if !mustRepeat {
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pomasCmd)
	pomasCmd.Flags().StringVarP(&devicePort, "port", "p", "", "Port where the test will run")
	pomasCmd.Flags().UintVarP(&baudRate, "baudrate", "b", 9600, "SerialCOM baudrate")

	pomasCmd.Flags().BoolVarP(&showAll, "show", "s", false, "Show all sensors registrated")
	pomasCmd.Flags().BoolVarP(&asTable, "table", "t", false, "Measure as table")
	pomasCmd.Flags().BoolVarP(&mustRepeat, "repeat", "r", false, "Repeat measures")

	// Sensors
	pomasSensors = make([]utils.ModbusSensor, 14)

	pomasSensors[0] = *utils.NewModbusSensor(uint(0x1), utils.SensorSoil)
	pomasSensors[1] = *utils.NewModbusSensor(uint(0x2), utils.SensorSoil)
	pomasSensors[2] = *utils.NewModbusSensor(uint(0x3), utils.SensorSoil)
	pomasSensors[3] = *utils.NewModbusSensor(uint(0x4), utils.SensorSoil)
	pomasSensors[4] = *utils.NewModbusSensor(uint(0x5), utils.SensorSoil)
	pomasSensors[5] = *utils.NewModbusSensor(uint(0x6), utils.SensorSoil)

	pomasSensors[6] = *utils.NewModbusSensor(uint(0x7), utils.SensorPh)
	pomasSensors[7] = *utils.NewModbusSensor(uint(0x8), utils.SensorPh)
	pomasSensors[8] = *utils.NewModbusSensor(uint(0x9), utils.SensorPh)
	pomasSensors[9] = *utils.NewModbusSensor(uint(0x10), utils.SensorPh)
	pomasSensors[10] = *utils.NewModbusSensor(uint(0x11), utils.SensorPh)
	pomasSensors[11] = *utils.NewModbusSensor(uint(0x12), utils.SensorPh)

	// pomasSensors[12] = *utils.NewModbusSensor(uint(0x13), utils.SensorLeaf)
	// pomasSensors[13] = *utils.NewModbusSensor(uint(0x14), utils.SensorLeaf)
}
