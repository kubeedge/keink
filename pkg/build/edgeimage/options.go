package edgeimage

import (
	"sigs.k8s.io/kind/pkg/log"
)

// Option is a configuration option supplied to Build
type Option interface {
	apply(*buildContext) error
}

type optionAdapter func(*buildContext) error

func (c optionAdapter) apply(o *buildContext) error {
	return c(o)
}

// WithImage configures a build to tag the built image with `image`
func WithImage(image string) Option {
	return optionAdapter(func(b *buildContext) error {
		b.image = image
		return nil
	})
}

// WithBaseImage configures a build to use `image` as the base image
func WithBaseImage(image string) Option {
	return optionAdapter(func(b *buildContext) error {
		b.baseImage = image
		return nil
	})
}

// WithKubeEdgeRoot sets the path to the KubeEdge source directory (if empty, the path will be autodetected)
func WithKubeEdgeRoot(root string) Option {
	return optionAdapter(func(b *buildContext) error {
		b.kubeEdgeRoot = root
		return nil
	})
}

// WithLogger sets the logger
func WithLogger(logger log.Logger) Option {
	return optionAdapter(func(b *buildContext) error {
		b.logger = logger
		return nil
	})
}

// WithArch sets the architecture to build for
func WithArch(arch string) Option {
	return optionAdapter(func(b *buildContext) error {
		if arch != "" {
			b.arch = arch
		}
		return nil
	})
}
