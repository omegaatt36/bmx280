package main

import (
	"log"
	"net/http"
	"time"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type bmp280Collector struct {
	temperatureMetric *prometheus.Desc
	pressureMetric    *prometheus.Desc
	altitudeMetric    *prometheus.Desc

	sensor *bsbmp.BMP
}

func newBmp280Collector(bmp280 *bsbmp.BMP) *bmp280Collector {
	return &bmp280Collector{
		temperatureMetric: prometheus.NewDesc(
			"bmp280_temperature", "bmp 280 temperature(Celsius)",
			nil, nil),
		pressureMetric: prometheus.NewDesc(
			"bmp280_pressure", "bmp 280 pressure(hPa)",
			nil, nil),
		altitudeMetric: prometheus.NewDesc(
			"bmp280_altitude", "bmp 280 altitude(m)",
			nil, nil),

		sensor: bmp280,
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *bmp280Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.temperatureMetric
	ch <- collector.pressureMetric
	ch <- collector.altitudeMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *bmp280Collector) Collect(ch chan<- prometheus.Metric) {
	// Read temperature in celsius degree
	temperature, err := collector.sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		log.Fatalln(err)
		return
	}
	// Read atmospheric pressure in pascal
	pressure, err := collector.sensor.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		log.Fatalln(err)
		return
	}
	// Read atmospheric altitude in meters above sea level, if we assume
	// that pressure at see level is equal to 101325 Pa.
	altitude, err := collector.sensor.ReadAltitude(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		log.Fatalln(err)
		return
	}

	mTemperature := prometheus.MustNewConstMetric(collector.temperatureMetric,
		prometheus.GaugeValue, float64(temperature))
	mPressure := prometheus.MustNewConstMetric(collector.pressureMetric,
		prometheus.GaugeValue, float64(pressure))
	mAltitude := prometheus.MustNewConstMetric(collector.altitudeMetric,
		prometheus.GaugeValue, float64(altitude))
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mTemperature)
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mPressure)
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mAltitude)
}

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

	collector := newBmp280Collector(sensor)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9110", nil))

	log.Println("exit")
}
