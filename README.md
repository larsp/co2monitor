[![Go Report Card](https://goreportcard.com/badge/github.com/larsp/co2monitor)](https://goreportcard.com/report/github.com/larsp/co2monitor)
[![GoDoc](https://godoc.org/github.com/larsp/co2monitor/meter?status.svg)](https://godoc.org/github.com/larsp/co2monitor/meter)

# CO₂ monitor

## Setup & Example
<img src="https://raw.githubusercontent.com/larsp/co2monitor/img/monitor.jpg" alt="Setup" height="300">
<img src="https://raw.githubusercontent.com/larsp/co2monitor/img/dashboard.png" alt="Dashboard" height="300">

## Motivation
Some time ago an [article](https://blog.wooga.com/woogas-office-weather-wow-67e24a5338) about a low cost CO₂ monitor 
came to our attention. A colleague quickly adopted the python [code](https://github.com/wooga/office_weather)
to fit in our prometheus setup. Since humans are sensitive to temperature and CO₂ level, we were now able to 
optimize HVAC settings in our office (Well, we mainly complained to our facility management).

For numerous reasons I wanted to replace the python code with a static Go binary.

## Hardware
- CO₂ meter: Can be found for around 70EUR/USD at [amazon.com](https://www.amazon.com/dp/B00H7HFINS) 
& [amazon.de](https://www.amazon.de/dp/B00TH3OW4Q/). Regardless of minor differences between both devices, both work.
- Some machine which can run the compiled Go binary, has USB and is reachable from your prometheus collector. 
A very first version of a raspberry pi is already sufficient.

## Software
You need prometheus to collect the metrics.

It might make things easier when you set up an `udev` rule e.g.
```bash
$ cat /etc/udev/rules.d/99-hidraw-permissions.rules 
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", MODE="0664", GROUP="plugdev"
```

## Run & Collect

Help
```bash
$ ./co2monitor --help      
usage: co2monitor [<flags>] <device> [<listen-address>]

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Args:
  <device>            CO2 Meter device, such as /dev/hidraw2
  [<listen-address>]  The address to listen on for HTTP requests.
```

Starting the meter export
```bash
$ ./co2monitor /dev/hidraw2
2018/01/18 13:09:31 Serving metrics at ':8080/metrics'
2018/01/18 13:09:31 Device '/dev/hidraw2' opened

```

## Credit

[Henryk Plötz](https://hackaday.io/project/5301-reverse-engineering-a-low-cost-usb-co-monitor/log/17909-all-your-base-are-belong-to-us)
& [wooga](https://github.com/wooga/office_weather)