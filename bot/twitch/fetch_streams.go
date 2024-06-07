package twitch

import "github.com/nicklaw5/helix/v2"

// FetchStreams returns the first 100 livestreams matching gameID and language
func FetchStreams(gameID string, language string) ([]helix.Stream, error) {
	params := &helix.StreamsParams{
		GameIDs: []string{gameID},
		First:   100,
	}

	if language != "" {
		// Only querying a single language at a time because we want all
		// 100 streams to be just that one
		params.Language = []string{language}
	}

	resp, err := Client.GetStreams(params)
	if err != nil {
		return nil, err
	}

	return resp.Data.Streams, err
}
