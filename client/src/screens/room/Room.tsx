import React, { useState, useEffect, useRef } from 'react'
import { useParams } from 'react-router'
import ReactPlayer from 'react-player'

import useWebsocket from '../../hooks/useWebsocket'
import { RoomProps as Props } from './types'


const Room: React.FC<Props> = () => {
  const playerRef = useRef<ReactPlayer | null>(null)

  const [sessionData, setSessionData] = useState()
  const [videoUrl, setVideo] = useState("")
  const [videoData, setVideoData] = useState({
    time: 0,
    url: "",
  })
  const [isMediaReady, setMediaReady] = useState(false)
  const [mediaPlaying, setMediaPlaying] = useState(false)

  const { roomID } = useParams()

  useEffect(() => {
    const loadedPlayer = playerRef.current

    if (videoData.time && loadedPlayer) {
      loadedPlayer.seekTo(videoData.time)
    }
  }, [videoData.time, playerRef])

  const messageListener = (ev: MessageEvent) => {
    const res = JSON.parse(ev.data)
    console.log(JSON.parse(ev.data))

    switch (res?.action) {
      case "PLAY_VIDEO": {
        console.log(res.data)

        if (!res.data) return

        setVideoData({
          url: res.data.url,
          time: Number(res.data.time),
        })

        playMedia()
        return
      }

      case "PAUSE_VIDEO": {
        if (!res.data) return

        syncVideo(res.data)
        return
      }

      default: console.log("Nothing", res.data)
    }

    setSessionData(res)
  }
  const { sendMessage } = useWebsocket(`ws://localhost:7777/ws/${roomID}`, messageListener)

  useEffect(() => {
    console.log('got new sessionData: ', sessionData)
  }, [sessionData])

  const handleRequestVideo = () => {
    sendMessage({
      action: "REQUEST",
      data: videoUrl
    })

    setVideoData({
      url: videoUrl,
      time: 0,
    })
  }

  const syncVideo = (newVideoData: any) => {
    setVideoData({
      url: newVideoData.url,
      time: newVideoData.time,
    })
    setMediaPlaying(newVideoData.playing)
  }

  const handlePlay = () => {
    sendMessage({
      action: "PLAY_VIDEO",
    })

    playMedia()
  }

  const handlePause = () => {
    sendMessage({
      action: "PAUSE_VIDEO",
    })
  }

  const playMedia = () => setMediaPlaying(true)
  const pauseMedia = () => setMediaPlaying(false)
  const handleMediaReady = (player: ReactPlayer) => {
    setMediaReady(true)
  }

  return (
    <div>
      <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
      <button onClick={handleRequestVideo}>Send</button>

      <ReactPlayer
        ref={playerRef}
        playing={mediaPlaying && isMediaReady}
        url={videoData.url || ''}
        onPlay={handlePlay}
        onPause={handlePause}
        onEnded={undefined}
        onReady={handleMediaReady}
      />
    </div>
  )
}

export default Room