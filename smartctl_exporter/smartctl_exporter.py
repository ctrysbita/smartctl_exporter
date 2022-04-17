#!/usr/bin/env python3

import asyncio
from typing import Dict, List

from prometheus_client import Gauge, start_http_server

from smartctl import get_ata, get_nvme, scan_open

devices: List[Dict[str, str]] = []
metrics: Dict[str, Gauge] = {}


async def collect():
    devices = await scan_open()

    for dev in devices:
        name = dev['name']
        if dev['type'] == 'sat':
            result = await get_ata(name)
        elif dev['type'] == 'nvme':
            result = await get_nvme(name)
        else:
            continue

        for k, v in result.items():
            if k not in metrics:
                metrics[k] = Gauge(
                    'smartctl_' + k,
                    'S.M.A.R.T. attribute ' + k,
                    labelnames=['device'],
                )
            metrics[k].labels(name).set(v)


async def run():
    while True:
        await collect()
        await asyncio.sleep(60)


if __name__ == "__main__":
    start_http_server(9111)
    asyncio.run(run())
