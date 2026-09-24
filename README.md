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


## Network discovery (fxcontroller)

Building and testing requires Go 1.25 or newer.

The service advertises its HTTP API through Zeroconf (mDNS/DNS-SD) by
default. A future fxcontroller can browse **`_fx._tcp` in `local.`** once
to find fxaudio, fxpixel, fxdmx, and fxtrigger. No controller is required
to run the services, and direct HTTP API access remains available if
multicast registration fails (a warning is logged).

```yaml
discovery:
  enabled: true
  name: "" # Optional friendly instance name, 1–63 UTF-8 bytes
  id: ""   # Optional unique, stable installation ID
```

The default name is `service-hostname-port`; the default ID is
`service:hostname:port`. Give each installation a distinct `discovery.id`
if identity must survive hostname or port changes. Names must also be
unique on the LAN. TXT values must fit the DNS-SD 255-byte record limit.

The discovery contract is the same across all four projects:

| DNS record | Meaning |
| --- | --- |
| SRV / A / AAAA | HTTP hostname, actual listening port, and addresses |
| TXT `txtvers=1` | Discovery metadata format version |
| TXT `service` | `fxaudio`, `fxpixel`, `fxdmx`, or `fxtrigger` |
| TXT `id` | Installation identity |
| TXT `api=v1` | API version |
| TXT `scheme=http` | API transport |
| TXT `path=/v1` | API base path |

Controllers should use the resolved SRV port and address, dispatch by
`service`, check supported metadata/API versions, and verify API reachability
separately. Discovery is not authentication; treat advertisements as untrusted
network input. This adds service advertisements, not a controller or UI.

Advertisements start only after the HTTP listener binds and are withdrawn
on SIGINT/SIGTERM or when the HTTP server stops. Discovery uses UDP 5353
multicast on the local network segment; firewalls, Wi-Fi client isolation,
VLANs, and container networking can prevent discovery. Containers typically
need host networking (where supported) or an mDNS relay. No internet service
or separately installed Bonjour/Avahi daemon is required.

On macOS, inspect live advertisements with:

```sh
dns-sd -B _fx._tcp local.
# Resolve one of the instance names returned above:
dns-sd -L "INSTANCE NAME" _fx._tcp local.
```

On Linux with Avahi tools: `avahi-browse -rt _fx._tcp`.
