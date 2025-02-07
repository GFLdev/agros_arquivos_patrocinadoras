export interface LoginRequest {
  name: string
  password: string
}

export interface DownloadRequest {
  user: string
  category: string
  year: number
}
