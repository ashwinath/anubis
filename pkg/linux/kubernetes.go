package linux

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ashwinath/anubis/pkg/config"
	"github.com/ashwinath/anubis/pkg/logger"
	"github.com/ashwinath/anubis/pkg/utils"
)

const kubernetesRepoLocation = "/etc/yum.repos.d/kubernetes.repo"
const kubernetesRepoContentsFormat = `[kubernetes]
name=Kubernetes
baseurl=https://pkgs.k8s.io/core:/stable:/%s/rpm/
enabled=1
gpgcheck=1
gpgkey=https://pkgs.k8s.io/core:/stable:/%s/rpm/repodata/repomd.xml.key
exclude=kubelet kubeadm kubectl cri-tools kubernetes-cni
`

var kubeletVersionRegex = regexp.MustCompile(`^Kubernetes (?P<version>v\d+\.\d+).*\n$`)

func Kubernetes(c config.KubernetesConfig) error {
	if version, ok := getKubernetesVersion(); ok {
		if version != c.Version {
			// TODO: change kubernetes version
		}
	} else {
		return installFreshKubernetes(c)
	}
	return nil
}

func installFreshKubernetes(c config.KubernetesConfig) error {
	logger.Infof("installing fresh set of kubernetes")

	// Disable swap
	if out, err := exec.Command("dnf", "-y", "remove", "zram-generator-defaults").CombinedOutput(); err != nil {
		return fmt.Errorf("could not remove zram-generator-defaults, output: %s, error: %s", out, err)
	}

	if out, err := exec.Command("swapoff", "-a").CombinedOutput(); err != nil {
		return fmt.Errorf("could not disable swap, output: %s, error: %s", out, err)
	}

	// Disable SELinux
	if out, err := exec.Command("setenforce", "0").CombinedOutput(); err != nil {
		return fmt.Errorf("could not setenforce 0, output: %s, error: %s", out, err)
	}

	if out, err := exec.Command("sed", "-i", "s/^SELINUX=enforcing$/SELINUX=permissive/", "/etc/selinux/config").CombinedOutput(); err != nil {
		return fmt.Errorf("could replace SELINUX enforcing to permissive, output: %s, error: %s", out, err)
	}

	// Install kubelet and kubeadm
	if _, err := os.Stat(kubernetesRepoLocation); err != nil {
		f, err := os.Create(kubernetesRepoLocation)
		if err != nil {
			return fmt.Errorf("could not create file at %s, error: %s", kubernetesRepoLocation, err)
		}

		f.Close()
	}

	if err := os.WriteFile(kubernetesRepoLocation, []byte(fmt.Sprintf(kubernetesRepoContentsFormat, c.Version, c.Version)), 0644); err != nil {
		return fmt.Errorf("Could not write to kubernetes.repo at %s", kubernetesRepoLocation)
	}

	if out, err := exec.Command("dnf", "-y", "install", "kubelet", "kubeadm", "--disableexcludes=kubernetes").CombinedOutput(); err != nil {
		return fmt.Errorf("could not install kubelet and kubeadm, output: %s, error: %s", out, err)
	}

	if out, err := exec.Command("systemctl", "enable", "--now", "kubelet").CombinedOutput(); err != nil {
		return fmt.Errorf("could not enable kubelet, output: %s, error: %s", out, err)
	}

	if out, err := exec.Command("systemctl", "enable", "--now", "crio").CombinedOutput(); err != nil {
		return fmt.Errorf("could not enable kubelet, output: %s, error: %s", out, err)
	}

	// Allow firewall settings
	if out, err := exec.Command("firewall-cmd", "--add-port=6443/tcp", "--add-port=10250/tcp", "--permanent").CombinedOutput(); err != nil {
		return fmt.Errorf("could not enable kubelet, output: %s, error: %s", out, err)
	}

	if out, err := exec.Command("firewall-cmd", "--reload").CombinedOutput(); err != nil {
		return fmt.Errorf("could not enable kubelet, output: %s, error: %s", out, err)
	}

	// Enable networking stuff
	if out, err := exec.Command("firewall-cmd", "--reload").CombinedOutput(); err != nil {
		return fmt.Errorf("could not enable kubelet, output: %s, error: %s", out, err)
	}

	if _, err := os.Stat("/etc/modules"); err != nil {
		f, err := os.Create("/etc/modules")
		if err != nil {
			return fmt.Errorf("could not create file at %s, error: %s", kubernetesRepoLocation, err)
		}

		f.Close()

		if err := os.WriteFile("/etc/modules", []byte("br_netfilter"), 0644); err != nil {
			return fmt.Errorf("Could not write to /etc/modules, error: %s", err)
		}
	}

	if out, err := exec.Command("modprobe", "br_netfilter").CombinedOutput(); err != nil {
		return fmt.Errorf("could not configure modprobe br_netfilter, output: %s, error: %s", out, err)
	}

	if out, err := exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").CombinedOutput(); err != nil {
		return fmt.Errorf("could not configure net.ipv4.ip_forward=1, output: %s, error: %s", out, err)
	}

	if c.IsMaster {
		logger.Infof("is master, configuring master")
		if out, err := exec.Command("kubeadm", "init", "--cri-socket", "/var/run/crio/crio.sock", "--pod-network-cidr", "10.244.0.0/16").CombinedOutput(); err != nil {
			return fmt.Errorf("could not kubeadmin init, output: %s, error: %s", out, err)
		}
		// Kubeconfig
		if _, err := os.Stat("/home/ashwin/.kube"); err != nil {
			if err := os.MkdirAll("/home/ashwin/.kube", 0644); err != nil {
				return fmt.Errorf("could not make .kube folder, error: %s", err)
			}

			if err := os.Chown("/home/ashwin/.kube", 1000, 1000); err != nil {
				return fmt.Errorf("could not chown .kube folder, error: %s", err)
			}
		}

		if err := utils.CopyFile("/etc/kubernetes/admin.conf", "/home/ashwin/.kube/config"); err != nil {
			return fmt.Errorf("could not copy kubeconfig to home folder, error: %s", err)
		}

		if err := os.Chown("/home/ashwin/.kube/config", 1000, 1000); err != nil {
			return fmt.Errorf("could not chown .kube/config folder, error: %s", err)
		}

		if err := os.Chmod("/home/ashwin/.kube/config", 0600); err != nil {
			return fmt.Errorf("could not chmod .kube/config folder, error: %s", err)
		}

		// Install flannel
		for i := 0; i < 100; i++ {
			logger.Infof("waiting for kube api server to be ready to apply flannel manifests")
			out, _ := exec.Command("bash", "-c", "kubectl --kubeconfig /etc/kubernetes/admin.conf get po -n kube-system -l component=kube-apiserver -o yaml | yq '.items[0].status.containerStatuses[0].ready'").Output()
			logger.Infof("kube-apiserver status: %s", string(out))
			if strings.Contains(string(out), "true") {
				logger.Infof("kube api server is ready")
				break
			}
			time.Sleep(10 * time.Second)
		}

		if out, err := exec.Command("kubectl", "--kubeconfig", "/etc/kubernetes/admin.conf", "apply", "-f", "https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml").CombinedOutput(); err != nil {
			return fmt.Errorf("could not apply flannel manifest, output: %s, error: %s", out, err)
		}

		// taint node
		if out, err := exec.Command("kubectl", "--kubeconfig", "/etc/kubernetes/admin.conf", "taint", "nodes", "--all", "node-role.kubernetes.io/control-plane-").CombinedOutput(); err != nil {
			return fmt.Errorf("could not taint nodes, output: %s, error: %s", out, err)
		}

		// Not sure why this is even needed.
		resetCNI()
	} else {
		logger.Infof("joining kubernetes cluster at %s", c.MasterIP)

		token := os.Getenv("KUBEADM_JOIN_TOKEN")
		if token == "" {
			return fmt.Errorf("KUBEADM_JOIN_TOKEN is not set")
		}

		hash := os.Getenv("KUBEADM_JOIN_HASH")
		if hash == "" {
			return fmt.Errorf("KUBEADM_JOIN_TOKEN is not set")
		}

		if out, err := exec.Command("kubeadm", "join", c.MasterIP, "--token", token, "--discovery-token-ca-cert-hash", hash).CombinedOutput(); err != nil {
			return fmt.Errorf("could not taint nodes, output: %s, error: %s", out, err)
		}
	}

	logger.Infof("done installing fresh set of kubernetes")
	return nil
}

func getKubernetesVersion() (string, bool) {
	out, _ := exec.Command("kubelet", "--version").CombinedOutput()
	m := utils.FindAllGroups(kubeletVersionRegex, string(out))
	if v, ok := m["version"]; ok {
		return v, true
	}

	return "", false
}

func resetCNI() {
	_, _ = exec.Command("ip", "link", "set", "cni0", "down").CombinedOutput()
	_, _ = exec.Command("ip", "link", "set", "flannel.1", "down").CombinedOutput()
	_, _ = exec.Command("ip", "link", "delete", "cni0").CombinedOutput()
	_, _ = exec.Command("ip", "link", "delete", "flannel.1").CombinedOutput()
	_, _ = exec.Command("systemctl", "restart", "crio").CombinedOutput()
	_, _ = exec.Command("systemctl", "restart", "kubelet").CombinedOutput()
}
