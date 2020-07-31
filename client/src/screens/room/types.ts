export type RoomData = {
  ID: string,
  available: boolean,
}

export interface RoomProps {
  roomData: RoomData,
}