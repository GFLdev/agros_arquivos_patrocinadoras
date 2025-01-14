export interface LoginRequest {
  username: string;
  password: string;
}

export interface DownloadRequest {
  user: string;
  category: string;
  year: number;
}
