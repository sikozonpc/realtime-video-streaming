import React, { useState, useEffect } from 'react'
import useWebsocket from '../../hooks/useWebsocket'
import { RoomProps as Props } from './types'
import { useParams } from 'react-router'


const Room: React.FC<Props> = () => {
  const [sessionData, setSessionData] = useState()
  const [videoUrl, setVideo] = useState("")

  const { roomID } = useParams()

  const messageListener = (ev: MessageEvent) => {
    console.log(ev)
    console.log(JSON.parse(ev.data))
    setSessionData(JSON.parse(ev.data))
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
  }

  return (
    <div>
      <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
      <button onClick={handleRequestVideo}>Send</button>
    </div>
  )
}

export default Room