package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	// Check if same
	if name == dst {
		return []Folder{},
			errors.New("cannot move a folder to itself")
	}
	// Get org Id of source and dest
	srcFolder := f.getFolder(name)
	dstFolder := f.getFolder(dst)

	// Error Checks
	if srcFolder.OrgId == uuid.Nil {
		return []Folder{},
			errors.New("source folder does not exist")
	}

	if dstFolder.OrgId == uuid.Nil {
		return []Folder{},
			errors.New("destination folder does not exist")
	}

	if srcFolder.OrgId != dstFolder.OrgId {
		return []Folder{},
			errors.New("cannot move a folder to a different organization")
	}

	// Get all child folders
	childFolders, _ := f.GetAllChildFolders(srcFolder.OrgId, name)

	// Check if destination is a child of Source
	for _, folder := range childFolders {
		if folder.Name == dst {
			return []Folder{},
				errors.New("cannot move a folder to a child of itself")
		}
	}
	// Get Destination Path
	dstPathSplit := strings.Split(dstFolder.Paths, ".")
	// Go through paths and append to dest path
	res := []Folder{}
	allFolders := f.GetFoldersByOrgID(srcFolder.OrgId)
	for _, folder := range allFolders {
		// Split paths by .
		splitPaths := strings.Split(folder.Paths, ".")
		// Find in path and splice in dest
		for i := 0; i < len(splitPaths); i++ {
			if splitPaths[i] == name {
				// Append child path to destPath
				splitPaths = append(dstPathSplit, splitPaths[i:]...)
				break
			}
		}
		// Join Paths back together and add to Folder struct
		joined := strings.Join(splitPaths, ".")
		folder.Paths = joined

		res = append(res, folder)
	}

	// Append all unedited orginastions to modified data system
	for _, folder := range f.folders {
		if folder.OrgId != srcFolder.OrgId {
			res = append(res, folder)
		}
	}

	return res, nil
}

// Helper function to get folder from name
func (f *driver) getFolder(name string) Folder {
	folder := Folder{}
	folders := f.folders
	for _, f := range folders {
		if f.Name == name {
			folder = f
		}
	}
	return folder
}
