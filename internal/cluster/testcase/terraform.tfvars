aws_cluster_name = "kubernauts"
aws_vpc_cidr_block = "10.250.192.0/18"
aws_cidr_subnets_private = "[&#34;10.250.192.0/20&#34;,&#34;10.250.208.0/20&#34;]"
aws_cidr_subnets_public = "[&#34;10.250.224.0/20&#34;,&#34;10.250.240.0/20&#34;]"

aws_bastion_size = "t2.medium"
aws_kube_master_num = "1"
aws_kube_master_size = "t2.medium"
aws_etcd_num = "1"

aws_etcd_size = "t2.medium"
aws_kube_worker_num = "2"
aws_kube_worker_size = "t2.medium"
aws_elb_api_port = "6443"
k8s_secure_api_port = "6443"
kube_insecure_apiserver_address = ""

default_tags = {
    Env = "devtest"
    Product = "kubernetes"
}