export interface LoginRequest {
  username: string
  password: string
}

export interface DownloadRequest {
  user: string
  category: string
  year: number
}

export interface UserRequest {
  username: string
  name: string
  password: string
}

export interface CategRequest {
  name: string
}
