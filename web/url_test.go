package web_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/anno/web"
)

func TestNiceURL(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		Source   string
		Expected string
	}{
		{
			"http://www.matryer.com/",
			"matryer.com",
		},
		{
			"http://www.matryer.com/something",
			"matryer.com/something",
		},
		{
			"www.matryer.com/something",
			"matryer.com/something",
		},

		// long paths are truncated
		{
			"www.matryer.com/something/else/is/looking/pretty/long",
			"matryer.com/something/else/is/lo...",
		},
	}

	for _, test := range tests {
		actual := web.ShortURL(test.Source, 35)
		is.Equal(actual, test.Expected)
	}

}
