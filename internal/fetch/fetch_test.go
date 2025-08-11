package fetch

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/nightscout"
	"gabe565.com/trmnl-nightscout/internal/nightscout/testproperties"
	"github.com/hhsnopek/etag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetch_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/entries.json":
			_, _ = w.Write([]byte("[]"))
		default:
			etag := etag.Generate(testproperties.JSON, true)

			if reqEtag := r.Header.Get("If-None-Match"); reqEtag == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			w.Header().Set("Etag", etag)
			_, _ = w.Write(testproperties.JSON)
		}
	}))

	t.Cleanup(server.Close)

	type fields struct {
		config         *config.Config
		tokenChecksum  string
		propertiesEtag string
		entriesEtag    string //nolint:unused
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               *Response
		wantPropertiesEtag string
		wantEntriesEtag    string
		wantErr            require.ErrorAssertionFunc
	}{
		{
			"no url",
			fields{config: &config.Config{UpdateInterval: 5 * time.Minute}},
			args{t.Context()},
			nil,
			"", "",
			require.Error,
		},
		{
			"success",
			fields{config: &config.Config{NightscoutURL: server.URL, UpdateInterval: 5 * time.Minute}},
			args{t.Context()},
			&Response{Properties: testproperties.Properties, Entries: []nightscout.SGVv1{}},
			testproperties.Etag, "",
			require.NoError,
		},
		{
			"same etag",
			fields{
				config:         &config.Config{NightscoutURL: server.URL, UpdateInterval: 5 * time.Minute},
				propertiesEtag: testproperties.Etag,
			},
			args{t.Context()},
			nil,
			testproperties.Etag, "",
			require.Error,
		},
		{
			"different etag",
			fields{
				config:         &config.Config{NightscoutURL: server.URL, UpdateInterval: 5 * time.Minute},
				propertiesEtag: etag.Generate([]byte("test"), true),
			},
			args{t.Context()},
			&Response{Properties: testproperties.Properties, Entries: []nightscout.SGVv1{}},
			testproperties.Etag, "",
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.fields.config)
			require.NoError(t, err)

			f.tokenChecksum = tt.fields.tokenChecksum
			f.propertiesEtag = tt.fields.propertiesEtag

			got, err := f.Do(tt.args.ctx)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantPropertiesEtag, f.propertiesEtag)
			assert.Equal(t, tt.wantEntriesEtag, f.entriesEtag)
		})
	}
}
