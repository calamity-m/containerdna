package heritage

import (
	"testing"

	"github.com/calamity-m/containerdna/pkg/containers"
	"github.com/containers/image/v5/types"
)

func Test_validateChildParentsImage(t *testing.T) {
	type args struct {
		relaxed bool
		child   containers.Image
		parents []containers.Image
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Empty child layers should always be invalid",
			args: args{
				relaxed: false,
				child: containers.Image{
					Layers: []types.BlobInfo{},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{},
						Name:   "",
						Err:    nil,
					},
					{
						Layers: []types.BlobInfo{},
						Name:   "",
						Err:    nil,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateChildParentsImage(tt.args.relaxed, tt.args.child, tt.args.parents...); got != tt.want {
				t.Errorf("validateChildParentsImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
