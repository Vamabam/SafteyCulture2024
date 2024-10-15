package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// feel free to change how the unit test is structured
func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	orgId := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
	orgId2 := uuid.FromStringOrNil("18b2873b-f73b-4b0e-b9d9-4fc4c23646a8")
	orgId3 := uuid.FromStringOrNil("26b2873b-f36b-6b1e-b9d9-4fc4c23646a8")

	sampleFolders := []folder.Folder{
		{Name: "creative-scalphunter", OrgId: orgId, Paths: "creative-scalphunter"},
		{Name: "clear-arclight", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight"},
		{Name: "topical-micromax", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax"},
		{Name: "bursting-lionheart", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart"},
		{Name: "striking-black-panther", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart.striking-black-panther"},
		{Name: "noble-vixen", OrgId: orgId2, Paths: "noble-vixen"},
	}

	tests := [...]struct {
		name    string
		orgId   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		{
			name:    "Valid Case",
			orgId:   orgId,
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "creative-scalphunter", OrgId: orgId, Paths: "creative-scalphunter"},
				{Name: "clear-arclight", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight"},
				{Name: "topical-micromax", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax"},
				{Name: "bursting-lionheart", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart"},
				{Name: "striking-black-panther", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart.striking-black-panther"},
			},
		},
		{
			name:    "Invalid Input",
			orgId:   orgId3,
			folders: sampleFolders,
			want:    []folder.Folder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			get := f.GetFoldersByOrgID(tt.orgId)

			assert.ElementsMatch(t, tt.want, get)
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	orgId := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
	orgId2 := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23646a7")
	// Pulled from sample data
	sampleFolders := []folder.Folder{
		{Name: "creative-scalphunter", OrgId: orgId, Paths: "creative-scalphunter"},
		{Name: "clear-arclight", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight"},
		{Name: "topical-micromax", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax"},
		{Name: "bursting-lionheart", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart"},
		{Name: "striking-black-panther", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart.striking-black-panther"},
		{Name: "noble-vixen", OrgId: orgId2, Paths: "noble-vixen"},
	}

	tests := [...]struct {
		name       string
		orgId      uuid.UUID
		folderName string
		folders    []folder.Folder
		want       []folder.Folder
		wantErr    error
	}{
		{
			name:       "Valid case with children",
			orgId:      orgId,
			folderName: "creative-scalphunter",
			folders:    sampleFolders,
			want: []folder.Folder{
				{Name: "clear-arclight", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight"},
				{Name: "topical-micromax", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax"},
				{Name: "bursting-lionheart", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart"},
				{Name: "striking-black-panther", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart.striking-black-panther"},
			},
			wantErr: nil,
		},
		{
			name:       "Non-existing folder",
			orgId:      orgId,
			folderName: "non-existing-folder",
			folders:    sampleFolders,
			want:       []folder.Folder{},
			wantErr:    errors.New("Folder doesn't exist"),
		},
		{
			name:       "Invalid orgID",
			orgId:      uuid.Must(uuid.NewV4()),
			folderName: "creative-scalphunter",
			folders:    sampleFolders,
			want:       []folder.Folder{},
			wantErr:    errors.New("OrgID does not exist"),
		},
		{
			name:       "Single Folder No Children",
			orgId:      orgId,
			folderName: "striking-black-panther",
			folders:    sampleFolders,
			want:       []folder.Folder{},
			wantErr:    nil,
		},
		{
			name:       "Folder exists in different org",
			orgId:      orgId,
			folderName: "noble-vixen",
			folders:    sampleFolders,
			want:       []folder.Folder{},
			wantErr:    errors.New("Folder does not exist in the specified organization"),
		},
		{
			name:       "Single Child",
			orgId:      orgId,
			folderName: "bursting-lionheart",
			folders:    sampleFolders,
			want: []folder.Folder{
				{Name: "striking-black-panther", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart.striking-black-panther"},
			},
			wantErr: nil,
		},
		{
			name:       "Middle Child",
			orgId:      orgId,
			folderName: "topical-micromax",
			folders:    sampleFolders,
			want: []folder.Folder{
				{Name: "bursting-lionheart", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart"},
				{Name: "striking-black-panther", OrgId: orgId, Paths: "creative-scalphunter.clear-arclight.topical-micromax.bursting-lionheart.striking-black-panther"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			got, err := f.GetAllChildFolders(tt.orgId, tt.folderName)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.want, got)
			}
		})
	}
}
