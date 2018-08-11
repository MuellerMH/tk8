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
	"C"
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/kubernauts/tk8/internal/cluster/oshelper"
	"github.com/spf13/viper"
)
import (
	"path"
	"path/filepath"
)

// AWS is the main structer of the platform controller
type AWS struct {
	Dists     map[string]DistOS
	Ec2IP     string
	Namespace string
	OSHelper  oshelper.OSHelper
}

type AwsCredentials struct {
	AwsAccessKeyID   string
	AwsSecretKey     string
	AwsAccessSSHKey  string
	AwsDefaultRegion string
}

func GetRootPath() string {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := path.Dir(e)
	fmt.Println(path)
	return path
}

// CreateFileFromTemplate create config files from templates
func (aws *AWS) CreateFileFromTemplate(templateName string, targetFileName string, awsInstanceOS string, data interface{}) bool {
	absPath, _ := filepath.Abs("../../" + templateName)
	_ = os.Mkdir(aws.Namespace, os.ModePerm)
	file, err := os.Create(targetFileName)
	if err != nil {
		aws.OSHelper.FatalLog("Cannot create file", err)
		return false
	}
	defer file.Close()
	template := template.Must(template.ParseFiles(absPath))
	if err != nil {
		aws.OSHelper.FatalLog(templateName, "for", awsInstanceOS, "could not parsed")
		return false
	}
	if data == nil {
		template.Execute(file, aws.Dists[awsInstanceOS])
		return true
	}
	template.Execute(file, data)
	return true
}

// GetConfig configs from viper
func (aws *AWS) GetConfig() (string, string, string) {
	//Read Configuration File
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/tk8")
	viper.AddConfigPath("./../..")
	verr := viper.ReadInConfig() // Find and read the config file
	if verr != nil {             // Handle errors reading the config file
		log.Fatal("config could not readed")
		return "", "", ""
	}
	awsAmiID := viper.GetString("aws.ami_id")
	awsInstanceOS := viper.GetString("aws.os")
	sshUser := viper.GetString("aws.ssh_user")
	return awsAmiID, awsInstanceOS, sshUser
}

// DistSelect choose the Dist and return sshUser and osLabel
func (aws *AWS) DistSelect() (string, string) {

	awsAmiID, awsInstanceOS, sshUser := aws.GetConfig()

	if awsAmiID != "" && sshUser != "" {
		awsInstanceOS = "custom"
		aws.Dists[awsInstanceOS] = DistOS{
			User:     sshUser,
			AmiOwner: "",
			OS:       awsAmiID,
		}
	}

	// TODO: clean up debug
	aws.OSHelper.Log(awsInstanceOS)

	if awsInstanceOS == "" && awsAmiID == "" {
		log.Fatal("Provide either of AMI ID or OS in the config file.")
		return "", ""
	}
	if awsAmiID != "" && sshUser == "" {
		log.Fatal("SSH Username is required when using custom AMI")
		return "", ""
	}
	// prepare config
	if !aws.CreateFileFromTemplate("/configs/templates/kubespray-aws-variables.tf", "./"+aws.Namespace+"/variables.tf", awsInstanceOS, nil) {
		return "", ""
	}
	if !aws.CreateFileFromTemplate("/configs/templates/kubespray-aws-create-infra.tf", "./"+aws.Namespace+"/create-infrastructure.tf", awsInstanceOS, nil) {
		return "", ""
	}

	return sshUser, awsInstanceOS
}

func (aws *AWS) GetCredentials() AwsCredentials {
	//Read Configuration File
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/tk8")
	viper.AddConfigPath("./../..")
	verr := viper.ReadInConfig() // Find and read the config file
	if verr != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", verr))
	}
	return AwsCredentials{
		AwsAccessKeyID:   viper.GetString("aws.aws_access_key_id"),
		AwsSecretKey:     viper.GetString("aws.aws_secret_access_key"),
		AwsAccessSSHKey:  viper.GetString("aws.aws_ssh_keypair"),
		AwsDefaultRegion: viper.GetString("aws.aws_default_region"),
	}
}

