package component

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/openshift/odo/pkg/devfile/adapters/common"
	adaptersCommon "github.com/openshift/odo/pkg/devfile/adapters/common"
	"github.com/openshift/odo/pkg/devfile/adapters/docker/utils"
	versionsCommon "github.com/openshift/odo/pkg/devfile/versions/common"
	"github.com/openshift/odo/pkg/lclient"
	"github.com/openshift/odo/pkg/log"
)

// New instantiantes a component adapter
func New(adapterContext common.AdapterContext, client lclient.Client) Adapter {
	return Adapter{
		Client:         client,
		AdapterContext: adapterContext,
	}
}

// Adapter is a component adapter implementation for Kubernetes
type Adapter struct {
	Client lclient.Client
	common.AdapterContext
}

// Push updates the component if a matching component exists or creates one if it doesn't exist
func (a Adapter) Push(parameters common.PushParameters) (err error) {
	componentExists := utils.ComponentExists(a.Client, a.ComponentName)

	if componentExists {
		err = a.updateComponent()
	} else {
		err = a.createComponent()
	}

	if err != nil {
		return errors.Wrap(err, "unable to create or update component")
	}

	return nil
}

// DoesComponentExist returns true if a component with the specified name exists, false otherwise
func (a Adapter) DoesComponentExist(cmpName string) bool {
	return utils.ComponentExists(a.Client, cmpName)
}

func (a Adapter) createComponent() (err error) {
	componentName := a.ComponentName

	volumeLabels := map[string]string{
		"component": componentName,
		"type":      "projects",
	}

	supportedComponents := adaptersCommon.GetSupportedComponents(a.Devfile.Data)
	if len(supportedComponents) == 0 {
		return fmt.Errorf("No valid components found in the devfile")
	}

	// Create a docker volume to store the project source code
	volume, err := a.Client.CreateVolume(volumeLabels)
	if err != nil {
		return errors.Wrapf(err, "Unable to create project source volume for component %s", componentName)
	}

	projectVolumeName := volume.Name

	// Loop over each component and start a container for it
	for _, comp := range supportedComponents {
		s := log.Spinner("Pulling image " + *comp.Image)
		defer s.End(false)

		// Pull the image (as the Docker daemon requires it to be on the system before starting it)
		err = a.Client.PullImage(*comp.Image)
		if err != nil {
			return errors.Wrapf(err, "Unable to pull %s image", *comp.Image)
		}
		s.End(true)

		// Start the container
		err = a.startContainer(componentName, projectVolumeName, comp)
		if err != nil {
			return errors.Wrapf(err, "Unable to start container for devfile component %s", *comp.Alias)
		}
		glog.V(3).Infof("Successfully created container %s for component %s", *comp.Image, componentName)
	}
	glog.V(3).Infof("Successfully created all containers for component %s", componentName)

	return nil
}

func (a Adapter) updateComponent() (err error) {
	glog.V(3).Info("The component already exists, attempting to update it")
	componentName := a.ComponentName

	// Get the project source volume
	volumeLabels := map[string]string{
		"component": componentName,
		"type":      "projects",
	}
	vols, err := a.Client.GetVolumesByLabel(volumeLabels)
	if err != nil {
		return errors.Wrapf(err, "Unable to retrieve source volume for component "+componentName)
	}
	if len(vols) == 0 {
		return fmt.Errorf("Unable to find source volume for component %s", componentName)
	}
	projectVolumeName := vols[0].Name

	supportedComponents := adaptersCommon.GetSupportedComponents(a.Devfile.Data)
	if len(supportedComponents) == 0 {
		return fmt.Errorf("No valid components found in the devfile")
	}

	for _, comp := range supportedComponents {
		// Check to see if this component is already running and if so, update it
		// If component isn't running, re-create it, as it either may be new, or crashed.
		containers, err := a.Client.GetContainersByComponentAndAlias(componentName, *comp.Alias)
		if err != nil {
			return errors.Wrapf(err, "unable to list containers for component %s", componentName)
		}
		if len(containers) == 0 {
			// Container doesn't exist, so need to pull its image (to be safe) and start a new container
			s := log.Spinner("Pulling image " + *comp.Image)
			defer s.End(false)

			err = a.Client.PullImage(*comp.Image)
			if err != nil {
				return errors.Wrapf(err, "Unable to pull %s image", *comp.Image)
			}
			s.End(true)

			// Start the container
			err = a.startContainer(componentName, projectVolumeName, comp)
			if err != nil {
				return errors.Wrapf(err, "Unable to start container for devfile component %s", *comp.Alias)
			}
			glog.V(3).Infof("Successfully created container %s for component %s", *comp.Image, componentName)
		} else {
			// Container already exists
			containerID := containers[0].ID

			// Get the associated container config from the container ID
			containerConfig, err := a.Client.GetContainerConfig(containerID)
			if err != nil {
				return errors.Wrapf(err, "unable to get the container config for component %s", componentName)
			}

			// See if the container needs to be updated
			if utils.DoesContainerNeedUpdating(comp, containerConfig) {
				glog.V(3).Info("Updating the component")
				// Remove the container
				err := a.Client.RemoveContainer(containerID)
				if err != nil {
					return errors.Wrapf(err, "Unable to remove container %s for component %s", containerID, *comp.Alias)
				}

				// Start the container
				err = a.startContainer(componentName, projectVolumeName, comp)
				if err != nil {
					return errors.Wrapf(err, "Unable to start container for devfile component %s", *comp.Alias)
				}
				glog.V(3).Infof("Successfully created container %s for component %s", *comp.Image, componentName)
			}
		}
	}
	return nil
}

func (a Adapter) startContainer(componentName string, projectVolumeName string, comp versionsCommon.DevfileComponent) error {
	containerConfig := a.getContainerConfig(componentName, comp)
	hostConfig := container.HostConfig{}

	// If the component set `mountSources` to true, add the source volume to it
	if comp.MountSources {
		addProjectVolumeToComp(projectVolumeName, &hostConfig)
	}

	// Create the docker container
	s := log.Spinner("Starting container for " + *comp.Image)
	defer s.End(false)
	err := a.Client.StartContainer(&containerConfig, &hostConfig, nil)
	if err != nil {
		return err
	}
	s.End(true)
	return nil
}
func (a Adapter) getContainerConfig(componentName string, comp versionsCommon.DevfileComponent) container.Config {
	// Convert the env vars in the Devfile to the format expected by Docker
	envVars := utils.ConvertEnvs(comp.Env)

	containerLabels := map[string]string{
		"component": componentName,
		"alias":     *comp.Alias,
	}

	containerConfig := a.Client.GenerateContainerConfig(*comp.Image, comp.Command, comp.Args, envVars, containerLabels)
	return containerConfig
}

func addProjectVolumeToComp(projectVolumeName string, hostConfig *container.HostConfig) *container.HostConfig {
	mount := mount.Mount{
		Type:   mount.TypeVolume,
		Source: projectVolumeName,
		Target: lclient.OdoSourceVolumeMount,
	}
	hostConfig.Mounts = append(hostConfig.Mounts, mount)

	return hostConfig
}
