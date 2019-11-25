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
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
	CTX "golang.org/x/net/context"
	yaml "gopkg.in/yaml.v2"
)

var profileTunningCommand = cli.Command{
	Name:      "tuning",
	Usage:     "dynamic bayes search optimal parameter sets",
	ArgsUsage: "PROJECT_YAML",
	Description: func() string {
		desc := `
	 tuning command usning bayes method dynamic search optimal parameter sets,
	 the PROJECT_YAML which you can refer to Documentation example.yaml.
	     example: atune-adm tuning ./example.yaml
	`
		return desc
	}(),
	Action: profileTunning,
}

func init() {
	svc := SVC.ProfileService{
		Name:    "opt.profile.tunning",
		Desc:    "opt profile system",
		NewInst: newProfileTuningCmd,
	}
	if err := SVC.AddService(&svc); err != nil {
		fmt.Printf("Failed to load profile analysis service : %s\n", err)
		return
	}
}

func newProfileTuningCmd(ctx *cli.Context, opts ...interface{}) (interface{}, error) {

	return profileTunningCommand, nil
}

func profileTunning(ctx *cli.Context) error {
	if err := utils.CheckArgs(ctx, 1, utils.ConstExactArgs); err != nil {
		return err
	}

	yamlPath := ctx.Args().Get(0)
	exist, err := utils.PathExist(yamlPath)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("project file %s is not exist", yamlPath)
	}

	if !strings.HasSuffix(yamlPath, ".yaml") && !strings.HasSuffix(yamlPath, ".yml") {
		return fmt.Errorf("project file is not ends with yaml or yml")
	}

	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	var yamlStruct interface{}
	if err := yaml.Unmarshal(data, &yamlStruct); err != nil {
		return err
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := PB.NewProfileMgrClient(c.Connection())
	stream, err := svc.Tuning(CTX.Background(), &PB.ProfileInfo{Name: string(data)})

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
