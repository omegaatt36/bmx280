from pico_i2c_lcd import I2cLcd
from machine import I2C
from machine import Pin
import utime as time
import bme280_float as bme280
 
 
i2c1 = I2C(id=1, scl=Pin(3), sda=Pin(2),freq=100000)
lcd = I2cLcd(i2c1, 0x27, 2, 16)

i2c2 = I2C(id=0, scl=Pin(1), sda=Pin(0))
bme=bme280.BME280(i2c=i2c2)

while True:
      temperature = bme.values[0]         #reading the value of temperature
      pressure = bme.values[1]            #reading the value of pressure
      humidity = bme.values[2]            #reading the value of humidity
      lcd.move_to(0, 0)
      lcd.putstr('H:' + str(humidity))
      lcd.move_to(8, 0)
      lcd.putstr('T:' + str(temperature))
      lcd.move_to(0, 1)
      lcd.putstr('P:' + str(pressure))
      time.sleep(1)
