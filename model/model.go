package model

type Artist struct {
	Name string
	Id   string
}

type Album struct {
	Name        string
	Id          string
	ArtistId    string
	ReleaseDate string
	Markets     []string
	Tracks      []Track
}

type AlbumList []Album

func (self AlbumList) RemoveDuplicates() AlbumList {
	filtered := make([]Album, 0, len(self))

	albumsByName := make(map[string]Album)
	for _, album := range self {
		key := album.ArtistId + ":" + album.Name

		_, exists := albumsByName[key]
		if !exists {
			albumsByName[key] = album
			filtered = append(filtered, album)
		}
	}

	return filtered
}

type Track struct {
	Name       string
	Id         string
	DurationMs int
}

func (self *Album) GetSampleTrack() *Track {
	if len(self.Tracks) > 3 {
		return &self.Tracks[2]
	}

	if len(self.Tracks) > 0 {
		return &self.Tracks[0]
	}

	return nil
}

type UserProfile struct {
	Id      string
	Country string
}

type TrackList []Track

func (self TrackList) GetUris() []string {
	ids := make([]string, 0, len(self))

	for _, track := range self {
		ids = append(ids, "spotify:track:"+track.Id)
	}

	return ids
}
