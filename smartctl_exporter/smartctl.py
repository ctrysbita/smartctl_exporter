import asyncio
import json
from typing import Dict, List


async def _run(cmd: List[str]) -> str:
    proc = await asyncio.create_subprocess_exec(
        *cmd,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )

    stdout, stderr = await proc.communicate()

    if proc.returncode != 0:
        if stderr:
            print(f'smartctl ended with error: {stderr}')
        if proc.returncode != 64:
            raise Exception(f'smartctl failed with code {proc.returncode}')

    return stdout.decode('utf-8')


async def scan_open() -> List[Dict[str, str]]:
    cmd_result = await _run(['smartctl', '--scan-open', '--json'])
    return json.loads(cmd_result)['devices']


def _get_common(result, cmd_result):
    try:
        result['temperature'] = cmd_result['temperature']['current']
    except:
        pass

    try:
        result['power_on_hours'] = cmd_result['power_on_time']['hours']
    except:
        pass

    try:
        passed = cmd_result['smart_status']['passed']
        result['smart_passed'] = 1 if passed else 0
    except:
        pass


async def get_ata(dev: str) -> Dict[str, str]:
    cmd_result = json.loads(await _run(['smartctl', '--xall', '--json', dev]))
    result = {}

    _get_common(result, cmd_result)

    try:
        attrs: List[Dict] = cmd_result['ata_smart_attributes']['table']
    except:
        return result

    for attr in attrs:
        name = attr['name'].lower().replace('-', '_')
        raw_str = attr['raw']['string']

        try:
            raw = int(raw_str)
        except:
            continue

        result[name] = raw

        if name == 'total_lbas_written':
            result['data_units_written'] = raw

    return result


async def get_nvme(dev: str) -> Dict[str, str]:
    cmd_result = json.loads(await _run(['smartctl', '--xall', '--json', dev]))
    result = {}

    _get_common(result, cmd_result)

    try:
        attrs: Dict = cmd_result['nvme_smart_health_information_log']
    except:
        return result

    for k, v in attrs.items():
        if isinstance(v, list):
            continue

        result[k] = v

    return result
