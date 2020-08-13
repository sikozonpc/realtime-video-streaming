import React from 'react'
import ReactPlayer from 'react-player'
import { RoomProps as Props } from './types'
import { useRoom } from './useRoom'
import c from "./Room.module.scss"


const Room: React.FC<Props> = () => {
  const {
    isMediaReady, playerRef, videoUrl, videoData,
    handleMediaReady, handlePause, handlePlay, handleSeek, handleRequestVideo, setVideo, handleMediaEnd,
  } = useRoom()

  return (
    <div className={c.Root}>
      <div className={c.PlaylistContainer}>
        <div className={c.VideoInput}>
          <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
          <button onClick={handleRequestVideo}>REQUEST TO PLAYLIST</button>
        </div>

        <ul className={c.List}>
          <li>Requesting playlist...</li>
        </ul>
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
