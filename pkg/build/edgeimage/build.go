package edgeimage

import (
	"runtime"

	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/log"

	internalkube "github.com/kubeedge/keink/pkg/build/edgeimage/internal/kube"
)

// Build builds a node image using the supplied options
func Build(options ...Option) error {
	// default options
	ctx := &buildContext{
		image:     DefaultImage,
		baseImage: DefaultBaseImage,
		logger:    log.NoopLogger{},
		arch:      runtime.GOARCH,
	}

	// apply user options
	for _, option := range options {
		if err := option.apply(ctx); err != nil {
			return err
		}
	}

	// verify that we're using a supported arch
	if !supportedArch(ctx.arch) {
		ctx.logger.Warnf("unsupported architecture %q", ctx.arch)
	}

	// locate sources if no KubeEdge source was specified
	if ctx.kubeEdgeRoot == "" {
		kubeEdgeRoot, err := internalkube.FindSource()
		if err != nil {
			return errors.Wrap(err, "error finding kubeedge root")
		}
		ctx.kubeEdgeRoot = kubeEdgeRoot
	}

	// initialize bits
	builder, err := internalkube.NewDockerBuilder(ctx.logger, ctx.kubeEdgeRoot, ctx.arch)
	if err != nil {
		return err
	}
	ctx.builder = builder

	// do the actual build
	return ctx.Build()
}

func supportedArch(arch string) bool {
	switch arch {
	default:
		return false
	// currently we nominally support building node images for these
	case "amd64":
	case "arm64":
	case "ppc64le":
	}
	return true
}