func (aws *AWS) GetClusterConfig() ClusterConfig {
	//Read Configuration File
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/tk8")
	viper.AddConfigPath("./../..")
	verr := viper.ReadInConfig() // Find and read the config file
	if verr != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", verr))
	}
	return ClusterConfig{
		AwsClusterName:               viper.GetString("aws.clustername"),
		AwsVpcCidrBlock:              viper.GetString("aws.aws_vpc_cidr_block"),
		AwsCidrSubnetsPrivate:        viper.GetString("aws.aws_cidr_subnets_private"),
		AwsCidrSubnetsPublic:         viper.GetString("aws.aws_cidr_subnets_public"),
		AwsBastionSize:               viper.GetString("aws.aws_bastion_size"),
		AwsKubeMasterNum:             viper.GetString("aws.aws_kube_master_num"),
		AwsKubeMasterSize:            viper.GetString("aws.aws_kube_master_size"),
		AwsEtcdNum:                   viper.GetString("aws.aws_etcd_num"),
		AwsEtcdSize:                  viper.GetString("aws.aws_etcd_size"),
		AwsKubeWorkerNum:             viper.GetString("aws.aws_kube_worker_num"),
		AwsKubeWorkerSize:            viper.GetString("aws.aws_kube_worker_size"),
		AwsElbAPIPort:                viper.GetString("aws.aws_elb_api_port"),
		K8sSecureAPIPort:             viper.GetString("aws.k8s_secure_api_port"),
		KubeInsecureApiserverAddress: viper.GetString("aws."),
	}
}

func (aws *AWS) Create() {
	if !aws.OSHelper.CheckDependency("terraform") {
		return
	}

	_, err := aws.OSHelper.Shell("terraform", "version")
	if err != nil {
		return // cancel process
	}
	if _, err := aws.OSHelper.FileInfo("/configs/templates/credentials.tfvars"); err != nil {
		if !aws.CreateFileFromTemplate("/configs/templates/kubespray-aws-credentials.tfvars", "./"+aws.Namespace+"/credentials.tfvars", aws.Namespace, aws.GetCredentials()) {
			return // cancel process
		}
	}
	if !aws.CreateFileFromTemplate("/configs/templates/kubespray-aws-terraform.tfvars", "./"+aws.Namespace+"/terraform.tfvars", aws.Namespace, aws.GetClusterConfig()) {
		return // cancel process
	}

	//TODO: extract to a builder same like oshelper
	terrInit, _ := aws.OSHelper.Shell("terraform", "init")
	terrInit.Dir = "./" + aws.Namespace + "/"
	out, _ := terrInit.StdoutPipe()
	terrInit.Start()
	scanInit := bufio.NewScanner(out)
	for scanInit.Scan() {
		m := scanInit.Text()
		aws.OSHelper.Log(m)
	}

	terrInit.Wait()

	terrSet, _ := aws.OSHelper.Shell("terraform", "apply", "-var-file=credentials.tfvars", "-auto-approve")

	terrSet.Dir = "./" + aws.Namespace + "/"
	stdout, _ := terrSet.StdoutPipe()
	terrSet.Stderr = terrSet.Stdout
	terrSet.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		aws.OSHelper.Log(m)
	}

	terrSet.Wait()
	os.Exit(0)

}

