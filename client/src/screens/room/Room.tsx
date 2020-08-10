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


  const messageListener = (ev: MessageEvent) => {
    const res = JSON.parse(ev.data)
    console.log(JSON.parse(ev.data))

    if (!res || !res.action) {
      console.warn("No action to handle")
      return
    }

    switch (res.action) {
      case "REQUEST": {
        if (!res.data) return

        syncVideoWithServer(res.data)
        return
      }

      case "SYNC": {
        if (!res.data || !res.data.url) return

        syncVideoWithServer(res.data)
        seekVideo(res.data.time)
        return
      }

      case "PLAY_VIDEO": {
        if (!res.data) return

        syncVideoWithServer(res.data)
        seekVideo(res.data.time)
        return
      }

      case "PAUSE_VIDEO": {
        if (!res.data) return

        syncVideoWithServer(res.data)
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

  const syncVideoWithServer = (newVideoData: VideoData) => {
    setVideoData({
      url: newVideoData.url,
      time: newVideoData.time,
      playing: newVideoData.playing,
    })
  }

  const seekVideo = (durationTime: number) => {
    if (durationTime > 0 && playerRef?.current) {
      setIsSeeking(true)
      playerRef.current.seekTo(durationTime, 'seconds')
      return
    }

    console.warn("Failed to seek", videoData)
  }

  const handlePlay = () => {
    if (!playerRef?.current) return

    if (!isSeeking) {
      sendMessage({
        action: "PLAY_VIDEO",
        data: {
          time: playerRef.current.getCurrentTime(),
          url: videoData.url,
          playing: true,
        }
      })
    }

    setIsSeeking(false)
  }

  const handlePause = () => {
    if (!playerRef?.current) return

    if (!isSeeking) {
      sendMessage({
        action: "PAUSE_VIDEO",
        data: {
          time: playerRef.current.getCurrentTime(),
          url: videoData.url,
        }
      })
    }

    setIsSeeking(false)
  }

  const handleSeek = () => {
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