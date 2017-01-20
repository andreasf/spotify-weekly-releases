package model

type Artist struct {
	Name string
	Id   string
}

type ArtistList []Artist

func (self ArtistList) GetIds() []string {
	ids := make([]string, 0, len(self))

	for _, artist := range self {
		ids = append(ids, artist.Id)
	}

	return ids
}

type Album struct {
	Name        string
	Id          string
	ArtistIds   []string
	ReleaseDate string
	Markets     []string
	Tracks      []Track
}

type AlbumList []Album

func (self AlbumList) RemoveDuplicates() AlbumList {
	filtered := make([]Album, 0, len(self))

	albumsByName := make(map[string]Album)
	for _, album := range self {
		key := album.ArtistIds[0] + ":" + album.Name

		_, exists := albumsByName[key]
		if !exists {
			albumsByName[key] = album
			filtered = append(filtered, album)
		}
	}

	return filtered
}

func (self AlbumList) Remove(albums []Album) AlbumList {
	filtered := make([]Album, 0, len(self))
	toRemove := make(map[string]bool)

	for _, album := range albums {
		toRemove[album.Id] = true
	}

	for _, album := range self {
		_, wantToRemove := toRemove[album.Id]
		if wantToRemove {
			continue
		}

		filtered = append(filtered, album)
	}

	return filtered
}

func (self AlbumList) GetArtistIds() []string {
	artistIds := make([]string, 0, len(self))

	for _, album := range self {
		artistIds = append(artistIds, album.ArtistIds...)
	}

	return artistIds
}

type Track struct {
	Name       string
	Id         string
	ArtistId   string
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

func (self TrackList) RemoveDuplicates() TrackList {
	filtered := make([]Track, 0, len(self))

	tracksByName := make(map[string]Track)
	for _, track := range self {
		key := track.ArtistId + ":" + track.Name

		_, exists := tracksByName[key]
		if !exists {
			tracksByName[key] = track
			filtered = append(filtered, track)
		}
	}

	return filtered
}
