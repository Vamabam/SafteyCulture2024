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

// Edited the function call from original to allow for error throwing
func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// Get all paths from org
	allFolders := f.GetFoldersByOrgID(orgID)

	// Check if Invalid Org
	if len(allFolders) < 1 {
		return []Folder{}, errors.New("invalid orgid")
	}

	children := []Folder{}
	for _, Folder := range allFolders {
		// Split paths by .
		splitPaths := strings.Split(Folder.Paths, ".")
		// Check for invalid path through length of splitPaths

		// Check through split paths if name exists
		// Only len - 1 to ensure only children added
		for i := 0; i < len(splitPaths); i++ {
			// If name is found add whole Folder struct to children
			if splitPaths[i] == name && len(splitPaths) > 1 {
				children = append(children, Folder)
			}
		}

	}
	// If file does not exist in this org then check if in other orgs
	if len(children) < 1 {
		folders := f.folders
		for _, f := range folders {
			if f.Name == name {
				return []Folder{},
					errors.New("Folder does not exist in the specified organization")
			}
		}
		// File isnt in the system at all
		return []Folder{},
			errors.New("Folder doesn't Exist")
	} else {
		return children, nil
	}
}
