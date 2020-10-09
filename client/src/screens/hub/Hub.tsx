import React from 'react'
import axios from 'axios'
import { useHistory } from 'react-router'
import c from './Hub.module.scss'

const Hub: React.FC = () => {
  const history = useHistory()

  const handleCreateRoom = () => {
    // Generate in API a random room that doesnt exist
    axios.get("http://localhost:8080/room")
      .then(d => {
        history.push(`/room/${d.data.ID}`)
      })
  }

  return (
    <div className={c.Root}>
      <button onClick={handleCreateRoom}>CREATE ROOM</button>
    </div>
  )
}

export default Hub
