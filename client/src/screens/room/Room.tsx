import React from 'react'
import ReactPlayer from 'react-player'
import { RoomProps as Props } from './types'
import { useRoom } from './useRoom'


const Room: React.FC<Props> = () => {
  const { handleMediaReady, handlePause, handlePlay, handleSeek, isMediaReady,
    videoData, handleRequestVideo, playerRef, setVideo, videoUrl, handleMediaEnd, syncVideoWithServer } = useRoom()

  console.log(videoData)

  return (
    <div>
      <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
      <button onClick={handleRequestVideo}>Send</button>

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
      />

      <button onClick={() => syncVideoWithServer({})}>SYNC</button>
    </div>
  )
}

export default Room