// NewAWS is the AWS Constructor
func NewAWS(namespace string, distOS map[string]DistOS, oshelper oshelper.OSHelper) AWS {

	aws := AWS{Namespace: namespace, Dists: distOS, OSHelper: oshelper}
	return aws
}
func AWSInstall() {
	// check if ansible is installed
	terr, err := exec.LookPath("ansible")
	if err != nil {
		log.Fatal("Ansible command not found, kindly check")
	}
	fmt.Printf("Found Ansible at %s\n", terr)
	rr, err := exec.Command("ansible", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(rr))

	//Start Kubernetes Installation

	//check if ansible host file exists

	if _, err := os.Stat("./kubespray/inventory/hosts"); err != nil {
		fmt.Println("./kubespray/inventory/host inventory file not found")
		os.Exit(1)
	}

	// Copy the configuraton files as indicated in the kubespray docs

	if _, err := os.Stat("./kubespray/inventory/awscluster"); err == nil {
		fmt.Println("Configuration folder already exists")
	} else {
		//os.MkdirAll("./kubespray/inventory/awscluster/group_vars", 0755)
		exec.Command("cp", "-rfp", "./kubespray/inventory/sample/", "./kubespray/inventory/awscluster/").Run()

		exec.Command("cp", "./kubespray/inventory/hosts", "./kubespray/inventory/awscluster/hosts").Run()

		//Enable load balancer api access and copy the kubeconfig file locally
		loadBalancerName, err := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./kubespray/inventory/hosts | cut -d'=' -f2").CombinedOutput()
		if err != nil {
			fmt.Println("Problem getting the load balancer domain name", err)
		} else {
			//Make a copy of kubeconfig on Ansible host
			f, err := os.OpenFile("./kubespray/inventory/awscluster/group_vars/k8s-cluster.yml", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()
			fmt.Fprintf(f, "kubeconfig_localhost: true\n")

			g, err := os.OpenFile("./kubespray/inventory/awscluster/group_vars/all.yml", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				panic(err)
			}

			defer g.Close()

			// Resolve Load Balancer Domain Name and pick the first IP

			s, _ := exec.Command("sh", "-c", "grep apiserver_loadbalancer_domain_name= ./kubespray/inventory/hosts | cut -d'=' -f2 | sed 's/\"//g'").CombinedOutput()
			// Convert the Domain name to string and strip all spaces so that Lookup does not return errors
			r := string(s)
			t := strings.TrimSpace(r)

			fmt.Println(t)
			node, err := net.LookupHost(t)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			ec2IP := node[0]

			fmt.Println(node)

			DomainName := strings.TrimSpace(string(loadBalancerName))
			loadBalancerDomainName := "apiserver_loadbalancer_domain_name: " + DomainName

			fmt.Fprintf(g, "#Set cloud provider to AWS\n")
			fmt.Fprintf(g, "cloud_provider: 'aws'\n")
			fmt.Fprintf(g, "#Load Balancer Configuration\n")
			fmt.Fprintf(g, "loadbalancer_apiserver_localhost: false\n")
			fmt.Fprintf(g, "%s\n", loadBalancerDomainName)
			fmt.Fprintf(g, "loadbalancer_apiserver:\n")
			fmt.Fprintf(g, "  address: %s\n", ec2IP)
			fmt.Fprintf(g, "  port: 6443\n")
		}
	}
	sshUser, osLabel := distSelect()
	kubeSet := exec.Command("ansible-playbook", "-i", "./inventory/awscluster/hosts", "./cluster.yml", "--timeout=60", "-e ansible_user="+sshUser, "-e bootstrap_os="+osLabel, "-b", "--become-user=root", "--flush-cache")
	kubeSet.Dir = "./kubespray/"
	stdout, _ := kubeSet.StdoutPipe()
	kubeSet.Stderr = kubeSet.Stdout
	kubeSet.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		//log.Printf(m)
	}

	kubeSet.Wait()

	os.Exit(0)
}

