export interface LoginRequest {
  username: string
  password: string
}

export interface UserRequest {
  username: string
  name: string
  password: string
}

export interface CategRequest {
  name: string
}

export interface FileRequest {
  name: string
  extension: string
  mimetype: string
  content: string
}

export interface UpdateCategRequest extends CategRequest {
  user_id: string
}

export interface UpdateFileRequest extends FileRequest {
  categ_id: string
}
