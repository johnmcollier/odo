package utils

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	adaptersCommon "github.com/openshift/odo/pkg/devfile/adapters/common"
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
	"github.com/openshift/odo/pkg/lclient"
	"github.com/openshift/odo/pkg/log"
	"github.com/pkg/errors"
)

// ComponentExists checks if Docker containers labeled with the specified component name exists
func ComponentExists(client lclient.Client, name string) bool {
	containers := GetComponentContainers(client, name)
	return len(containers) != 0
}

// GetComponentContainers returns a list of the running component containers
func GetComponentContainers(client lclient.Client, name string) (containers []types.Container) {
	containerList, err := client.GetContainerList()
	if err != nil {
		return
	}
	containers = client.GetContainersByComponent(name, containerList)

	return
}

// ConvertEnvs converts environment variables from the devfile structure to an array of strings, as expected by Docker
func ConvertEnvs(vars []common.DockerimageEnv) []string {
	dockerVars := []string{}
	for _, env := range vars {
		envString := fmt.Sprintf("%s=%s", *env.Name, *env.Value)
		dockerVars = append(dockerVars, envString)
	}
	return dockerVars
}

// DoesContainerNeedUpdating returns true if a given container needs to be removed and recreated
// This function compares values in the container vs corresponding values in the devfile component.
// If any of the values between the two differ, a restart is required (and this function returns true)
// Unlike Kube, Docker doesn't provide a mechanism to update a container in place only when necesary
// so this function is necessary to prevent having to restart the container on every odo pushs
func DoesContainerNeedUpdating(component common.DevfileComponent, containerConfig *container.Config, devfileMounts []mount.Mount, containerMounts []types.MountPoint) bool {
	// If the image was changed in the devfile, the container needs to be updated
	if *component.Image != containerConfig.Image {
		return true
	}

	// Update the container if the volumes were updated in the devfile
	for _, devfileMount := range devfileMounts {
		if !containerHasMount(devfileMount, containerMounts) {
			return true
		}
	}

	// Update the container if the env vars were updated in the devfile
	// Need to convert the devfile envvars to the format expected by Docker
	devfileEnvVars := ConvertEnvs(component.Env)
	for _, envVar := range devfileEnvVars {
		if !containerHasEnvVar(envVar, containerConfig.Env) {
			return true
		}
	}
	return false

}

// AddVolumeToComp adds the project volume to the container host config
func AddVolumeToComp(volumeName, volumeMount string, hostConfig *container.HostConfig) *container.HostConfig {
	mount := mount.Mount{
		Type:   mount.TypeVolume,
		Source: volumeName,
		Target: volumeMount,
	}
	hostConfig.Mounts = append(hostConfig.Mounts, mount)

	return hostConfig
}

// GetProjectVolumeLabels returns the label selectors used to retrieve the project/source volume for a given component
func GetProjectVolumeLabels(componentName string) map[string]string {
	volumeLabels := map[string]string{
		"component": componentName,
		"type":      "projects",
	}
	return volumeLabels
}

// containerHasEnvVar returns true if the specified env var (and value) exist in the list of container env vars
func containerHasEnvVar(envVar string, containerEnv []string) bool {
	for _, env := range containerEnv {
		if envVar == env {
			return true
		}
	}
	return false
}

// containerHasMount returns true if the specified volume is mounted in the given container
func containerHasMount(devfileMount mount.Mount, containerMounts []types.MountPoint) bool {
	for _, mount := range containerMounts {
		if devfileMount.Source == mount.Name && devfileMount.Target == mount.Destination {
			return true
		}
	}
	return false
}

// GetSupervisordVolumeLabels returns the label selectors used to retrieve the supervisord volume
func GetSupervisordVolumeLabels() map[string]string {
	supervisordLabels := map[string]string{
		"name": adaptersCommon.SupervisordVolumeName,
		"type": "supervisord",
	}
	return supervisordLabels
}

// CreateAndInitSupervisordVolume creates the supervisord volume and initializes it with odo init
func CreateAndInitSupervisordVolume(client lclient.Client) (string, error) {
	supervisordLabels := GetSupervisordVolumeLabels()
	supervisordVolume, err := client.CreateVolume("", supervisordLabels)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to create supervisord volume for component")
	}
	supervisordVolumeName := supervisordVolume.Name

	err = StartBootstrapSupervisordInitContainer(client, supervisordVolumeName)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to start supervisord container for component")
	}

	return supervisordVolumeName, nil
}

// StartBootstrapSupervisordInitContainer ...
func StartBootstrapSupervisordInitContainer(client lclient.Client, supervisordVolumeName string) error {
	supervisordLabels := GetSupervisordVolumeLabels()

	image := adaptersCommon.GetBootstrapperImage()
	command := []string{"/usr/bin/cp"}
	args := []string{
		"-r",
		adaptersCommon.OdoInitImageContents,
		adaptersCommon.SupervisordMountPath,
	}

	s := log.Spinner("Pulling image " + image)
	err := client.PullImage(image)
	if err != nil {
		s.End(false)
		return errors.Wrapf(err, "Unable to pull %s image", image)
	}
	s.End(true)

	containerConfig := client.GenerateContainerConfig(image, command, args, nil, supervisordLabels)
	hostConfig := container.HostConfig{}

	AddVolumeToComp(supervisordVolumeName, adaptersCommon.SupervisordMountPath, &hostConfig)

	// Create the docker container
	s = log.Spinner("Starting container for " + image)
	defer s.End(false)
	err = client.StartContainer(&containerConfig, &hostConfig, nil)
	if err != nil {
		return err
	}
	s.End(true)

	return nil
}
