export interface LoginResponse {
  token: string
  message: string
  id: string
  name: string
  admin: boolean
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
  updated_at: number
}

export interface CreateResponse {
  id: string
  message: string
}

export interface QueryResponse {
  message: string
  code: number
}

export interface GetAllResponse<T> extends QueryResponse {
  data: T[] | null
}

export interface GetOneResponse<T> extends QueryResponse {
  data: T | null
}
