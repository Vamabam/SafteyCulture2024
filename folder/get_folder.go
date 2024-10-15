package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}

// - Gets all the child folders of the given folder within Org
// - Inputs: OrgID and name of folder
// - Returns: Array of Child Folders and Errors
func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// Get all paths from org
	allFolders := f.GetFoldersByOrgID(orgID)

	// Check if invalid org
	if len(allFolders) < 1 {
		return []Folder{}, errors.New("OrgID does not exist")
	}

	// Construct the path prefix to search for
	folder := f.getFolder(name)
	namePrefix := folder.Paths + "."

	children := []Folder{}
	// Iterate over all folders in the org and find children folders
	for _, folder := range allFolders {
		// Check if the folder path starts with "name."
		if strings.HasPrefix(folder.Paths, namePrefix) {
			children = append(children, folder)
		}
	}

	// If file does not exist in this org then check other orgs
	if len(children) < 1 {
		folders := f.folders
		for _, f := range folders {
			if f.Name == name && f.OrgId != orgID {
				return []Folder{},
					errors.New("Folder does not exist in the specified organization")
			} else if f.Name == name {
				// Has no children
				return []Folder{}, nil
			}
		}
		// File isn't in the system at all
		return []Folder{}, errors.New("Folder doesn't exist")
	} else {
		return children, nil
	}
}
