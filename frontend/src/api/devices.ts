export interface Device {
  device_id: string
  display_name: string
  active_state: string
  online: boolean
  lat: number
  lng: number
  heading: number
  speed: number
  dt_tracker: string
  drive_status: string
}
