package edgeimage

import (
	"fmt"
	"math/rand"
	"path"
	"strings"
	"time"

	"sigs.k8s.io/kind/pkg/build/nodeimage/shared/container/docker"
	"sigs.k8s.io/kind/pkg/build/nodeimage/shared/kube"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/log"
)

// buildContext is used to build the keink node image, and contains
// build configuration
type buildContext struct {
	// option fields
	image        string
	baseImage    string
	logger       log.Logger
	arch         string
	kubeEdgeRoot string
	// non-option fields
	builder kube.Builder
}

// Build builds the cluster node image, the sourcedir must be set on
// the buildContext
func (c *buildContext) Build() (err error) {
	// ensure kubernetes build is up to date first
	c.logger.V(0).Info("Starting to build KubeEdge")
	bits, err := c.builder.Build()
	if err != nil {
		c.logger.Errorf("Failed to build KubeEdge: %v", err)
		return errors.Wrap(err, "failed to build KubeEdge")
	}
	c.logger.V(0).Info("Finished building KubeEdge")

	// then the perform the actual docker image build
	c.logger.V(0).Info("Building edge image ...")
	return c.buildImage(bits)
}

func (c *buildContext) buildImage(bits kube.Bits) error {
	// create build container
	// NOTE: we are using docker run + docker commit so we can install
	// debian packages without permanently copying them into the image.
	// if docker gets proper squash support, we can rm them instead
	// This also allows the KubeBit implementations to perform programmatic
	// install in the image
	containerID, err := c.createBuildContainer()
	cmder := docker.ContainerCmder(containerID)

	// ensure we will delete it
	if containerID != "" {
		defer func() {
			_ = exec.Command("docker", "rm", "-f", "-v", containerID).Run()
		}()
	}
	if err != nil {
		c.logger.Errorf("Image build Failed! Failed to create build container: %v", err)
		return err
	}

	c.logger.V(0).Info("Building in " + containerID)

	// helper we will use to run "build steps"
	execInBuild := func(command string, args ...string) error {
		return exec.InheritOutput(cmder.Command(command, args...)).Run()
	}

	// make artifacts directory
	if err = execInBuild("mkdir", "/kubeedge/"); err != nil {
		c.logger.Errorf("Image build Failed! Failed to make directory %v", err)
		return err
	}

	// make kubeedge directory
	if err = execInBuild("mkdir", "-p", "/etc/kubeedge/config/"); err != nil {
		c.logger.Errorf("Image build Failed! Failed to make directory /etc/kubeedge/config/ %v", err)
		return err
	}
	if err = execInBuild("mkdir", "-p", "/etc/kubeedge/crds/"); err != nil {
		c.logger.Errorf("Image build Failed! Failed to make directory /etc/kubeedge/crds/ %v", err)
		return err
	}

	// copy artifacts in
	for _, binary := range bits.BinaryPaths() {

		// kubeedge binaries should be /usr/local/bin, service file expects /usr/local/bin/edgecore and /usr/local/bin/cloudcore
		nodePath := "/usr/local/bin/" + path.Base(binary)
		if strings.Contains(path.Base(binary), ".service") {
			// service file should be put /etc/systemd/system/
			nodePath = "/etc/systemd/system/" + path.Base(binary)
		} else if strings.Contains(path.Base(binary), ".yaml") {
			// CRDs yaml should be /etc/kubeedge/crds/
			nodePath = "/etc/kubeedge/crds/" + path.Base(binary)
		}

		if err := exec.Command("docker", "cp", binary, containerID+":"+nodePath).Run(); err != nil {
			return err
		}
		if err := execInBuild("chmod", "+x", nodePath); err != nil {
			return err
		}
		if err := execInBuild("chown", "root:root", nodePath); err != nil {
			return err
		}
	}

	// Save the image changes to a new image
	cmd := exec.Command(
		"docker", "commit",
		// we need to put this back after changing it when running the image
		"--change", `ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]`,
		containerID, c.image,
	)
	exec.InheritOutput(cmd)
	if err = cmd.Run(); err != nil {
		c.logger.Errorf("Image build Failed! Failed to save image: %v", err)
		return err
	}

	c.logger.V(0).Infof("Image %q build completed.", c.image)
	return nil
}

func (c *buildContext) createBuildContainer() (id string, err error) {
	// attempt to explicitly pull the image if it doesn't exist locally
	// we don't care if this errors, we'll still try to run which also pulls
	_ = docker.Pull(c.logger, c.baseImage, dockerBuildOsAndArch(c.arch), 4)
	// this should be good enough: a specific prefix, the current unix time,
	// and a little random bits in case we have multiple builds simultaneously
	random := rand.New(rand.NewSource(time.Now().UnixNano())).Int31()
	id = fmt.Sprintf("kind-build-%d-%d", time.Now().UTC().Unix(), random)
	err = docker.Run(
		c.baseImage,
		[]string{
			"-d", // make the client exit while the container continues to run
			// the container should hang forever so we can exec in it
			"--entrypoint=sleep",
			"--name=" + id,
			"--platform=" + dockerBuildOsAndArch(c.arch),
		},
		[]string{
			"infinity", // sleep infinitely to keep the container around
		},
	)
	if err != nil {
		return id, errors.Wrap(err, "failed to create build container")
	}
	return id, nil
}

func dockerBuildOsAndArch(arch string) string {
	return "linux/" + arch
}
