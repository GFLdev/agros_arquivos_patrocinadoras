export interface ErrorResponse {
  error: string
  message: string
}

export interface LoginResponse {
  user: string
  authenticated: boolean
}

export interface UserModel {
  user_id: string
  username: string
  name: string
  password: string
  updated_at: string
}

export interface CategModel {
  categ_id: string
  user_id: string
  name: string
  updated_at: string
}
export interface FileModel {
  file_id: string
  categ_id: string
  name: string
  extension: string
  mimetype: string
  blob: string
  updated_at: string
}
