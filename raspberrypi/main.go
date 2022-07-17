package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":9110", "The address to listen on for HTTP requests.")
)

type bme280Collector struct {
	temperatureMetric *prometheus.Desc
	humidityMetric    *prometheus.Desc
	pressureMetric    *prometheus.Desc
	altitudeMetric    *prometheus.Desc

	sensor *bsbmp.BMP
}

func newBme280Collector(bme280 *bsbmp.BMP) *bme280Collector {
	return &bme280Collector{
		temperatureMetric: prometheus.NewDesc(
			"bme280_temperature", "bme 280 temperature(Celsius)",
			nil, nil),
		humidityMetric: prometheus.NewDesc(
			"bme280_humidity", "bme 280 humidity(RH)",
			nil, nil),
		pressureMetric: prometheus.NewDesc(
			"bme280_pressure", "bme 280 pressure(hPa)",
			nil, nil),
		altitudeMetric: prometheus.NewDesc(
			"bme280_altitude", "bme 280 altitude(m)",
			nil, nil),

		sensor: bme280,
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *bme280Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.temperatureMetric
	ch <- collector.pressureMetric
	ch <- collector.altitudeMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *bme280Collector) Collect(ch chan<- prometheus.Metric) {
	// Read temperature in celsius degree
	temperature, err := collector.sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		log.Fatalln(err)
		return
	}
	// Read humidity in RH
	_, humidity, err := collector.sensor.ReadHumidityRH(bsbmp.ACCURACY_STANDARD)
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
	mHumidity := prometheus.MustNewConstMetric(collector.humidityMetric,
		prometheus.GaugeValue, float64(humidity))
	mPressure := prometheus.MustNewConstMetric(collector.pressureMetric,
		prometheus.GaugeValue, float64(pressure))
	mAltitude := prometheus.MustNewConstMetric(collector.altitudeMetric,
		prometheus.GaugeValue, float64(altitude))
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mTemperature)
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mHumidity)
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mPressure)
	ch <- prometheus.NewMetricWithTimestamp(time.Now(), mAltitude)
}

func main() {
	flag.Parse()

	// Create new connection to i2c-bus on 1 line with address 0x76.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x76, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()
	// Uncomment next line to supress verbose output
	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	sensor, err := bsbmp.NewBMP(bsbmp.BME280, i2c)
	if err != nil {
		log.Fatal(err)
	}
	// Uncomment next line to supress verbose output
	logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)

	collector := newBme280Collector(sensor)
	prometheus.MustRegister(collector)

	log.Printf("handle %s/metrics \n", *addr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))

	log.Println("exit")
}
