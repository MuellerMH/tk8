package cluster

import (
	"reflect"
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

func TestAWS_GetAbsPath(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			if got := aws.GetAbsPath(tt.args.filePath); got != tt.want {
				t.Errorf("AWS.GetAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWS_GetConfig(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  string
		want2  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			got, got1, got2 := aws.GetConfig()
			if got != tt.want {
				t.Errorf("AWS.GetConfig() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AWS.GetConfig() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("AWS.GetConfig() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestAWS_DistSelect(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			got, got1 := aws.DistSelect()
			if got != tt.want {
				t.Errorf("AWS.DistSelect() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AWS.DistSelect() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_distSelect(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := distSelect()
			if got != tt.want {
				t.Errorf("distSelect() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("distSelect() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAWS_GetCredentials(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	tests := []struct {
		name   string
		fields fields
		want   AwsCredentials
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			if got := aws.GetCredentials(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AWS.GetCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWS_GetClusterConfig(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	tests := []struct {
		name   string
		fields fields
		want   ClusterConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			if got := aws.GetClusterConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AWS.GetClusterConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAWS_Create(t *testing.T) {
	type fields struct {
		Dists     map[string]DistOS
		Ec2IP     string
		Namespace string
		OSHelper  oshelper.OSHelper
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := &AWS{
				Dists:     tt.fields.Dists,
				Ec2IP:     tt.fields.Ec2IP,
				Namespace: tt.fields.Namespace,
				OSHelper:  tt.fields.OSHelper,
			}
			aws.Create()
		})
	}
}
