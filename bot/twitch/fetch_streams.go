package twitch

import "github.com/nicklaw5/helix/v2"

// FetchStreams returns the first 100 livestreams matching gameID and language
func FetchStreams(gameID string, language string) ([]helix.Stream, error) {
	resp, err := Client.GetStreams(&helix.StreamsParams{
		GameIDs: []string{gameID},
		// Only querying a single language at a time because we want all
		// 100 streams to be just that one
		Language: []string{language},
		First:    100,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data.Streams, err
}
