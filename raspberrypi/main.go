package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

func main() {
	// Create new connection to i2c-bus on 1 line with address 0x76.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x76, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()
	// Uncomment next line to supress verbose output
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	sensor, err := bsbmp.NewBMP(bsbmp.BMP280, i2c)
	if err != nil {
		log.Fatal(err)
	}
	// Uncomment next line to supress verbose output
	logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)

	ctx := context.Background()
	ch := make(chan os.Signal, 1)

	go func(ctx context.Context) {
		for {
			ctx.Done()

			// Read temperature in celsius degree
			t, err := sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
			if err != nil {
				log.Fatal(err)
			}
			// Read atmospheric pressure in pascal
			pPa, err := sensor.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
			if err != nil {
				log.Fatal(err)
			}
			// Read atmospheric altitude in meters above sea level, if we assume
			// that pressure at see level is equal to 101325 Pa.
			a, err := sensor.ReadAltitude(bsbmp.ACCURACY_STANDARD)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Temprature = %vC\tPressure = %v Pa\tAltitude = %v m\n", t, pPa, a)

			time.Sleep(time.Second)
		}
	}(ctx)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("terminating...")

	ctx.Done()
	log.Println("exit")
}
