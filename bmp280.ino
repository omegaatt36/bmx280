/* s-Sense BME280 I2C / s-Sense BMP280 I2C sensor breakout read and environmental calculation example - v1.0/20190524. 
 * 
 * Compatible with:
 *    s-Sense BME280 I2C sensor breakout - temperature, humidity and pressure - [PN: SS-BME280#I2C, SKU: ITBP-6002], info https://itbrainpower.net/sensors/BME280-TEMPERATURE-HUMIDITY-PRESSURE-I2C-sensor-breakout 
 *    s-Sense BMP280 I2C sensor breakout - temperature and pressure - [PN: SS-BMP280#I2C, SKU: ITBP-6001], info https://itbrainpower.net/sensors/BMP280-TEMPERATURE-HUMIDITY-I2C-sensor-breakout
 *
 * BME280/BMP280 environmental sensor, default settings. Read temperature, humidity (unavailable for BMP280) and pressure (pulling at 1sec), then calculate the 
 * environmental data: altitude, dew point (unavailable for BMP280) and equivalent sea level pressure - code based on BME280-2.3.0 library originally written 
 * by Tyler Glenn and forked by Alex Shavlovsky. Some part of code was written by Brian McNoldy. 
 * Amazing work folks! 
 * 
 * We've just select the relevant functions, add some variables, functions and fuctionalities.
 * 
 * 
 * Mandatory wiring:
 *    Common for 3.3V and 5V Arduino boards:
 *        sensor I2C SDA  <------> Arduino I2C SDA
 *        sensor I2C SCL  <------> Arduino I2C SCL
 *        sensor GND      <------> Arduino GND
 *    For Arduino 3.3V compatible:
 *        sensor Vin      <------> Arduino 3.3V
 *    For Arduino 5V compatible:
 *        sensor Vin      <------> Arduino 5V
 * 
 * Leave other sensor PADS not connected.
 * 
 * SPECIAL note for some ARDUINO boards:
 *        SDA (Serial Data)   ->  A4 on Uno/Pro-Mini, 20 on Mega2560/Due, 2 Leonardo/Pro-Micro
 *        SCK (Serial Clock)  ->  A5 on Uno/Pro-Mini, 21 on Mega2560/Due, 3 Leonardo/Pro-Micro
 * 
 * WIRING WARNING: wrong wiring may damage your Arduino board MCU or your sensor! Double check what you've done.
 * 
 * READ BME280 documentation! https://itbrainpower.net/sensors/BME280-TEMPERATURE-HUMIDITY-PRESSURE-I2C-sensor-breakout
 * READ BMP280 documentation! https://itbrainpower.net/sensors/BMP280-TEMPERATURE-PRESSURE-I2C-sensor-breakout
 * 
 * We ask you to use this SOFTWARE only in conjunction with s-Sense BME280 I2C or s-Sense BMP280 I2C sensor breakout usage. Modifications, derivates 
 * and redistribution of this SOFTWARE must include unmodified this notice. You can redistribute this SOFTWARE and/or modify it under the 
 * terms of this notice. 
 * 
 * This SOFTWARE is distributed is provide "AS IS" in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of 
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *  
 * itbrainpower.net invests significant time and resources providing those how to and in design phase of our IoT products.
 * Support us by purchasing our environmental and air quality sensors from https://itbrainpower.net/order#s-Sense
 *  
 *  
 * Dragos Iosub, Bucharest 2019.
 * https://itbrainpower.net
 */

#define SERIAL_SPEED 19200

#include <BMx280_EnvCalc.h>
#include <sSense-BMx280I2C.h>
#include <Wire.h>


BMx280I2C ssenseBMx280;   // Default : forced mode, standby time = 1000 ms
                          // Oversampling = pressure ×1, temperature ×1, humidity ×1, filter off,


//////////////////////////////////////////////////////////////////
void setup()
{
  delay(5000);
  DebugPort.begin(SERIAL_SPEED);

  while(!DebugPort) {} // Wait

  Wire.begin();

  while(!ssenseBMx280.begin())
  {
    DebugPort.println("Could not find BME280 sensor!");
    delay(1000);
  }

  switch(ssenseBMx280.chipModel())
  {
     case BME280::ChipModel_BME280:
       DebugPort.println("Found BME280 sensor! Humidity available.");
       break;
     case BME280::ChipModel_BMP280:
       DebugPort.println("Found BMP280 sensor! No Humidity available.");
       break;
     default:
       DebugPort.println("Found UNKNOWN sensor! Error!");
  }
}

//////////////////////////////////////////////////////////////////
void loop()
{
   printBMx280Data(&DebugPort);
   delay(500);
}

//////////////////////////////////////////////////////////////////
void printBMx280Data( Stream* client )
{
   float temp(NAN), hum(NAN), pres(NAN);

   BME280::TempUnit tempUnit(BME280::TempUnit_Celsius);
   BME280::PresUnit presUnit(BME280::PresUnit_Pa);

   ssenseBMx280.read(pres, temp, hum, tempUnit, presUnit);

   client->print("Temp: ");
   client->print(temp);
   client->print(String(tempUnit == BME280::TempUnit_Celsius ? "C" :"F"));
   client->print("\t\tHumidity: ");
   client->print(hum);
   client->print("% RH");
   client->print("\t\tPressure: ");
   client->print(pres);
   client->print(" Pa");

   BMx280_EnvCalc::AltitudeUnit envAltUnit  =  BMx280_EnvCalc::AltitudeUnit_Meters;
   BMx280_EnvCalc::TempUnit     envTempUnit =  BMx280_EnvCalc::TempUnit_Celsius;

   float altitude = BMx280_EnvCalc::Altitude(pres, envAltUnit);
   float dewPoint = BMx280_EnvCalc::DewPoint(temp, hum, envTempUnit);
   float seaLevel = BMx280_EnvCalc::EquivalentSeaLevelPressure(altitude, temp, pres);
   // seaLevel = BMx280_EnvCalc::SealevelAlitude(altitude, temp, pres); // Deprecated. See EquivalentSeaLevelPressure().

   client->print("\r\nAltitude: ");
   client->print(altitude);
   client->print((envAltUnit == BMx280_EnvCalc::AltitudeUnit_Meters ? "m" : "ft"));
   client->print("\tDew point: ");
   client->print(dewPoint);
   client->print(String(envTempUnit == BMx280_EnvCalc::TempUnit_Celsius ? "C" :"F"));
   client->print("\t\tEquivalent Sea Level Pressure: ");
   client->print(seaLevel);
   client->println(" Pa\r\n");

   delay(1000);
}
