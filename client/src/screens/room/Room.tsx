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
  const [isSeeking, setIsSeeking] = useState(false)

  const { roomID } = useParams()

  //TODO: THIS NEEDS AJUSTMENTS
/* 
  useEffect(() => {
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
      data: {
        url: videoUrl,
      }
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
    console.log("PLAYING...")
    if (!playerRef.current) {
      console.log("No REF to media player found")
      return
    }
    
    if (!isSeeking) {
      console.log("SENDING!!")
      sendMessage({
        action: "PLAY_VIDEO",
        data: {
          time: playerRef.current.getCurrentTime(),
          url: videoData.url,
        }
      })
    }

    setIsSeeking(false)
    playMedia()
  }

  const handlePause = () => {
    if (!playerRef.current) {
      console.log("No REF to media player found")
      return
    }

        setIsSeeking(false)

    sendMessage({
      action: "PAUSE_VIDEO",
      data: {
        time: playerRef.current.getCurrentTime(),
        url: videoData.url,
      }
    })

    pauseMedia()
  }

  const handleSeek = () => {
    console.log("SEEKING")
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
        onSeek={handleSeek}
        onPlay={handlePlay}
        onPause={handlePause}
        onEnded={undefined}
        onReady={handleMediaReady}
      />
    </div>
  )
}

export default Room