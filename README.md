Reading mavlink v1 and sending LTM Telemetry to serial

```
Usage: 
  -baud int
        Example: -baud 2400 (default 2400)
  -mavlink string
        Example: -mavlink 127.0.0.1:14550 (default ":14550")
  -serial string
        Example: -serial /dev/cu.SLAB_USBtoUART
```


Example read from :14550 and send to /dev/cu.SLAB_USBtoUART at 2400baud
```
./mavlink2ltm -serial /dev/cu.SLAB_USBtoUART
```