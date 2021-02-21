# LinuxMonitorControl

Linux monitor control using [DDC](https://en.wikipedia.org/wiki/Display_Data_Channel) (your monitor must support it)

Only brightness and contrast are supported.

Requires: [ddcutil](https://www.ddcutil.com/)

## Build & Run

```bash
make build

# Set brightness on all monitors
./LinuxMonitorControl -b 50

# Set brightness only on monitor 1
./LinuxMonitorControl -b 50 -d 1

# Set contrast only on monitor 1
./LinuxMonitorControl -c 50 -d 1
```


## Build local ddcutil

See required dependencies [here](https://www.ddcutil.com/building/)

```bash
make dependencies

# If build succeeds
./ddcutil/bin/ddcutil detect
```
