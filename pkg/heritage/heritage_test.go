package heritage

import (
	"testing"

	"github.com/calamity-m/containerdna/pkg/containers"
	"github.com/containers/image/v5/types"
)

func Test_ValidateHeritage(t *testing.T) {
	type args struct {
		relaxed bool
		child   string
		parents []string
	}
	type testArgs struct {
		name string
		args args
		want bool
	}
	tests := []testArgs{
		{
			name: "Strict check",
			args: args{
				relaxed: false,
				child:   "docker://alpine",
				parents: []string{"docker://alpine", "docker://nginx"},
			},
			want: false,
		},
		{
			name: "Relaxed check",
			args: args{
				relaxed: true,
				child:   "docker://alpine",
				parents: []string{"docker://alpine", "docker://nginx"},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ValidateHeritage(tc.args.relaxed, tc.args.child, tc.args.parents...)
			if err != nil {
				t.Errorf("ValidateHeritage ecountered an error - %v", err)
			}

			if got != tc.want {
				t.Errorf("validateChildParentsImage() = %v, want %v", got, tc.want)
			}
		})
	}
}

func Test_validateChildParentsImageInvalidCases(t *testing.T) {
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
						Layers: []types.BlobInfo{{
							Digest: "A",
						}},
						Name: "",
						Err:  nil,
					},
					{
						Layers: []types.BlobInfo{{
							Digest: "B",
						}},
						Name: "",
						Err:  nil,
					},
				},
			},
			want: false,
		},
		{
			name: "Parent with different layers causes an invalid check. Layers must match from first onwards.",
			args: args{
				relaxed: false,
				child: containers.Image{
					Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}, {Digest: "LAYER3"}},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{{
							Digest: "LAYER2",
						}},
						Name: "",
						Err:  nil,
					},
				},
			},
			want: false,
		},
		{
			name: "Parent with more layers causes an invalid check",
			args: args{
				relaxed: false,
				child: containers.Image{
					Layers: []types.BlobInfo{{Digest: "LAYER1"}},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{{
							Digest: "LAYER1",
						}, {
							Digest: "LAYER2",
						}},
						Name: "",
						Err:  nil,
					},
				},
			},
			want: false,
		},
		{
			name: "Relaxed disabled; should be invalid as one parent matches",
			args: args{
				relaxed: false,
				child: containers.Image{
					Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}, {Digest: "LAYER3"}},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}},
						Name:   "",
						Err:    nil,
					},
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER3"}},
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

func Test_validateChildParentsImageValidCases(t *testing.T) {
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
		{
			name: "Relaxed enabled; should be valid as one parent matches",
			args: args{
				relaxed: true,
				child: containers.Image{
					Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}, {Digest: "LAYER3"}},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}},
						Name:   "",
						Err:    nil,
					},
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER3"}},
						Name:   "",
						Err:    nil,
					},
				},
			},
			want: true,
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

func Test_validateChildParentsImageRelaxed(t *testing.T) {
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
		{
			name: "Equal layers should provide valid check",
			args: args{
				relaxed: false,
				child: containers.Image{
					Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}},
						Name:   "",
						Err:    nil,
					},
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}},
						Name:   "",
						Err:    nil,
					},
				},
			},
			want: true,
		},
		{
			name: "Child and parent share same initial layers",
			args: args{
				relaxed: false,
				child: containers.Image{
					Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}, {Digest: "LAYER3"}},
					Name:   "",
					Err:    nil,
				},
				parents: []containers.Image{
					{
						Layers: []types.BlobInfo{{Digest: "LAYER1"}, {Digest: "LAYER2"}},
						Name:   "",
						Err:    nil,
					},
				},
			},
			want: true,
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
