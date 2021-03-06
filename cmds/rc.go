package cmds

import (
	"github.com/codegangsta/cli"

	"github.com/fabric8io/kansible/ansible"
	"github.com/fabric8io/kansible/log"

	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// RC creates or updates the kansible ReplicationController for some hosts in an Ansible inventory
func RC(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		log.Die("Expected argument [hosts] for the name of the hosts in the ansible inventory file")
	}
	hosts := args[0]

	f := cmdutil.NewFactory(nil)
	if f == nil {
		log.Die("Failed to create Kuberentes client factory!")
	}
	kubeclient, _ := f.Client()
	if kubeclient == nil {
		log.Die("Failed to create Kuberentes client!")
	}
	ns, _, _ := f.DefaultNamespace()
	if len(ns) == 0 {
		ns = "default"
	}

	inventory, err := osExpandAndVerify(c, "inventory")
	if err != nil {
		fail(err)
	}

	hostEntries, err := ansible.LoadHostEntries(inventory, hosts)
	if err != nil {
		fail(err)
	}
	log.Info("Found %d host entries in the Ansible inventory for %s", len(hostEntries), hosts)

	rcFile := "kubernetes/" + hosts + "/rc.yml"

	_, err = ansible.UpdateKansibleRC(hostEntries, hosts, kubeclient, ns, rcFile)
	if err != nil {
		fail(err)
	}
}

