/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kube

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/kind/pkg/build/nodeimage/shared/kube"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/log"
)

// TODO(bentheelder): plumb through arch

// dockerBuilder implements Bits for a local docker-ized make / bash build
type dockerBuilder struct {
	kubeEdgeRoot string
	arch         string
	logger       log.Logger
}

var _ Builder = &dockerBuilder{}

// NewDockerBuilder returns a new Bits backed by the docker-ized build,
// given kubeRoot, the path to the KubeEdge source directory
func NewDockerBuilder(logger log.Logger, kubeEdgeRoot, arch string) (Builder, error) {
	return &dockerBuilder{
		kubeEdgeRoot: kubeEdgeRoot,
		arch:         arch,
		logger:       logger,
	}, nil
}

// Build implements Bits.Build
func (b *dockerBuilder) Build() (kube.Bits, error) {
	// cd to KubeEdge source
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// make sure we cd back when done
	defer func() {
		// TODO(bentheelder): set return error?
		_ = os.Chdir(cwd)
	}()
	if err := os.Chdir(b.kubeEdgeRoot); err != nil {
		return nil, err
	}

	// we will pass through the environment variables, prepending defaults
	// NOTE: if env are specified multiple times the last one wins
	// NOTE: currently there are no defaults so this is essentially a deep copy
	env := append([]string{}, os.Environ()...)
	// compile kubeedge we cannot use the container build mode
	// or it will report error: the input device is not a TTY
	// we will compile kubeedge with local machine environment
	env = append(env, "BUILD_WITH_CONTAINER=false")

	// binaries we want to build
	what := []string{
		// binaries we use directly
		"keadm",
		"cloudcore",
		"edgecore",
	}

	// build binaries
	for _, component := range what {
		cmd := exec.Command("make",
			"all",
			"WHAT="+component,
		).SetEnv(env...)
		exec.InheritOutput(cmd)
		if err := cmd.Run(); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to build %s", component))
		}
	}

	binDir := filepath.Join(b.kubeEdgeRoot,
		"_output", "local", "bin",
	)
	crdDir := filepath.Join(b.kubeEdgeRoot,
		"build", "crds",
	)

	// use edgecore.service and cloudcore.service file under this repo
	serviceDir := filepath.Join(cwd,
		"build", "tools",
	)

	return &bits{
		binaryPaths: []string{
			// binaries for kubeedge
			filepath.Join(filepath.Join(binDir, "keadm")),
			filepath.Join(filepath.Join(binDir, "cloudcore")),
			filepath.Join(filepath.Join(binDir, "edgecore")),

			// CRDs required by KubeEdge
			filepath.Join(crdDir, "devices", "devices_v1alpha2_device.yaml"),
			filepath.Join(crdDir, "devices", "devices_v1alpha2_devicemodel.yaml"),
			filepath.Join(crdDir, "reliablesyncs", "cluster_objectsync_v1alpha1.yaml"),
			filepath.Join(crdDir, "reliablesyncs", "objectsync_v1alpha1.yaml"),
			filepath.Join(crdDir, "router", "router_v1_rule.yaml"),
			filepath.Join(crdDir, "router", "router_v1_ruleEndpoint.yaml"),

			// cloudcore.service and edgecore.service
			filepath.Join(serviceDir, "edgecore.service"),
			filepath.Join(serviceDir, "cloudcore.service"),
		},
	}, nil
}
