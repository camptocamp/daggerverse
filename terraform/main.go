// Terraform
//
// Work with Terraform infrastructure as code tool.

// Copyright Camptocamp SA
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"context"
	"dagger/terraform/internal/dagger"
	"fmt"
	"path"
	"strings"
)

const (
	// Name of Terraform executable binary
	BinaryName = "terraform"

	// Name of Terraform version file
	VersionFileName = ".terraform-version"

	// Name of Terraform plan file
	PlanFileName = "terraform.tfplan"

	// Location of Terraform plugin cache
	PluginCacheDir = "/var/cache/terraform"
)

type Terraform struct {
	// +private
	Version string
}

// Terraform constructor
func New(
	// Terraform version to get
	// +optional
	version string,
) *Terraform {
	terraform := &Terraform{
		Version: version,
	}

	return terraform
}

// Get a Terraform executable binary
func (terraform *Terraform) Binary(
	ctx context.Context,
	// Platform to get Terraform for
	// +optional
	platform dagger.Platform,
) (*dagger.File, error) {
	if terraform.Version == "" {
		return nil, fmt.Errorf("Terraform version must be specified")
	}

	if platform == "" {
		defaultPlatform, err := dag.DefaultPlatform(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get platform: %s", err)
		}

		platform = defaultPlatform
	}

	platformElements := strings.Split(string(platform), "/")

	os := platformElements[0]
	arch := platformElements[1]

	downloadURL := "https://releases.hashicorp.com/terraform/" + terraform.Version

	archiveName := fmt.Sprintf("terraform_%s_%s_%s.zip", terraform.Version, os, arch)
	checksumsName := fmt.Sprintf("terraform_%s_SHA256SUMS", terraform.Version)
	checksumsSignatureName := checksumsName + ".sig"

	binaryName := "terraform"

	if os == "windows" {
		binaryName += ".exe"
	}

	archive := dag.HTTP(downloadURL + "/" + archiveName)
	checksums := dag.HTTP(downloadURL + "/" + checksumsName)
	checksumsSignature := dag.HTTP(downloadURL + "/" + checksumsSignatureName)

	const hashicorpPGPKeyName = "hashicorp.pgp"

	container := dag.Redhat().Container().
		With(dag.Redhat().Packages([]string{
			"gpg",
			"unzip",
		}).Installed).
		WithMountedFile(hashicorpPGPKeyName, dag.CurrentModule().Source().File(hashicorpPGPKeyName)).
		WithExec([]string{"gpg", "--import", hashicorpPGPKeyName}).
		WithMountedFile(archiveName, archive).
		WithMountedFile(checksumsName, checksums).
		WithMountedFile(checksumsSignatureName, checksumsSignature).
		WithExec([]string{"gpg", "--verify", checksumsSignatureName, checksumsName}).
		WithExec([]string{"sh", "-c", "grep -w " + archiveName + " " + checksumsName + " | sha256sum -c"}).
		WithExec([]string{"unzip", archiveName})

	binary := container.File(binaryName)

	return binary, nil
}

// Get a Terraform root filesystem overlay
func (terraform *Terraform) Overlay(
	ctx context.Context,
	// Platform to get Terraform for
	// +optional
	platform dagger.Platform,
	// Filesystem prefix under which to install Terraform
	// +optional
	prefix string,
) (*dagger.Directory, error) {
	if prefix == "" {
		prefix = "/usr/local"
	}

	binary, err := terraform.Binary(ctx, platform)

	if err != nil {
		return nil, fmt.Errorf("failed to get Terraform binary: %s", err)
	}

	overlay := dag.Directory().
		WithDirectory(prefix, dag.Directory().
			WithDirectory("bin", dag.Directory().
				WithFile(BinaryName, binary),
			),
		)

	return overlay, nil
}

// Install Terraform in a container
func (terraform *Terraform) Installation(
	ctx context.Context,
	// Container in which to install Terraform
	container *dagger.Container,
	// Filesystem prefix under which to install Terraform
	// +optional
	prefix string,
) (*dagger.Container, error) {
	platform, err := container.Platform(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get container platform: %s", err)
	}

	overlay, err := terraform.Overlay(ctx, platform, prefix)

	if err != nil {
		return nil, fmt.Errorf("failed to get Terraform overlay: %s", err)
	}

	container = container.
		WithDirectory("/", overlay).
		WithMountedCache(PluginCacheDir, dag.CacheVolume("terraform")).
		WithEnvVariable("TF_PLUGIN_CACHE_DIR", PluginCacheDir)

	return container, nil
}

