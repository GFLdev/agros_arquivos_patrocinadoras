import apiClient from '@/services/axios.ts'
import type { AxiosError, AxiosResponse } from 'axios'
import type { CategModel, UserModel, FileModel } from '@/@types/Responses.ts'

// Função para gerenciar a obtenção dos dados de todos os usuários
export async function getAllUsers(): Promise<UserModel[] | null> {
  try {
    const res: AxiosResponse<UserModel[], never> = await apiClient.get('/auth/user')
    const users: UserModel[] = res.data

    if (users) {
      return users
    } else {
      console.warn('Nenhum usuário recebido')
      return null
    }
  } catch (e: unknown) {
    const error = e as AxiosError
    if (error.response && error.response.status === 401) {
      // TODO: Exiba uma mensagem amigável para o usuário
      console.error('Credenciais inválidas')
    } else {
      // TODO: Tratar erros e exibir mensagens relevantes ao usuário
      console.error('Erro ao buscar usuários:', error.message || error)
    }
  }

  return null
}

// Função para gerenciar a obtenção dos dados de todas as categorias
export async function getAllCategories(userId: string): Promise<CategModel[] | null> {
  try {
    const res: AxiosResponse<CategModel[], never> = await apiClient.get(`/auth/user/${userId}/category`)
    const categories: CategModel[] = res.data

    if (categories) {
      return categories
    } else {
      console.warn('Nenhuma categoria recebida')
      return null
    }
  } catch (e: unknown) {
    const error = e as AxiosError
    if (error.response && error.response.status === 401) {
      // TODO: Exiba uma mensagem amigável para o usuário
      console.error('Credenciais inválidas')
    } else {
      // TODO: Tratar erros e exibir mensagens relevantes ao usuário
      console.error('Erro ao buscar categorias:', error.message || error)
    }
  }

  return null
}

// Função para gerenciar a obtenção dos dados de todos os arquivos
export async function getAllFiles(userId: string, categId: string): Promise<FileModel[] | null> {
  try {
    const res: AxiosResponse<FileModel[], never> = await apiClient.get(`/auth/user/${userId}/category/${categId}/file`)
    const files: FileModel[] = res.data

    if (files) {
      return files
    } else {
      console.warn('Nenhum arquivo recebido')
      return null
    }
  } catch (e: unknown) {
    const error = e as AxiosError
    if (error.response && error.response.status === 401) {
      // TODO: Exiba uma mensagem amigável para o usuário
      console.error('Credenciais inválidas')
    } else {
      // TODO: Tratar erros e exibir mensagens relevantes ao usuário
      console.error('Erro ao buscar arquivos:', error.message || error)
    }
  }

  return null
}
