{{define "ClusterConfig"}}
aws_cluster_name = {{.AwsClusterName}}
aws_vpc_cidr_block = {{.AwsVpcCidrBlock}}
aws_cidr_subnets_private = {{.AwsCidrSubnetsPrivat}}
aws_cidr_subnets_public = {{.AwsCidrSubnetsPubli}}

aws_bastion_size = {{.AwsBastionSize}}
aws_kube_master_num = {{.AwsKubeMasterNu}}
aws_kube_master_size = {{.AwsKubeMasterSize}}
aws_etcd_num = {{.AwsEtcdNu}}

aws_etcd_size = {{.AwsEtcdSize}}
aws_kube_worker_num = {{.AwsKubeWorkerNu}}
aws_kube_worker_size = {{.AwsKubeWorkerSize}}
aws_elb_api_port = {{.AwsElbAPIPor}}
k8s_secure_api_port = {{.K8sSecureAPIPor}}
kube_insecure_apiserver_address = {{.KubeInsecureApiserverAddress}}

default_tags = {
    Env = 'devtest',
    Product = 'kubernetes'
}
{{end}}