// Get a Terraform container from a base container
func (terraform *Terraform) Container(
	ctx context.Context,
	// Base container
	container *dagger.Container,
	// Filesystem prefix under which to install Terraform
	// +optional
	prefix string,
) (*dagger.Container, error) {
	container, err := terraform.Installation(ctx, container, prefix)

	if err != nil {
		return nil, fmt.Errorf("failed to install Terraform: %s", err)
	}

	container = container.
		WithEntrypoint([]string{path.Join(prefix, BinaryName)})

	return container, nil
}

// Get a Terraform container from a Red Hat Universal Base Image container
func (terraform *Terraform) RedhatContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Container(dagger.RedhatContainerOpts{Platform: platform}).
		With(dag.Redhat().Packages([]string{
			"git",
			"diffutils",
		}).Installed)

	return terraform.Container(ctx, container, "")
}

// Get a Terraform container from a Red Hat Minimal Universal Base Image container
func (terraform *Terraform) RedhatMinimalContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Minimal().Container(dagger.RedhatMinimalContainerOpts{Platform: platform}).
		With(dag.Redhat().Minimal().Packages([]string{
			"git",
			"diffutils",
		}).Installed)

	return terraform.Container(ctx, container, "")
}

// Get a Terraform container from a Red Hat Micro Universal Base Image container
//
// Some features of Terraform may not be available in a Red Hat Micro Universal Base Image container.
func (terraform *Terraform) RedhatMicroContainer(
	ctx context.Context,
	// Platform to get container for
	// +optional
	platform dagger.Platform,
) (*dagger.Container, error) {
	container := dag.Redhat().Micro().Container(dagger.RedhatMicroContainerOpts{Platform: platform}).
		With(dag.Redhat().CaCertificates)

	return terraform.Container(ctx, container, "")
}

type TerraformWorkspace struct {
	// Get Terraform version
	Version string
	// +private
	Source *dagger.Directory
	// Get a Terraform container with mounted source directory and environment variables set
	Container *dagger.Container
}

// Get a Terraform workspace from a source directory
func (terraform *Terraform) Workspace(
	ctx context.Context,
	// Terraform configuration source directory
	// +optional
	// +defaultPath="/terraform"
	// +ignore=["terraform.tfstate", "terraform.tfstate.backup", "terraform.tfplan"]
	source *dagger.Directory,
) (*TerraformWorkspace, error) {
	if terraform.Version == "" {
		version, err := source.File(VersionFileName).Contents(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to get Terraform version from source: %s", err)
		}

		terraform.Version = strings.TrimSpace(version)
	}

	container, err := terraform.RedhatMinimalContainer(ctx, "")

	if err != nil {
		return nil, fmt.Errorf("failed to get Terraform container: %s", err)
	}

	container = container.
		WithMountedDirectory(".", source)

	workspace := &TerraformWorkspace{
		Version:   terraform.Version,
		Source:    source,
		Container: container,
	}

	return workspace, nil
}

// Set an environment variable in the Terraform workspace
func (workspace *TerraformWorkspace) WithEnvVariable(
	// Environment variable name
	name string,
	// Environment variable value
	value string,
) *TerraformWorkspace {
	workspace.Container = workspace.Container.
		WithEnvVariable(name, value)

	return workspace
}

// Set a secret environment variable in the Terraform workspace
func (workspace *TerraformWorkspace) WithSecretVariable(
	// Environment variable name
	name string,
	// Environment variable secret value
	secret *dagger.Secret,
) *TerraformWorkspace {
	workspace.Container = workspace.Container.
		WithSecretVariable(name, secret)

	return workspace
}

// Initialize the Terraform workspace
func (workspace *TerraformWorkspace) Init(
	// Arguments to pass to Terraform init command
	// +optional
	args ...string,
) *TerraformWorkspace {
	command := append([]string{"init"}, args...)

	workspace.Container = workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true})

	return workspace
}

// Format the Terraform workspace
func (workspace *TerraformWorkspace) Format(
	// Arguments to pass to Terraform fmt command
	// +optional
	args ...string,
) *TerraformWorkspace {
	command := append([]string{"fmt"}, args...)

	workspace.Container = workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true})

	return workspace
}

