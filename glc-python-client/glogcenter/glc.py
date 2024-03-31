import os
import json
import datetime
import requests
import platform
import socket
import netifaces
from decimal import Decimal
import random
import sys


cached_ip = None
server_name = None

def get_server_name():
    global server_name
    if server_name is None:
        server_name = platform.node()
    return server_name

def get_intranet_ip():
    global cached_ip

    if cached_ip is None:
        interfaces = netifaces.interfaces()
        ip_addresses = []

        for interface in interfaces:
            addresses = netifaces.ifaddresses(interface).get(socket.AF_INET)
            if addresses:
                for address in addresses:
                    ip = address['addr']
                    if ip.startswith('192.') or ip.startswith('172.') or ip.startswith('10.'):
                        ip_addresses.append(ip)

        # 对IP地址列表按照特定优先级排序(192优先)
        sorted_ips = sorted(ip_addresses, key=lambda x: (x.startswith('192.'), x.startswith('172.'), x.startswith('10.')))
        if sorted_ips:
            cached_ip = sorted_ips[0]
        else:
            cached_ip = '127.0.0.1'

    return cached_ip

def hash_string(input_str = ''):
    rs = 53653
    i = len(input_str) if input_str is not None else 0
    while i > 0:
        rs = (rs * 33) ^ ord(input_str[i - 1])
        i -= 1
    return str(Decimal(rs & 0xFFFFFFFF).to_eng_string())

# yyyy-MM-dd HH:mm:ss.SSS
def get_current_time():
    current_time = datetime.datetime.now()
    return current_time.strftime('%Y-%m-%d %H:%M:%S.') + str(1000 + current_time.microsecond % 1000)[-3:]

class GlcData:
    def __init__(self, logLevel = ''):
        self.text = ''
        self.date = get_current_time()
        self.system = os.getenv('GLC_SYSTEM')
        self.serverName = get_server_name()
        self.serverIp = get_intranet_ip()
        self.clientIp = ''
        self.traceId = os.getenv('GLC_TRACE_ID')
        self.user = ''


def post_glc_data(glc_data, logLevel):
    if os.getenv('GLC_ENABLE') != 'true':
        return

    url = os.getenv('GLC_API_URL')
    if not url:
        return

    data = {
        'text': glc_data.text,
        'date': glc_data.date,
        'system': glc_data.system,
        'servername': glc_data.serverName,
        'serverip': glc_data.serverIp,
        'clientip': glc_data.clientIp,
        'traceid': glc_data.traceId,
        'loglevel': logLevel,
        'user': glc_data.user
    }
    json_data = json.dumps(data)

    headers = {'Content-Type': 'application/json', 'X-GLC-AUTH': 'glogcenter'}
    glc_api_key = os.getenv('GLC_API_KEY') # X-GLC-AUTH:glogcenter
    if glc_api_key and ':' in glc_api_key:
        key_parts = glc_api_key.split(':')
        headers[key_parts[0]] = key_parts[1]

    requests.post(url, data=json_data, headers=headers)


def argsToGlcData(*args):
    text = ''
    glc_data = None

    # 将非空且非GlcData实例的参数转换为字符串并拼接
    for arg in args:
        if arg is not None and not isinstance(arg, GlcData):
            text += ' ' + str(arg)
    text = text.strip()

    # 无内容时返回空
    if text == '':
        return None

    # 处理GlcData实例参数
    glc_args = [arg for arg in args if isinstance(arg, GlcData)]

    if glc_args:
        # 如果有GlcData实例参数，取出最后一个作为glc_data
        glc_data = glc_args[-1]
    else:
        # 如果没有GlcData实例参数，新建一个GlcdData实例
        glc_data = GlcData()

    # 将第一步得到的text赋值给glc_data的text属性
    glc_data.text = text

    # 相应字段为空时设定默认值
    if glc_data.system == '':
        glc_data.system = 'default'
    if glc_data.traceId == '':
        glc_data.traceId = hash_string(str(random.randint(10000, sys.maxsize)))

    return glc_data

def get_log_level():
    level = os.getenv('GLC_LOG_LEVEL')
    if level is None:
        return 0

    level = level.lower()

    if 'debug' in level:
        return 0
    if 'info' in level:
        return 1
    if 'warn' in level:
        return 2
    if 'error' in level:
        return 3


def debug(*args):
    if get_log_level() <= 0:
        glc_data = argsToGlcData(*args)
        if not glc_data:
            return

        if os.getenv('GLC_ENABLE_CONSOLE_LOG') != 'false':
            print(get_current_time(), '[DEBUG]', glc_data.text)
        post_glc_data(glc_data, 'DEBUG')

def info(*args):
    if get_log_level() <= 1:
        glc_data = argsToGlcData(*args)
        if not glc_data:
            return

        if os.getenv('GLC_ENABLE_CONSOLE_LOG') != 'false':
            print(get_current_time(), ' [INFO]', glc_data.text)
        post_glc_data(glc_data, 'INFO')

def warn(*args):
    if get_log_level() <= 2:
        glc_data = argsToGlcData(*args)
        if not glc_data:
            return

        if os.getenv('GLC_ENABLE_CONSOLE_LOG') != 'false':
            print(get_current_time(), ' [WARN]', glc_data.text)
        post_glc_data(glc_data, 'WARN')

def error(*args):
    if get_log_level() <= 3:
        glc_data = argsToGlcData(*args)
        if not glc_data:
            return

        if os.getenv('GLC_ENABLE_CONSOLE_LOG') != 'false':
            print(get_current_time(), '[ERROR]', glc_data.text)
        post_glc_data(glc_data, 'ERROR')
