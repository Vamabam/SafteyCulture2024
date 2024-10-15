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

// - Inputs: OrgID and name of folder
// - Returns: Array of Child Folders and Errors
// - Gets all the child folders of the given folder within Org
func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// Get all paths from org
	allFolders := f.GetFoldersByOrgID(orgID)

	// Check if Invalid Org
	if len(allFolders) < 1 {
		return []Folder{}, errors.New("OrgID does not exist")
	}

	children := []Folder{}
	for _, Folder := range allFolders {
		// Split paths by .
		splitPaths := strings.Split(Folder.Paths, ".")

		// Check if name exists in split paths
		// len(splitPaths) - 1 as final folder in path has no children
		for i := 0; i < len(splitPaths)-1; i++ {
			// If name is found add whole Folder struct to children
			if splitPaths[i] == name && len(splitPaths) > 1 {
				children = append(children, Folder)
			}
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
		// File isnt in the system at all
		return []Folder{}, errors.New("Folder doesn't exist")
	} else {
		return children, nil
	}
}
