import React from 'react'
import ReactPlayer from 'react-player'
import { RoomProps as Props } from './types'
import { useRoom } from '../../hooks/useRoom'
import c from "./Room.module.scss"
import Playlist from 'components/Playlist'


const Room: React.FC<Props> = () => {
  const {
    isMediaReady, playerRef, videoUrl, videoData, playlist,
    handleMediaReady, handlePause, handlePlay, handleSeek, handleRequestVideo, setVideo, handleMediaEnd,
  } = useRoom()

  return (
    <div className={c.Root}>
      <div className={c.PlaylistContainer}>
        <div className={c.VideoInput}>
          <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
          <button onClick={handleRequestVideo}>Add to playlist</button>
        </div>

        <Playlist
          className={c.List}
          videosUrls={playlist}
        />
      </div>

      <div className={c.VideoContainer}>
        <ReactPlayer
          ref={playerRef}
          playing={videoData.playing && isMediaReady}
          url={videoData.url || ''}
          onSeek={handleSeek}
          onPlay={handlePlay}
          onPause={handlePause}
          onEnded={handleMediaEnd}
          onReady={handleMediaReady}
          controls
          height='100%'
          width='100%'
        />
      </div>
    </div>
  )
}

export default Room
