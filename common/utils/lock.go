/*
 * Copyright (c) 2019 Huawei Technologies Co., Ltd.
 * A-Tune is licensed under the Mulan PSL v1.
 * You can use this software according to the terms and conditions of the Mulan PSL v1.
 * You may obtain a copy of Mulan PSL v1 at:
 *     http://license.coscl.org.cn/MulanPSL
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v1 for more details.
 * Create: 2019-10-29
 */

package utils

import (
	"sync/atomic"
)

const (
	locked int32 = 1
)

// MutexLock : the type implement the TryLock function
type MutexLock struct {
	state int32
}

// TryLock method if get the lock success, otherwise return false
func (m *MutexLock) TryLock() bool {
	if !atomic.CompareAndSwapInt32(&m.state, 0, 1) {
		return false
	}
	return true
}

// Unlock method unlock lock flag
func (m *MutexLock) Unlock() {
	atomic.CompareAndSwapInt32(&m.state, 1, 0)

}

// IsLocked method return wether the lock is already locked
func (m *MutexLock) IsLocked() bool {
	if atomic.LoadInt32(&m.state) == locked {
		return true
	}

	return false

}
