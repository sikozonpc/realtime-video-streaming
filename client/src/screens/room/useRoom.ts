import { useState, useRef, useCallback } from "react"
import { VideoData } from "./types"
import ReactPlayer from "react-player"
import useWebsocket, { ActionType } from "hooks/useWebsocket"
import { useParams } from "react-router"


const initialVideoData = {
  time: 0,
  url: "",
  playing: false,
}


export const useRoom = () => {
  const playerRef = useRef<ReactPlayer | null>(null)

  const [videoData, setVideoData] = useState<VideoData>(initialVideoData)
  const [isMediaReady, setMediaReady] = useState(false)
  const [isSeeking, setIsSeeking] = useState(false)
  const [videoUrl, setVideo] = useState("")
  const stateRef = useRef<VideoData>()

  const { roomID } = useParams()

  //Hack to fix a callback react variable not updating https://stackoverflow.com/questions/57847594/react-hooks-accessing-up-to-date-state-from-within-a-callback
  stateRef.current = videoData;


  const seekVideo = (durationTime: number) => {
    if (durationTime > 0 && playerRef?.current) {
      setIsSeeking(true)
      playerRef.current.seekTo(durationTime, 'seconds')
      return
    }

    console.warn("Failed to seek", videoData)
  }

  const messageListener = (ev: MessageEvent) => {
    const res = JSON.parse(ev.data)
    console.log(JSON.parse(ev.data))

    const action: ActionType = res.action

    if (!res || !action) {
      console.warn("No action to handle")
      return
    }

    switch (action) {
      case ActionType.REQUEST: {
        console.log(videoData)

        if (!res.data) return

        const isRequestingNewVideo = stateRef.current?.url != res.data.url
        if (isRequestingNewVideo) return

        syncVideoWithServer(res.data)
        return
      }

      case ActionType.END_VIDEO: {
        return
      }

      case ActionType.SYNC: {
        if (!res.data || !res.data.url) return

        syncVideoWithServer(res.data)
        seekVideo(res.data.time)
        return
      }

      case ActionType.PLAY_VIDEO: {
        if (!res.data) return

        syncVideoWithServer(res.data)
        seekVideo(res.data.time)
        return
      }

      case ActionType.PAUSE_VIDEO: {
        if (!res.data) return

        syncVideoWithServer(res.data)
        return
      }

      default: return console.log("Nothing", res.data)
    }
  }

  const { sendMessage } = useWebsocket(`ws://localhost:7777/ws/${roomID}`, messageListener)


  const handleRequestVideo = () => {
    sendMessage({
      action: ActionType.REQUEST,
      data: {
        url: videoUrl,
      }
    })

    console.log(videoData.url)

    if (!videoData.url) {
      setVideoData({
        url: videoUrl,
      })
    }
  }

  const syncVideoWithServer = useCallback((newVideoData: VideoData) => {
    setVideoData({
      url: newVideoData.url,
      time: newVideoData.time,
      playing: newVideoData.playing,
    })
  }, [videoData])

  const handlePlay = () => {
    if (!playerRef?.current) return

    if (!isSeeking) {
      sendMessage({
        action: ActionType.PLAY_VIDEO,
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
        action: ActionType.PAUSE_VIDEO,
        data: {
          time: playerRef.current.getCurrentTime(),
          url: videoData.url,
          playing: false,
        }
      })
    }

    setIsSeeking(false)
  }

  const handleSeek = () => {
    sendMessage({
      action: ActionType.PLAY_VIDEO,
      data: {
        time: playerRef?.current?.getCurrentTime() || 0,
        url: videoData.url,
        playing: true,
      }
    })
  }

  const handleMediaEnd = () => {
    sendMessage({
      action: ActionType.END_VIDEO,
      data: {}
    })
  }

  const handleMediaReady = (_player: ReactPlayer) => {
    setMediaReady(true)
  }

  return {
    isMediaReady,
    videoData,
    playerRef,
    videoUrl,

    handleMediaEnd,
    handleMediaReady,
    handleSeek,
    handlePause,
    handlePlay,
    handleRequestVideo,
    setVideo,
    syncVideoWithServer,
  }
}
