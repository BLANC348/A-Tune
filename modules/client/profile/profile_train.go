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
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	CTX "golang.org/x/net/context"
)

var trainCommand = cli.Command{
	Name:      "train",
	Usage:     "train a new model with self collected data",
	UsageText: "atune-adm train OPTIONS",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "data_path,d",
			Usage: "the path where has data for training",
			Value: "",
		},
		cli.StringFlag{
			Name:  "output_file,o",
			Usage: "the model to be generated",
			Value: "",
		},
	},
	Description: func() string {
		desc := `
	 training a new model with the self collected data, data_path option specified
	 the path that storage the collected data, the collected data must have more 
	 than two workload type. output_file specified the file path where to store
	 the trained model, which must be end with .m.
	     example: atune-adm train --data_path=./data --output_file=./model/trained.m`
		return desc
	}(),
	Action: train,
}

func init() {
	svc := SVC.ProfileService{
		Name:    "opt.profile.train",
		Desc:    "opt profile system",
		NewInst: newTrainCmd,
	}
	if err := SVC.AddService(&svc); err != nil {
		fmt.Printf("Failed to load collection service : %s\n", err)
		return
	}
}

func newTrainCmd(ctx *cli.Context, opts ...interface{}) (interface{}, error) {
	return trainCommand, nil
}

func checkTrainCtx(ctx *cli.Context) error {
	dataPath := ctx.String("data_path")
	if dataPath == "" {
		cli.ShowCommandHelp(ctx, "train")
		return fmt.Errorf("error: data_path must be specified")
	}

	outputPath := ctx.String("output_file")
	if outputPath == "" {
		cli.ShowCommandHelp(ctx, "train")
		return fmt.Errorf("error: output_file must be specified")
	}

	return nil
}

func train(ctx *cli.Context) error {
	if err := checkTrainCtx(ctx); err != nil {
		return err
	}

	dataPath, err := filepath.Abs(ctx.String("data_path"))
	if err != nil {
		return err
	}
	outputPath, err := filepath.Abs(ctx.String("output_file"))
	if err != nil {
		return err
	}

	dir := filepath.Dir(outputPath)

	exist, err := utils.PathExist(dir)
	if err != nil {
		return err
	}
	if !exist {
		os.MkdirAll(dir, 0660)
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := PB.NewProfileMgrClient(c.Connection())
	stream, err := svc.Training(CTX.Background(), &PB.TrainMessage{DataPath: dataPath, OutputPath: outputPath})

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
	fmt.Println("the model generate path:", outputPath)
	return nil
}
