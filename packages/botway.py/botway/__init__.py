"""
botway.py

Python client package for Botway.
"""

__author__ = 'abdfnx'
__license__ = 'MIT'
__copyright__ = 'Copyright (c) 2022-now Abdfn'
__version__ = '0.0.5'

import yaml
import json
from os import path
from pathlib import Path

botwayStream = open(path.join(Path().home(), '.botway', 'botway.json'), 'r')
botwayConfigData = json.load(botwayStream)

botStream = open('.botway.yaml', 'r')
botConfigData = yaml.load(botStream, Loader=yaml.FullLoader)

def find(d, i):
    if i in d:
        yield d[i]

    for k, v in d.items():
        if isinstance(v, dict):
            for i in find(v, i):
                yield i

def getBotInfo(value):
    for val in find(botConfigData, 'bot'):
        return val[value]

if getBotInfo('lang') != 'python':
    raise RuntimeError('ERROR: Botway is not running in Python')

def GetToken():
    for val in find(botwayConfigData, 'botway'):
        return val['bots'][getBotInfo('name')]['bot_token']

def GetAppId():
    if getBotInfo('type') == 'slack':
        for val in find(botwayConfigData, 'botway'):
            return val['bots'][getBotInfo('name')]['bot_app_token']
    else:
        for val in find(botwayConfigData, 'botway'):
            return val['bots'][getBotInfo('name')]['bot_app_id']

def GetGuildId(serverName):
    if getBotInfo('type') != 'discord':
        raise RuntimeError('ERROR: This function/feature is only working with discord bots.')
    else:
        for val in find(botwayConfigData, 'botway'):
            return val['bots'][getBotInfo('name')]['guilds'][serverName]['server_id']

def GetSigningSecret():
    if getBotInfo('type') != 'slack':
        raise RuntimeError('ERROR: This function/feature is only working with slack bots.')
    else:
        for val in find(botwayConfigData, 'botway'):
            return val['bots'][getBotInfo('name')]['signing_secret']
