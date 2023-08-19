# fxtrigger [![Build and release](https://github.com/danesparza/fxtrigger/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxtrigger/actions/workflows/release.yaml)
REST service for Raspberry Pi GPIO/Sensor -> webhooks.  Made with ❤️ for makers, DIY craftsmen, prop makers and professional soundstage designers everywhere

## Prerequisites
fxtrigger uses Raspberry Pi GPIO to listen for input pin button presses or sensor events.  You'll need to make sure those buttons and sensors are wired up and working before using fxTrigger to connect those triggers to your webhook endpoints.

For motion sensing, I would recommend using the [Adafruit PIR (motion) sensor](https://www.adafruit.com/product/189) as well -- just get a [Pi with headers](https://www.adafruit.com/product/3708) and connect the PIR to power, ground, and to a GPIO data pin  (and be sure to follow the [PIR motion sensor guide](https://learn.adafruit.com/pir-passive-infrared-proximity-motion-sensor/) for the board).  In fxtrigger, specifiy the GPIO pin you hook it up to (not the physical pin) when create the trigger.  See the [Raspberry Pi Pinout interactive reference](https://pinout.xyz/#) for more information. 


## Installing
Installing fxtrigger is also really simple.  Grab the .deb file from the [latest release](https://github.com/danesparza/fxtrigger/releases/latest) and then install it using dpkg:


```bash
sudo dpkg -i fxtrigger-1.0.40_armhf.deb 
````

This automatically installs the **fxtrigger** service with a default configuration and starts the service. 

You can then use the service at http://localhost:3020

See the REST API documentation at http://localhost:3020/v1/swagger/

## Removing 
Uninstalling is just as simple:

```bash
sudo dpkg -r fxtrigger
````
