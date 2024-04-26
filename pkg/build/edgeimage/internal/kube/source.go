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
	"go/build"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
)

// ImportPath is the canonical import path for the KubeEdge root package
// this is used by FindSource
const ImportPath = "kubeedge"

func FindSource() (root string, err error) {
	// look up the source the way go build would
	pkg, err := build.Default.Import(ImportPath, build.Default.GOPATH, build.FindOnly|build.IgnoreVendor)
	if err == nil && maybeKubeDir(pkg.Dir) {
		return pkg.Dir, nil
	}
	path, err := findOrCloneKubeEdge(ImportPath)
	if err == nil && path != "" {
		return path, nil
	}
	return "", errors.New("could not find kubeedge source")
}

func maybeKubeDir(dir string) bool {
	// TODO(bentheelder): consider adding other sanity checks
	// check if 'go.mod' exists in the directory
	goModPath := dir + "/go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func findOrCloneKubeEdge(importPath string) (string, error) {
	pkg, err := packages.Load(&packages.Config{Mode: packages.NeedFiles}, importPath)
	if err == nil && len(pkg) > 0 && pkg[0].GoFiles != nil {
		return filepath.Dir(pkg[0].GoFiles[0]), nil
	}

	localDir := filepath.Join(build.Default.GOPATH, "src", importPath)
	fmt.Println("Cloning KubeEdge from GitHub to", localDir)

	if err := gitClone("https://github.com/kubeedge/kubeedge.git", localDir); err != nil {
		return "", err
	}

	return localDir, nil
}

func gitClone(repoURL, localDir string) error {
	cmd := exec.Command("git", "clone", repoURL, localDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone KubeEdge repository: %w", err)
	}
	return nil
}
