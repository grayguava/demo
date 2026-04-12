#!/usr/bin/env python3
"""formseal-sync - fetch encrypted form submissions"""

import sys
import os

if os.name == 'nt':
    try:
        os.system('chcp 65001 >nul')
    except:
        pass

try:
    sys.stdout.reconfigure(encoding='utf-8')
    sys.stderr.reconfigure(encoding='utf-8')
except:
    pass

from commands import config as cmd_config
from commands import fetch as cmd_fetch
from commands import status as cmd_status
from commands import help as cmd_help


def main():
    args = sys.argv[1:]
    command = args[0] if args else None

    match command:
        case 'config':
            cmd_config.run()
        case 'set':
            cmd_config.set(args[1:])
        case 'fetch':
            cmd_fetch.run()
        case 'status':
            cmd_status.run()
        case 'help' | '--help' | '-h':
            cmd_help.run()
        case None:
            cmd_help.run()
        case _:
            print(f'Unknown command: {command}')
            print('Run formseal-sync help for usage')


if __name__ == '__main__':
    main()