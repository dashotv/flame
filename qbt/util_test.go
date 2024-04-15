package qbt

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_magentInfoHash(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "dragon's disciple",
			url:  "magnet:?xt=urn:btih:TOQQBFD5PZ3PO63FSELAXR4L2DVON2S6&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&dn=%5BWEB%5D%20Dragon%27s%20Disciple%20%5B1080P%5D%5BHEVC%5D%5Bx265%5D%5B10-bit%5D%5BHFR%5D%5BAAC%5D",
			want: "9ba100947d7e76f77b6591160bc78bd0eae6ea5e",
		},
		{
			name: "way of choices",
			url:  "magnet:?xt=urn:btih:ee1b48e6eed216440e3940f5031ce5ef2bfa0fbf&dn=Ze%20Tian%20Ji%20%28Way%20of%20Choices%29&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce",
			want: "ee1b48e6eed216440e3940f5031ce5ef2bfa0fbf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			assert.NoError(t, err)

			got, err := magnetInfohash(u)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
