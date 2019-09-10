package kclient

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/redhat-developer/odo-fork/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreatePVC creates a PVC resource in the cluster with the given name, size and
// labels
func (c *Client) CreatePVC(name string, size string, labels map[string]string) (*corev1.PersistentVolumeClaim, error) {
	quantity, err := resource.ParseQuantity(size)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse size: %v", size)
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: quantity,
				},
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
		},
	}

	createdPvc, err := c.KubeClient.CoreV1().PersistentVolumeClaims(c.Namespace).Create(pvc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create PVC")
	}
	return createdPvc, nil
}

// AddPVCToDeployment adds the given PVC to the given Deployment
// at the given path
func (c *Client) AddPVCToDeployment(dep *appsv1.Deployment, pvc string, path string) error {
	volumeName := generateVolumeNameFromPVC(pvc)

	dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc,
			},
		},
	})

	// Validating dc.Spec.Template.Spec.Containers[] is present before dereferencing
	if len(dep.Spec.Template.Spec.Containers) == 0 {
		return fmt.Errorf("DeploymentConfig %s doesn't have any Containers defined", dep.Name)
	}
	dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(dep.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      volumeName,
		MountPath: path,
	},
	)
	return nil
}

// UpdatePVCLabels updates the given PVC with the given labels
func (c *Client) UpdatePVCLabels(pvc *corev1.PersistentVolumeClaim, labels map[string]string) error {
	pvc.Labels = labels
	_, err := c.KubeClient.CoreV1().PersistentVolumeClaims(c.Namespace).Update(pvc)
	if err != nil {
		return errors.Wrap(err, "unable to remove storage label from PVC")
	}
	return nil
}

// DeletePVC deletes the given PVC by name
func (c *Client) DeletePVC(name string) error {
	return c.KubeClient.CoreV1().PersistentVolumeClaims(c.Namespace).Delete(name, nil)
}

// IsAppSupervisorDVolume checks if the volume is a supervisorD volume
func (c *Client) IsAppSupervisorDVolume(volumeName, dcName string) bool {
	if volumeName == getAppRootVolumeName(dcName) {
		return true
	}
	return false
}

// getVolumeNamesFromPVC returns the name of the volume associated with the given
// PVC in the given Deployment Config
func (c *Client) getVolumeNamesFromPVC(pvc string, dep *appsv1.Deployment) []string {
	var volumes []string
	for _, volume := range dep.Spec.Template.Spec.Volumes {

		// If PVC does not exist, we skip (as this is either EmptyDir or "shared-data" from SupervisorD
		if volume.PersistentVolumeClaim == nil {
			glog.V(4).Infof("Volume has no PVC, skipping %s", volume.Name)
			continue
		}

		// If we find the PVC, add to volumes to be returned
		if volume.PersistentVolumeClaim.ClaimName == pvc {
			volumes = append(volumes, volume.Name)
		}

	}
	return volumes
}

// removeVolumeFromDeployment removes the volume from the given Deployment and
// returns true. If the given volume is not found, it returns false.
func removeVolumeFromDeployment(vol string, dep *appsv1.Deployment) bool {
	found := false
	for i, volume := range dep.Spec.Template.Spec.Volumes {
		if volume.Name == vol {
			found = true
			dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes[:i], dep.Spec.Template.Spec.Volumes[i+1:]...)
		}
	}
	return found
}

// removeVolumeMountFromDeployment removes the volumeMount from all the given containers
// in the given Deployment Config and return true. If the given volumeMount is
// not found, it returns false
func removeVolumeMountFromDeployment(vm string, dep *appsv1.Deployment) bool {
	found := false
	for i, container := range dep.Spec.Template.Spec.Containers {
		for j, volumeMount := range container.VolumeMounts {
			if volumeMount.Name == vm {
				found = true
				dep.Spec.Template.Spec.Containers[i].VolumeMounts = append(dep.Spec.Template.Spec.Containers[i].VolumeMounts[:j], dep.Spec.Template.Spec.Containers[i].VolumeMounts[j+1:]...)
			}
		}
	}
	return found
}

// generateVolumeNameFromPVC generates a random volume name based on the name
// of the given PVC
func generateVolumeNameFromPVC(pvc string) string {
	return fmt.Sprintf("%v-%v-volume", pvc, util.GenerateRandomString(nameLength))
}

// addOrRemoveVolumeAndVolumeMount mounts or unmounts PVCs from the given deployment
func addOrRemoveVolumeAndVolumeMount(client *Client, dep *appsv1.Deployment, storageToMount map[string]*corev1.PersistentVolumeClaim, storageUnMount map[string]string) error {

	if len(dep.Spec.Template.Spec.Containers) == 0 || len(dep.Spec.Template.Spec.Containers) > 1 {
		return fmt.Errorf("more than one container found in deployment")
	}

	// find the volume mount to be unmounted from the dc
	for i, volumeMount := range dep.Spec.Template.Spec.Containers[0].VolumeMounts {
		if _, ok := storageUnMount[volumeMount.MountPath]; ok {
			dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(dep.Spec.Template.Spec.Containers[0].VolumeMounts[:i], dep.Spec.Template.Spec.Containers[0].VolumeMounts[i+1:]...)

			// now find the volume to be deleted from the dc
			for j, volume := range dep.Spec.Template.Spec.Volumes {
				if volume.Name == volumeMount.Name {
					dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes[:j], dep.Spec.Template.Spec.Volumes[j+1:]...)
				}
			}
		}
	}

	for path, pvc := range storageToMount {
		err := client.AddPVCToDeployment(dep, pvc.Name, path)
		if err != nil {
			return errors.Wrap(err, "unable to add pvc to deployment")
		}
	}
	return nil
}
