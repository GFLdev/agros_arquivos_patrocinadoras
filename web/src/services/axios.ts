import axios, { type AxiosInstance, type InternalAxiosRequestConfig } from 'axios'

// A instância do Axios configurada para interagir com a API.
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
})

// Interceptor para adicionar o token.
apiClient.interceptors.request.use((config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
  const token: string | null = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default apiClient
