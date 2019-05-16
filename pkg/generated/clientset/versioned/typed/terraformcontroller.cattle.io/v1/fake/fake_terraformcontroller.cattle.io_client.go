// Code generated by main. DO NOT EDIT.

package fake

import (
	v1 "github.com/rancher/terraform-controller/pkg/generated/clientset/versioned/typed/terraformcontroller.cattle.io/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeTerraformcontrollerV1 struct {
	*testing.Fake
}

func (c *FakeTerraformcontrollerV1) Executions(namespace string) v1.ExecutionInterface {
	return &FakeExecutions{c, namespace}
}

func (c *FakeTerraformcontrollerV1) ExecutionRuns(namespace string) v1.ExecutionRunInterface {
	return &FakeExecutionRuns{c, namespace}
}

func (c *FakeTerraformcontrollerV1) Modules(namespace string) v1.ModuleInterface {
	return &FakeModules{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeTerraformcontrollerV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
