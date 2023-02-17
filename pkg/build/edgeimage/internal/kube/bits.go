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

	"sigs.k8s.io/kind/pkg/build/nodeimage/shared/kube"
)

type Bits kube.Bits

// shared real bits implementation for now
type bits struct {
	// computed at build time
	binaryPaths []string
	// TODO: maybe we can include cloudcore image also
	// imagePaths  []string
}

var _ Bits = &bits{}

func (b *bits) BinaryPaths() []string {
	return b.binaryPaths
}

func (b *bits) ImagePaths() []string {
	fmt.Println("Unimplemented ImagePaths method")
	return nil
}

func (b *bits) Version() string {
	fmt.Println("Unimplemented Version method")
	return ""
}
