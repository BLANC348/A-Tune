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

package profile

import (
	PB "atune/api/profile"
	"atune/common/client"
	SVC "atune/common/service"
	"atune/common/utils"
	"fmt"
	"io"

	"github.com/urfave/cli"
	CTX "golang.org/x/net/context"
)

var profileCheckCommand = cli.Command{
	Name:      "check",
	Usage:     "check system basic information",
	ArgsUsage: "[arguments...]",
	Description: func() string {
		desc := "\n    check system basic information\n"
		return desc
	}(),
	Action: profileCheck,
}

func init() {
	svc := SVC.ProfileService{
		Name:    "opt.profile.check",
		Desc:    "opt profile system",
		NewInst: newProfileCheckCmd,
	}
	if err := SVC.AddService(&svc); err != nil {
		fmt.Printf("Failed to load profile check service : %s\n", err)
		return
	}
}

func newProfileCheckCmd(ctx *cli.Context, opts ...interface{}) (interface{}, error) {

	return profileCheckCommand, nil
}

func profileCheck(ctx *cli.Context) error {
	appname := ""
	if ctx.NArg() > 2 {
		return fmt.Errorf("only one or zero argument required")
	}
	if ctx.NArg() == 1 {
		appname = ctx.Args().Get(0)
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := PB.NewProfileMgrClient(c.Connection())
	stream, err := svc.CheckInitProfile(CTX.Background(), &PB.ProfileInfo{Name: appname})

	for {
		reply, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		utils.Print(reply)
	}

	return nil
}
