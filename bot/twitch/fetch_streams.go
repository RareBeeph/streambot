package twitch

import "github.com/nicklaw5/helix/v2"

// FetchStreams returns the first 100 livestreams matching gameID
func FetchStreams(gameID string) ([]helix.Stream, error) {
	resp, err := Client.GetStreams(&helix.StreamsParams{
		GameIDs: []string{gameID},
		First:   100,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data.Streams, err
}
