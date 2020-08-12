package hub

// Playlist is a collection of videos
type Playlist []VideoData

func (p Playlist) Unqueue() Playlist {
	_, p = p[0], p[1:]
	return p
}

func (p Playlist) Enqueue(video VideoData) Playlist {
	return append(p, video)
}
