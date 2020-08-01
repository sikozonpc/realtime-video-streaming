import React, { useState, useEffect } from 'react'
import { useParams } from 'react-router'
import ReactPlayer from 'react-player'

import useWebsocket from '../../hooks/useWebsocket'
import { RoomProps as Props } from './types'


const Room: React.FC<Props> = () => {
  const [sessionData, setSessionData] = useState()
  const [videoUrl, setVideo] = useState("")

  const [videoData, setVideoData] = useState({
    time: 0,
    url: "",
  })

  const { roomID } = useParams()

  const messageListener = (ev: MessageEvent) => {
    const res = JSON.parse(ev.data)
    console.log(JSON.parse(ev.data))

    switch (res?.action) {
      case "PLAY_VIDEO": {
        console.log(res.data)
        setVideoData({
          url: res.data.url,
          time: res.data.time,
        })
      }

      //case "PAUSE_VIDEO": {
        //  return setVideoData(prev => ({
        //  url: prev.url,
        // time: res.data.time,
        //}))
        //}

      case "SYNC_VIDEO": {
        // handleSeek()
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

  const handlePlay = () => {
    sendMessage({
      action: "PLAY_VIDEO",
    })
  }

  const handlePause = () => {
    sendMessage({
      action: "PAUSE_VIDEO",
    })
  }

  const handleSeek = () => {}

  return (
    <div>
      <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
      <button onClick={handleRequestVideo}>Send</button>


      <ReactPlayer
          url={videoData.url}
          onPlay={handlePlay}
          onPause={handlePause}
          onEnded={undefined}
      />
    </div>
  )
}

export default Room