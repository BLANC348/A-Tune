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
The sub class of the Configurator, used to change the /proc/sys/* config.
"""

import sys
import logging
import subprocess
import re

if __name__ == "__main__":
    sys.path.insert(0, "./../../")
from configurator.common import *

logger = logging.getLogger(__name__)


class Sysctl(Configurator):
    """To change the /proc/sys/* config"""
    _module = "SYSCTL"
    _submod = "SYSCTL"

    def __init__(self, user=None):
        Configurator.__init__(self, user)
        self.__cmd = "sysctl"

    def _get(self, key):
        with open('/dev/null', 'w') as no_print:
            output = subprocess.check_output("{cmd} -n {key}".format(
                cmd=self.__cmd, key=key).split(),
                stderr=no_print)
        return output.decode()

    def _set(self, key, value):
        with open('/dev/null', 'w') as no_print:
            return subprocess.call('{cmd} -w {key}="{val}"'.format(
                cmd=self.__cmd, key=key, val=value), shell=True,
                stdout=no_print, stderr=no_print)

    def _check(self, config1, config2):
        config1 = re.sub("\s{1,}", " ", config1)
        config2 = re.sub("\s{1,}", " ", config2)
        return config1 == config2


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print('usage: ' + sys.argv[0] + ' key=value')
        sys.exit(-1)
    ct = Sysctl("UT")
    print(ct.set(sys.argv[1]))
    print(ct.get(ct._getcfg(sys.argv[1])[0]))
