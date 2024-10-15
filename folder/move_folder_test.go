package folder_test

import (
	"errors"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()

	orgId := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23643a7")
	orgId2 := uuid.FromStringOrNil("38b9879b-f73b-4b0e-b9d9-4fc4c23646a7")

	sampleFolders := []folder.Folder{
		{Name: "alpha", OrgId: orgId, Paths: "alpha"},
		{Name: "bravo", OrgId: orgId, Paths: "alpha.bravo"},
		{Name: "charlie", OrgId: orgId, Paths: "alpha.bravo.charlie"},
		{Name: "delta", OrgId: orgId, Paths: "alpha.delta"},
		{Name: "echo", OrgId: orgId, Paths: "alpha.delta.echo"},
		{Name: "foxtrot", OrgId: orgId2, Paths: "foxtrot"},
		{Name: "golf", OrgId: orgId, Paths: "golf"},
	}

	tests := [...]struct {
		name    string
		src     string
		dst     string
		folders []folder.Folder
		want    []folder.Folder
		wantErr error
	}{
		// Test Cases
		{
			name:    "Valid case multiple Paths",
			src:     "bravo",
			dst:     "delta",
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "alpha", OrgId: orgId, Paths: "alpha"},
				{Name: "bravo", OrgId: orgId, Paths: "alpha.delta.bravo"},
				{Name: "charlie", OrgId: orgId, Paths: "alpha.delta.bravo.charlie"},
				{Name: "delta", OrgId: orgId, Paths: "alpha.delta"},
				{Name: "echo", OrgId: orgId, Paths: "alpha.delta.echo"},
				{Name: "foxtrot", OrgId: orgId2, Paths: "foxtrot"},
				{Name: "golf", OrgId: orgId, Paths: "golf"},
			},
			wantErr: nil,
		},
		{
			name:    "Valid case signle path",
			src:     "bravo",
			dst:     "golf",
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "alpha", OrgId: orgId, Paths: "alpha"},
				{Name: "bravo", OrgId: orgId, Paths: "golf.bravo"},
				{Name: "charlie", OrgId: orgId, Paths: "golf.bravo.charlie"},
				{Name: "delta", OrgId: orgId, Paths: "alpha.delta"},
				{Name: "echo", OrgId: orgId, Paths: "alpha.delta.echo"},
				{Name: "foxtrot", OrgId: orgId2, Paths: "foxtrot"},
				{Name: "golf", OrgId: orgId, Paths: "golf"},
			},
			wantErr: nil,
		},
		{
			name:    "Move to Child",
			src:     "bravo",
			dst:     "charlie",
			folders: sampleFolders,
			want:    []folder.Folder{},
			wantErr: errors.New("cannot move a folder to a child of itself"),
		},
		{
			name:    "Move to itself",
			src:     "bravo",
			dst:     "bravo",
			folders: sampleFolders,
			want:    []folder.Folder{},
			wantErr: errors.New("cannot move a folder to itself"),
		},
		{
			name:    "Move to another organisation",
			src:     "bravo",
			dst:     "foxtrot",
			folders: sampleFolders,
			want:    []folder.Folder{},
			wantErr: errors.New("cannot move a folder to a different organization"),
		},
		{
			name:    "Invalid Source",
			src:     "invalid_folder",
			dst:     "foxtrot",
			folders: sampleFolders,
			want:    []folder.Folder{},
			wantErr: errors.New("source folder does not exist"),
		},
		{
			name:    "Invalid Destination",
			src:     "bravo",
			dst:     "invalid_folder",
			folders: sampleFolders,
			want:    []folder.Folder{},
			wantErr: errors.New("destination folder does not exist"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			got, err := f.MoveFolder(tt.src, tt.dst)

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
