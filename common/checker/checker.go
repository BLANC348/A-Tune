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

package checker

import (
	PB "atune/api/profile"
	"atune/common/log"
	"atune/common/utils"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

// MemTopo represent the memory topology type
type MemTopo struct {
	Path string
}

type children struct {
	ID          string `xml:"id"`
	Class       string `xml:"class"`
	Claimed     bool   `xml:"claimed"`
	Handle      string `xml:"handle"`
	Description string `xml:"description"`
	Product     string `xml:"product"`
	Vendor      string `xml:"vendor"`
	Physid      string `xml:"physid"`
	Serial      string `xml:"serial"`
	Slot        string `xml:"slot"`
	Units       string `xml:"units"`
	Size        int64  `xml:"size"`
	Width       int    `xml:"width"`
	Clock       int64  `xml:"clock"`
}

type memorysInfo struct {
	ID        string     `xml:"id"`
	Physid    string     `xml:"physid"`
	Childrens []children `xml:"children"`
}

type topology struct {
	XMLName xml.Name      `xml:"topology"`
	Memorys []memorysInfo `xml:"memorys"`
}

/*
Check method check the memory topolog, whether the memory interpolation is balanced.
*/
func (m *MemTopo) Check(ch chan *PB.AckCheck) error {
	file, err := os.Open(m.Path)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	topology := topology{}
	err = xml.Unmarshal(data, &topology)
	if err != nil {
		return err
	}

	reg := regexp.MustCompile(`DIMM.*?(\d)(\d)(\d)\s.*`)
	var size int64
	memNum := 0
	size = 0

	maxSocket := 0
	maxChannel := 0
	maxSlot := 0

	for _, memory := range topology.Memorys {
		for _, child := range memory.Childrens {
			if params := reg.FindStringSubmatch(child.Slot); params != nil {
				socket, _ := strconv.Atoi(params[1])
				channel, _ := strconv.Atoi(params[2])
				slot, _ := strconv.Atoi(params[3])
				if socket > maxSocket {
					maxSocket = socket
				}
				if channel > maxChannel {
					maxChannel = channel
				}
				if slot > maxSlot {
					maxSlot = slot
				}
			}
		}
	}

	memTotal := (maxSocket + 1) * (maxChannel + 1) * (maxSlot + 1)
	memLocation := make([]bool, memTotal)

	for _, memory := range topology.Memorys {
		for _, child := range memory.Childrens {
			if child.Size != 0 {
				if params := reg.FindStringSubmatch(child.Slot); params != nil {
					socket, _ := strconv.Atoi(params[1])
					channel, _ := strconv.Atoi(params[2])
					slot, _ := strconv.Atoi(params[3])
					index := socket*(maxChannel+1)*(maxSlot+1) + channel*(maxSlot+1) + slot
					memLocation[index] = true
				}

				memNum++
				if size != 0 && child.Size != size {
					sendChanToAdm(ch, "memory sieze", utils.FAILD, "memory size is not the same")
				} else if size == 0 {
					size = child.Size
				}
			}
		}
	}

	log.Infof("memory total num is : %d", memNum)
	log.Info("memory location is:", memNum)

	if memNum == memTotal {
		sendChanToAdm(ch, "memory number", utils.SUCCESS, fmt.Sprintf("memory num is %d", memNum))
		return nil
	}

	if memNum%(maxChannel+1) != 0 {
		sendChanToAdm(ch, "memory number", utils.FAILD, fmt.Sprintf("memory num is %d, not recommend, recommand 8,16 or 32", memNum))
		return nil
	}

	memHalf := memTotal / 2

	for i := 0; i < memHalf; i++ {
		if memLocation[i] != memLocation[i+memHalf] {
			sendChanToAdm(ch, "memory location", utils.FAILD, fmt.Sprintf("memory location maybe not balanced"))
			return nil
		}
	}
	sendChanToAdm(ch, "memory location", utils.SUCCESS, "OK")
	return nil
}

func sendChanToAdm(ch chan *PB.AckCheck, item string, status string, description string) {
	if ch == nil {
		return
	}

	ch <- &PB.AckCheck{Name: item, Status: status, Description: description}
}
