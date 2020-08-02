import React, { useState, useEffect, useRef } from 'react'
import { useParams } from 'react-router'
import ReactPlayer from 'react-player'

import useWebsocket from '../../hooks/useWebsocket'
import { RoomProps as Props, VideoData } from './types'


const Room: React.FC<Props> = () => {
  const playerRef = useRef<ReactPlayer | null>(null)

  const [videoUrl, setVideo] = useState("")
  const [videoData, setVideoData] = useState<VideoData>({
    time: 0,
    url: "",
    playing: false,
  })
  const [isMediaReady, setMediaReady] = useState(false)
  const [isSeeking, setIsSeeking] = useState(false)

  const { roomID } = useParams()


  /*   useEffect(() => {
      const mediaPlayer = playerRef.current
  
      if (videoData.time && mediaPlayer) {
        console.log("SEEKING TO", videoData.time)
        mediaPlayer.seekTo(videoData.time, 'seconds')
        setIsSeeking(true)
      }
    }, [videoData.time]) */

  const messageListener = (ev: MessageEvent) => {
    const res = JSON.parse(ev.data)
    console.log(JSON.parse(ev.data))

    switch (res?.action) {
      case "SYNC": {
        if (!res.data) return

        syncVideo(res.data)
        return
      }

      case "PLAY_VIDEO": {
        if (!res.data) return

        syncVideo(res.data)
        return
      }

      case "PAUSE_VIDEO": {
        if (!res.data) return

        setVideoData({
          url: res.data.url,
          time: res.data.time,
          playing: res.data.playing,
        })
    
        return
      }

      default: console.log("Nothing", res.data)
    }
  }

  const { sendMessage } = useWebsocket(`ws://localhost:7777/ws/${roomID}`, messageListener)


  const handleRequestVideo = () => {
    sendMessage({
      action: "REQUEST",
      data: {
        url: videoUrl,
      }
    })

    setVideoData({
      url: videoUrl,
      time: 0,
      playing: false,
    })
  }

  const syncVideo = (newVideoData: VideoData) => {
    setVideoData({
      url: newVideoData.url,
      time: newVideoData.time,
      playing: newVideoData.playing,
    })

    if (playerRef.current) {
      setIsSeeking(true)
      playerRef.current.seekTo(newVideoData.time, 'seconds')
    }
  }

  const handlePlay = () => {
    console.log("PLAYING...")
    if (!playerRef.current) {
      console.log("No REF to media player found")
      return
    }

    if (!isSeeking) {
      console.log("SENDING PLAY!!")
      sendMessage({
        action: "PLAY_VIDEO",
        data: {
          time: playerRef.current.getCurrentTime(),
          url: videoData.url,
          playing: true,
        }
      })
    }

    console.log(playerRef.current.getCurrentTime())

    setIsSeeking(false)
  }

  const handlePause = () => {
    if (!playerRef.current) {
      console.log("No REF to media player found")
      return
    }

    if (!isSeeking) {
      console.log("SENDING PAUSE!!")
      sendMessage({
        action: "PAUSE_VIDEO",
        data: {
          time: playerRef.current.getCurrentTime(),
          url: videoData.url,
        }
      })
    }

    console.log(playerRef.current.getCurrentTime())

    setIsSeeking(false)
  }

  const handleSeek = () => {
    console.log("SEEKING")
    sendMessage({
      action: "PLAY_VIDEO",
      data: {
        time: playerRef?.current?.getCurrentTime(),
        url: videoData.url,
        playing: true,
      }
    })
  }

  const handleMediaReady = (player: ReactPlayer) => {
    setMediaReady(true)
  }

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
        onEnded={undefined} // TODO: Notify video ended 
        onReady={handleMediaReady}
        controls
      />
    </div>
  )
}

export default Room