package folder

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
)

//   - Moves source folder and all children to destination while
//     keeping folder hierarchy
//   - Inputs: Source folder, destination folder
//   - Returns: Full data structure with folders moved to destination
func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	// Check if source and dest are the same value
	if name == dst {
		return []Folder{},
			errors.New("cannot move a folder to itself")
	}
	// Get orgID of source and dest
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

	allFolders := f.GetFoldersByOrgID(srcFolder.OrgId)

	// Get all children of the source folder and map to hashmap
	childFoldersMap := make(map[string]bool)
	childFolders, _ := f.GetAllChildFolders(srcFolder.OrgId, name)
	for _, folder := range childFolders {
		childFoldersMap[folder.Name] = true
	}

	// Check if destination is a child of the source
	if childFoldersMap[dst] {
		return []Folder{}, errors.New("cannot move a folder to a child of itself")
	}

	// Prepare the destination path prefix
	dstPathPrefix := dstFolder.Paths + "."
	// Prepare the result list
	res := []Folder{}

	// Update paths for all folders that need to be moved
	for _, folder := range allFolders {
		// If the folder's path starts with the source folder's path, update it
		if strings.HasPrefix(folder.Paths, srcFolder.Paths) {
			newPath := strings.Replace(folder.Paths, srcFolder.Paths, dstPathPrefix+name, 1)
			folder.Paths = newPath
			res = append(res, folder)
		} else {
			// Keep the other folders unchanged
			res = append(res, folder)
		}
	}

	// Append all unedited orginastions to modified data system
	for _, folder := range f.folders {
		if folder.OrgId != srcFolder.OrgId {
			res = append(res, folder)
		}
	}

	return res, nil
}

// - Gets the Folder data of the folder that has the name given
// - Input: Name of folder
// - Returns: Folder data type
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
