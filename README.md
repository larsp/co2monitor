[![Go Report Card](https://goreportcard.com/badge/github.com/larsp/co2monitor)](https://goreportcard.com/report/github.com/larsp/co2monitor)
[![GoDoc](https://godoc.org/github.com/larsp/co2monitor/meter?status.svg)](https://godoc.org/github.com/larsp/co2monitor/meter)

# CO₂ monitor

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
> cat /etc/udev/rules.d/99-hidraw-permissions.rules 
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", MODE="0664", GROUP="plugdev"
```

## Credit

[Henryk Plötz](https://hackaday.io/project/5301-reverse-engineering-a-low-cost-usb-co-monitor/log/17909-all-your-base-are-belong-to-us)
& [wooga](https://github.com/wooga/office_weather)