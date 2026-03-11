package tpm_test

/*
func TestFdeStatus(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg     testutils.MockConfig
		want    string
		wantErr bool
	}{
		"Success default value": {
			want: "active",
		},
		"Success custom value": {
			cfg:  testutils.MockConfig{FdeStatusValue: "inactive"},
			want: "inactive",
		},

		"Error from snapd": {
			cfg:     testutils.MockConfig{FdeStatusError: true},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			ctx := testutils.ContextLoggerWithDebug(t)
			mockClient := testutils.NewMockSnapdClient(tc.cfg)

			res, err := tpm.FdeStatus(ctx, mockClient)

			if tc.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
			is.Equal(tc.want, res)
		})
	}
}
*/
