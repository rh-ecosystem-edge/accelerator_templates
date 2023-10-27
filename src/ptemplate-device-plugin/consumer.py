#!/usr/bin/python

"""
    A Toy program to read and write the current time to the /dev/ptemplate-X devices
    the device(s) to be used are passed as environemnt variables
    e.g. 
        export ptemplate-0=/dev/ptemplate-0
        export ptemplate-4=/dev/ptemplate-4
    and the location of devfs can be passed in the DEVFS envvar
        export DEVFS=/dev
    this allows for a container to map the hosts /dev to /hostdevinside the container
"""

import os
import os.path
import sys
import time

ENVVAR="ptemplate"

ENVVAR_LEN = len(ENVVAR)
DEVFS = os.getenv("DEVFS", "/hostdev") + "/"

## get the devices to use from their envvars
raw_devices = [ os.environ[i]  for i in os.environ if i[:ENVVAR_LEN] == ENVVAR ]

## envvars point to devices in /dev/ by default so we need to repoint them to the value of DEVFS
devices = [DEVFS + os.path.basename(i) for i in raw_devices]

if len(devices) == 0:
    sys.stdout.write(f"failed to read envvar {ENVVAR}\n")
    sys.exit(1)

while True:
    for dev in devices:
        with open(dev, "r", encoding="utf8") as f:
            content = f.read()
            sys.stdout.write(f"read from {dev}: {content}\n")

        with open(dev, "w", encoding="utf8") as f:
            sys.stdout.write(f"writing to {dev}: {time.time()}\n")
            f.write(f"write time {time.time()}")

    sys.stdout.write("sleep for 10 seconds\n\n")
    sys.stdout.flush()
    time.sleep(10)
