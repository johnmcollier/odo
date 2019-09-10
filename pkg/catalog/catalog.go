package catalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
)

const DefaultIDPRepo = "https://raw.githubusercontent.com/johnmcollier/iterative-dev-packs/master"
const DefaultIDPCatalog = DefaultIDPRepo + "/index.json"

// CatalogEntry represents an entry in the index.json file for IDPs
type CatalogEntry struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Language    string            `json:"language"`
	Framework   string            `json:"framework"`
	Devpacks    map[string]string `json:"devpacks"`
}

// List lists all the available Iterative-Dev Packs
func List() ([]CatalogEntry, error) {
	// See if we have an index.json already cached
	var idpList []CatalogEntry
	indexJSONFile := path.Join(os.TempDir(), ".kdo", "index.json")
	// Load the index.json file into memory
	var jsonBytes []byte
	if _, err := os.Stat(indexJSONFile); os.IsNotExist(err) {
		jsonBytes, err = downloadIDPs()
		if err != nil {
			return nil, fmt.Errorf("unable to download index.json from %s", DefaultIDPCatalog)
		}
	} else {
		file, err := os.Open(indexJSONFile)
		if err != nil {
			return nil, err
		}
		jsonBytes, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}

	err := json.Unmarshal(jsonBytes, &idpList)
	return idpList, err
}

// Search searches for the IDP in the catalog
func Search(name string) ([]string, error) {
	var result []string
	idpList, err := List()
	if err != nil {
		return nil, errors.Wrap(err, "unable to list Iterative-Dev Packs")
	}

	// do a partial search in all the Iterative-Dev Packs
	for _, idp := range idpList {
		// we only show IDP entries that contain the search term in their name and that have at least one devpack yaml linked
		if strings.Contains(idp.Name, name) {
			result = append(result, idp.Name)
		}
	}

	return result, nil
}

// Exists returns true if the given iterative-dev pack type is valid, false if not
func Exists(idpName string) (bool, error) {

	catalogList, err := List()
	if err != nil {
		return false, errors.Wrapf(err, "unable to list catalog")
	}

	for _, supported := range catalogList {
		if idpName == supported.Name {
			return true, nil
		}
	}
	return false, nil
}

// VersionExists checks if a IDP of the specified name and version exists
func VersionExists(client *kclient.Client, idpName string, idpVersion string) (bool, error) {

	// Loading status
	glog.V(4).Info("Checking Iterative-Dev pack version")

	// Retrieve the catalogList
	catalogList, err := List()
	if err != nil {
		return false, errors.Wrapf(err, "unable to list catalog")
	}

	// Find the IDP and then return true if the version has been found
	for _, supported := range catalogList {
		if idpName == supported.Name {
			// Now check to see if that version matches that IDP's version
			// here we use the AllTags, because if the user somehow got hold of a version that was hidden
			// then it's safe to assume that this user went to a lot of trouble to actually use that version,
			// so let's allow it
			for version := range supported.Devpacks {
				if idpVersion == version {
					return true, nil
				}
			}
		}
	}

	// Else return false if nothing is found
	return false, nil
}

// downloadIDPs downloads the index.json to a temp directory for KDO to access
func downloadIDPs() ([]byte, error) {
	// Download the IDP index.json
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(DefaultIDPCatalog)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return jsonBytes, err
}
