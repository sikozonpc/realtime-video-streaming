import React from 'react'
import axios from 'axios'
import { useHistory } from 'react-router'

const Hub: React.FC = () => {
  const history = useHistory()

  const handleCreateRoom = () => {
    // Generate in API a random room that doesnt exist
    axios.get("http://localhost:7777/room")
      .then(d => {
        history.push(`/room/${d.data.ID}`)
      })
  }

  return (
    <div className="App">
      <button onClick={handleCreateRoom}>Create room</button>
    </div>
  )
}

export default Hub