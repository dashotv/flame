package qbt

import (
	"fmt"
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

func Test_httpInfoHash(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "martial god asura",
			url:  "http://10.0.4.41:9117/dl/animetosho/?jackett_apikey=hcedgf1s4w9etni4pkjb66vh291efylw&path=Q2ZESjhJTHhYZ3Z3d3MxTHZVZ0UtbHJsWVU4OEh2emk3LVVadGhySDQ3bVlkVWhIYWdPUTcyeFhFUzB4ako2NTdjZU5ZQ25CYndGOHFJTGhDTHNsdW82WlRJdUhlQTAwSXBNTmIxQWRaRUlUTVYxSjVvTzRRNjBDTGhPdFBWMUNiVlMzN245M2JPQ1pSdmowbDltaE9tT2RkUkFwNFVmYlNWT1l6RGZWTURaRDR1OVctWG1rVVBlMWw3VGhaRHB5QjBUZV9HSkFxUHJDYkxLMkozeTB6M1c3MlI0UkxHVmd5QWxTRm10VGt4VnJqMmVlR2d3Y1A2MTcweDhzRWstLUJBa1VwVjF2WFU3SkVESklTcTdMb3pKMUZ2Y1VOcUdtUnV0LW1oZTdXaFFSOGxPNUw1R2hGbjNfYmlaWU9CeGtEcjI2cmR6eVZra2xMa0ZIQ0JjUTc2TUtiTkU&file=%5BHuangSubs%5D+%E4%BF%AE%E7%BD%97%E6%AD%A6%E7%A5%9E+%7C+Xiuluo+Wu+Shen+%7C+Martial+God+Asura+S1+(4K+%7C+2160p)",
			want: "9ba100947d7e76f77b6591160bc78bd0eae6ea5e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			assert.NoError(t, err)

			got, err := httpInfohash(u)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_downloadJackett(t *testing.T) {
	s := "http://10.0.4.41:9117/dl/animetosho/?jackett_apikey=hcedgf1s4w9etni4pkjb66vh291efylw&path=Q2ZESjhJTHhYZ3Z3d3MxTHZVZ0UtbHJsWVU4OEh2emk3LVVadGhySDQ3bVlkVWhIYWdPUTcyeFhFUzB4ako2NTdjZU5ZQ25CYndGOHFJTGhDTHNsdW82WlRJdUhlQTAwSXBNTmIxQWRaRUlUTVYxSjVvTzRRNjBDTGhPdFBWMUNiVlMzN245M2JPQ1pSdmowbDltaE9tT2RkUkFwNFVmYlNWT1l6RGZWTURaRDR1OVctWG1rVVBlMWw3VGhaRHB5QjBUZV9HSkFxUHJDYkxLMkozeTB6M1c3MlI0UkxHVmd5QWxTRm10VGt4VnJqMmVlR2d3Y1A2MTcweDhzRWstLUJBa1VwVjF2WFU3SkVESklTcTdMb3pKMUZ2Y1VOcUdtUnV0LW1oZTdXaFFSOGxPNUw1R2hGbjNfYmlaWU9CeGtEcjI2cmR6eVZra2xMa0ZIQ0JjUTc2TUtiTkU&file=%5BHuangSubs%5D+%E4%BF%AE%E7%BD%97%E6%AD%A6%E7%A5%9E+%7C+Xiuluo+Wu+Shen+%7C+Martial+God+Asura+S1+(4K+%7C+2160p)"
	u, err := url.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, s, u.String())

	file, err := downloadURL(u)
	assert.NoError(t, err)
	assert.NotNil(t, file)
	fmt.Printf("file: %s\n", file)
}
