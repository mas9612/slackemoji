package slackemoji

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// EmojiConfig represents the configuration to generate new Slack emoji.
type EmojiConfig struct {
	Color  string
	Font   string
	Public bool
}

// EmojiOption is the option used for generating emoji.
type EmojiOption func(c *EmojiConfig)

// Color sets font color of emoji.
func Color(colorCode string) EmojiOption {
	return func(c *EmojiConfig) {
		c.Color = colorCode
	}
}

// Font sets font family of emoji.
func Font(font string) EmojiOption {
	return func(c *EmojiConfig) {
		c.Font = font
	}
}

// Public sets the publicity of new emoji.
func Public(public bool) EmojiOption {
	return func(c *EmojiConfig) {
		c.Public = public
	}
}

// Emoji is the Emoji image in []byte representation.
type Emoji []byte

func buildParams(text string, c *EmojiConfig) url.Values {
	values := url.Values{}
	values.Add("align", "center")
	values.Add("back_color", "FFFFFF00")
	values.Add("size_fixed", "false")
	values.Add("stretch", "true")

	values.Add("color", c.Color)
	values.Add("font", c.Font)
	values.Add("public_fg", strconv.FormatBool(c.Public))
	// values.Add("text", url.QueryEscape(text))
	values.Add("text", text)
	return values
}

// GenerateEmoji generates new Slack emoji and returns it.
func GenerateEmoji(text string, options ...EmojiOption) (Emoji, error) {
	config := &EmojiConfig{
		Color:  "EC71A1FF",
		Font:   "notosans-mono-bold",
		Public: false,
	}
	for _, o := range options {
		o(config)
	}

	endpoint, err := url.Parse(Endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse endpoint")
	}

	values := buildParams(text, config)
	req := &http.Request{
		Method: http.MethodGet,
		URL:    endpoint,
	}
	req.URL.RawQuery = values.Encode()
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from server")
	}
	defer res.Body.Close()

	emoji, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	return emoji, nil
}