// Validate the Terraform workspace
func (workspace *TerraformWorkspace) Validate(
	// Arguments to pass to Terraform validate command
	// +optional
	args ...string,
) *TerraformWorkspace {
	command := append([]string{"validate"}, args...)

	workspace.Container = workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true})

	return workspace
}

type TerraformPlan struct {
	// +private
	Workspace *TerraformWorkspace
	// Get the Terraform plan file
	File *dagger.File
}

// Create a Terraform plan for the workspace
func (workspace *TerraformWorkspace) Plan(
	// Arguments to pass to Terraform plan command
	// +optional
	args ...string,
) *TerraformPlan {
	command := append([]string{"plan", "-out=" + PlanFileName}, args...)

	file := workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		File(PlanFileName)

	plan := &TerraformPlan{
		Workspace: workspace,
		File:      file,
	}

	return plan
}

// Show a Terraform plan
func (plan *TerraformPlan) Show(
	ctx context.Context,
	// Arguments to pass to Terraform show command
	// +optional
	args ...string,
) (string, error) {
	return plan.Workspace.Show(ctx, plan.File, args...)
}

// Apply a Terraform plan
func (plan *TerraformPlan) Apply(
	// Arguments to pass to Terraform apply command
	// +optional
	args ...string,
) *TerraformWorkspace {
	return plan.Workspace.Apply(plan.File, args...)
}

// Show a Terraform plan or the Terraform workspace current state
func (workspace *TerraformWorkspace) Show(
	ctx context.Context,
	// Terraform plan file to show (defaults to current state if not specified)
	// +optional
	planFile *dagger.File,
	// Arguments to pass to Terraform show command
	// +optional
	args ...string,
) (string, error) {
	command := append([]string{"show"}, args...)

	if planFile != nil {
		workspace.Container = workspace.Container.
			WithMountedFile(PlanFileName, planFile)

		command = append(command, PlanFileName)
	}

	output, err := workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return output, err
}

// Apply a Terraform plan
func (workspace *TerraformWorkspace) Apply(
	// Terraform plan file to apply
	// +optional
	planFile *dagger.File,
	// Arguments to pass to Terraform apply command
	// +optional
	args ...string,
) *TerraformWorkspace {
	command := append([]string{"apply"}, args...)

	if planFile != nil {
		workspace.Container = workspace.Container.
			WithMountedFile(PlanFileName, planFile)

		command = append(command, PlanFileName)
	}

	workspace.Container = workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true})

	return workspace
}

// Get output variables from a Terraform workspace
func (workspace *TerraformWorkspace) Output(
	ctx context.Context,
	// Output variable (defaults to all variables if empty)
	// +optional
	name string,
	// Arguments to pass to Terraform output command
	// +optional
	args ...string,
) (string, error) {
	command := append([]string{"output"}, args...)

	if name != "" {
		command = append(command, name)
	}

	output, err := workspace.Container.
		WithExec(command, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return output, err
}

// Get Terraform workspace source directory changes
func (workspace *TerraformWorkspace) Changes() *dagger.Changeset {
	return workspace.Container.Directory(".").Changes(workspace.Source)
}

// Get combined buffered standard output and standard error stream of the last executed command in the Terraform workspace container
func (workspace *TerraformWorkspace) CombinedOutput(
	ctx context.Context,
) (string, error) {
	return workspace.Container.CombinedOutput(ctx)
}

// Force evaluation of the Terraform workspace commands
func (workspace *TerraformWorkspace) Sync(
	ctx context.Context,
) (*TerraformWorkspace, error) {
	var err error

	workspace.Container, err = workspace.Container.
		Sync(ctx)

	return workspace, err
}

// Check a Terraform workspace
// +check
func (terraform *Terraform) Check(
	ctx context.Context,
	// +defaultPath="/terraform"
	// +ignore=[".terraform/", "terraform.tfstate", "terraform.tfstate.backup", "terraform.tfplan"]
	source *dagger.Directory,
) error {
	workspace, err := terraform.
		Workspace(ctx, source)

	if err != nil {
		return fmt.Errorf("failed to get Terraform workspace: %s", err)
	}

	_, err = workspace.
		Init("-backend=false").
		Format("-check", "-recursive", "-diff").
		Validate().
		Sync(ctx)

	return err
}