func AWSDestroy() {
	// check if terraform is installed
	terr, err := exec.LookPath("terraform")
	if err != nil {
		log.Fatal("Terraform command not found, kindly check")
	}
	fmt.Printf("Found terraform at %s\n", terr)
	rr, err := exec.Command("terraform", "version").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(rr))

	// Remove ssh bastion file

	if _, err := os.Stat("./kubespray/ssh-bastion.conf"); err == nil {
		os.Remove("./kubespray/ssh-bastion.conf")
	}

	// Remove the cluster inventory folder
	err = os.RemoveAll("./kubespray/inventory/awscluster")
	if err != nil {
		fmt.Println(err)
	}

	// Check if credentials file exist, if it exists skip asking to input the AWS values
	if _, err := os.Stat("./kubespray/contrib/terraform/aws/credentials.tfvars"); err == nil {
		fmt.Println("Credentials file already exists, creation skipped")
	} else {

		fmt.Println("Please enter your AWS access key ID")
		var awsAccessKeyID string
		fmt.Scanln(&awsAccessKeyID)

		fmt.Println("Please enter your AWS SECRET ACCESS KEY")
		var awsSecretKey string
		fmt.Scanln(&awsSecretKey)

		fmt.Println("Please enter your AWS SSH Key Name")
		var awsAccessSSHKey string
		fmt.Scanln(&awsAccessSSHKey)

		fmt.Println("Please enter your AWS Default Region")
		var awsDefaultRegion string
		fmt.Scanln(&awsDefaultRegion)

		file, err := os.Create("./kubespray/contrib/terraform/aws/credentials.tfvars")
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()

		fmt.Fprintf(file, "AWS_ACCESS_KEY_ID = %s\n", awsAccessKeyID)
		fmt.Fprintf(file, "AWS_SECRET_ACCESS_KEY = %s\n", awsSecretKey)
		fmt.Fprintf(file, "AWS_SSH_KEY_NAME = %s\n", awsAccessSSHKey)
		fmt.Fprintf(file, "AWS_DEFAULT_REGION = %s\n", awsDefaultRegion)
	}
	terrSet := exec.Command("terraform", "destroy", "-var-file=credentials.tfvars", "-force")
	terrSet.Dir = "./kubespray/contrib/terraform/aws/"
	stdout, _ := terrSet.StdoutPipe()
	terrSet.Stderr = terrSet.Stdout
	error := terrSet.Start()
	if error != nil {
		fmt.Println(error)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		//log.Printf(m)
	}

	terrSet.Wait()

	os.Exit(0)
}

// depricated
var ec2IP string

// depricated
func distSelect() (string, string) {
	var sshUser, osLabel string

	centos := map[string]string{
		"user":      "centos",
		"ami_owner": "688023202711",
		"os":        "dcos-centos7",
	}

	ubuntu := map[string]string{
		"user":      "ubuntu",
		"ami_owner": "099720109477",
		"os":        "ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64",
	}

	//Read Configuration File
	viper.SetConfigName("config")

	viper.AddConfigPath(".")
	verr := viper.ReadInConfig() // Find and read the config file
	if verr != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", verr))
	}

	awsAmiID := viper.GetString("aws.ami_id")
	awsInstanceOS := viper.GetString("aws.os")
	sshUser = viper.GetString("aws.ssh_user")

	// Think of a better way to do this
	if awsInstanceOS != "" {
		fmt.Println(awsInstanceOS)
		switch awsInstanceOS {
		case "centos":
			exec.Command("sh", "-c", "sed -i \"\" -e 's/dcos-centos7/"+centos["os"]+"/g' ./kubespray/contrib/terraform/aws/variables.tf").Run()
			exec.Command("sh", "-c", "sed -i \"\" -e 's/688023202711/"+centos["ami_owner"]+"/g' ./kubespray/contrib/terraform/aws/variables.tf").Run()
			sshUser = centos["user"]
			osLabel = "centos"
		case "ubuntu":
			exec.Command("sh", "-c", "sed -i \"\" -e 's#dcos-centos7#"+ubuntu["os"]+"#g' ./kubespray/contrib/terraform/aws/variables.tf").Run()
			exec.Command("sh", "-c", "sed -i \"\" -e 's/688023202711/"+ubuntu["ami_owner"]+"/g' ./kubespray/contrib/terraform/aws/variables.tf").Run()
			sshUser = ubuntu["user"]
			osLabel = "ubuntu"
		// Will only work with 'https://github.com/kubernetes-incubator/kubespray'
		default:
			sshUser = "core"
			osLabel = "coreos"
			return sshUser, osLabel
		}
	} else if awsAmiID != "" && sshUser != "" {
		err := exec.Command("sh", "-c", "sed -i \"\" -e 's/${data.aws_ami.distro.id}/"+awsAmiID+"/g' ./kubespray/contrib/terraform/aws/create-infrastructure.tf").Run()
		if err != nil {
			log.Fatal("Cannot replace AMI ID in Infrastructure template", err)
		}
		osLabel = "Custom-AMI"
	} else if awsAmiID != "" && sshUser == "" {
		log.Fatal("SSH Username is required when using custom AMI")
		return "", ""
	} else {
		log.Fatal("Provide either of AMI ID or OS in the config file.")
		return "", ""
	}

	return sshUser, osLabel
}
