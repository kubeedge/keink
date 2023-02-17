package edgeimage

import "sigs.k8s.io/kind/pkg/apis/config/defaults"

// DefaultImage is the default name:tag for the built image
const DefaultImage = "kubeedge/node:latest"

// DefaultBaseImage is the default base image used
// we will add KubeEdge components based on this image
// Just keep the same with kind default node image
const DefaultBaseImage = defaults.Image
