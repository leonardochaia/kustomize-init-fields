// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package main implements an injection function for resource reservations and
// is run with `kustomize fn run -- DIR/`.
package main

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		for _, item := range items {
			if err := initContainers(item); err != nil {
				return nil, err
			}
			if err := initVolumes(item); err != nil {
				return nil, err
			}
		}
		return items, nil
	}
	p := framework.SimpleProcessor{Config: nil, Filter: kio.FilterFunc(fn)}
	cmd := command.Build(p, command.StandaloneDisabled, false)
	command.AddGenerateDockerfile(cmd)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initContainers(r *yaml.RNode) error {
	// lookup the containers field
	containers, err := r.Pipe(yaml.Lookup("spec", "template", "spec", "containers"))
	if err != nil {
		s, _ := r.String()
		return fmt.Errorf("%v: %s", err, s)
	}
	if containers == nil {
		// doesn't have containers, skip the Resource
		return nil
	}

	// visit each container and apply init the fields
	return containers.VisitElements(func(node *yaml.RNode) error {
		err := node.PipeE(
			yaml.Tee(
				yaml.LookupCreate(yaml.SequenceNode, "volumeMounts"),
			),
			yaml.Tee(
				yaml.LookupCreate(yaml.SequenceNode, "envFrom"),
			),
			yaml.Tee(
				yaml.LookupCreate(yaml.SequenceNode, "env"),
			),
		)
		if err != nil {
			s, _ := r.String()
			return fmt.Errorf("%v: %s", err, s)
		}

		return nil
	})
}

func initVolumes(r *yaml.RNode) error {
	// lookup the spec field
	spec, err := r.Pipe(yaml.Lookup("spec", "template", "spec"))
	if err != nil {
		s, _ := r.String()
		return fmt.Errorf("%v: %s", err, s)
	}
	if spec == nil {
		// doesn't have spec, skip the Resource
		return nil
	}

	return spec.PipeE(
		yaml.LookupCreate(yaml.SequenceNode, "volumes"),
	)
}
