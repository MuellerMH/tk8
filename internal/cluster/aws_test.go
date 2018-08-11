// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"testing"

	"github.com/kubernauts/tk8/internal/cluster/oshelper"
)

func TestAWS_CreateFileFromTemplate(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	type args struct {
		templateName   string
		targetFileName string
		awsInstanceOS  string
		data           interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"test_variables",
			fields{
				map[string]DistOS{
					"centos": DistOS{
						User:     "centos",
						AmiOwner: "688023202711",
						OS:       "dcos-centos7-*",
					},
				},
				"",
				"testcase",
				oshelper.NewOSHelper(),
			},
			args{
				"configs/templates/kubespray-aws-variables.tf", "testcase/variables.tf", "centos", nil,
			},
			true,
		},
		{
			"test_infra",
			fields{
				map[string]DistOS{
					"centos": DistOS{
						User:     "centos",
						AmiOwner: "688023202711",
						OS:       "dcos-centos7-*",
					},
				},
				"",
				"testcase",
				oshelper.NewOSHelper(),
			},
			args{
				"configs/templates/kubespray-aws-create-infra.tf", "testcase/create-infrastructer.tf", "centos", nil,
			},
			true,
		},
		{
			"test_terraform",
			fields{
				map[string]DistOS{
					"centos": DistOS{
						User:     "centos",
						AmiOwner: "688023202711",
						OS:       "dcos-centos7-*",
					},
				},
				"",
				"testcase",
				oshelper.NewOSHelper(),
			},
			args{
				"configs/templates/kubespray-aws-terraform.tfvars", "testcase/terraform.tfvars", "centos", nil,
			},
			true,
		},
		{
			"test_credentials",
			fields{
				map[string]DistOS{
					"centos": DistOS{
						User:     "centos",
						AmiOwner: "688023202711",
						OS:       "dcos-centos7-*",
					},
				},
				"",
				"testcase",
				oshelper.NewOSHelper(),
			},
			args{
				"configs/templates/kubespray-aws-credentials.tfvars", "testcase/credentials.tfvars", "centos", nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			switch tt.name {
			case "test_terraform":
				tt.args.data = aws.GetClusterConfig()
			case "test_credentials":
				tt.args.data = aws.GetCredentials()
			}

			if got := aws.CreateFileFromTemplate(tt.args.templateName, tt.args.targetFileName, tt.args.awsInstanceOS, tt.args.data); got != tt.want {
				t.Errorf("AWS.CreateFileFromTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
