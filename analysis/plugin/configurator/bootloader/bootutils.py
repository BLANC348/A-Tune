#!/usr/bin/python3
# -*- coding: utf-8 -*-
# Copyright (c) 2019 Huawei Technologies Co., Ltd.
# A-Tune is licensed under the Mulan PSL v1.
# You can use this software according to the terms and conditions of the Mulan PSL v1.
# You may obtain a copy of Mulan PSL v1 at:
#     http://license.coscl.org.cn/MulanPSL
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v1 for more details.
# Create: 2019-10-29

"""
Utils class.
"""


class Utils:
    """Utils class"""

    @staticmethod
    def get_keypos(str, key):
        keys = [" " + key + "=", " " + key + " ", " " + key + "\n"]
        ret = []
        for k in keys:
            pos = str.rfind(k)
            if pos != -1:
                pos += 1
            ret.append(pos)
        return max(ret)

    @staticmethod
    def get_value(key):
        with open("/proc/cmdline", 'r') as f:
            active_cmd = f.read()
        keypos = Utils.get_keypos(active_cmd, key)
        active = None
        if keypos != -1:
            config = active_cmd[keypos:].split()[0]
            if config.find("=") != -1:
                active = config.split("=")[1]
        return active
