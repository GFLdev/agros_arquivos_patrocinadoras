import apiClient from '@/services/axios.ts'
import type { AxiosError, AxiosResponse } from 'axios'
import type {
  CategModel,
  UserModel,
  FileModel,
  CreateResponse,
  QueryResponse,
  GetAllResponse,
  GetOneResponse,
} from '@/@types/Responses.ts'
import type {
  UserRequest,
  CategRequest,
  FileRequest,
  UpdateCategRequest,
  UpdateFileRequest,
} from '@/@types/Requests.ts'
import router from '@/router'

/**
 * Lida com erros lançados por requisições Axios e retorna uma resposta padronizada.
 *
 * @param {unknown} error - O objeto de erro lançado por uma requisição Axios.
 * @param {string} entityName - O nome da entidade para formatação das mensagens de erro.
 * @return {Promise<QueryResponse>} Uma Promise resolvendo para um objeto contendo a mensagem de erro e o código de
 * status.
 */
async function handleAxiosError(error: unknown, entityName: string = 'Entidade'): Promise<QueryResponse> {
  const SERVER_ERROR_MSG = 'Erro no servidor. Tente novamente mais tarde.'
  const UNAUTHORIZED_MSG = 'Credenciais inválidas.'
  const CONFLICTM_MSG = `${entityName} já existe.`
  const UNKNOWN_ERROR_MSG = 'Erro desconhecido. Tente novamente mais tarde.'

  const axiosError: AxiosError = error as AxiosError
  if (!axiosError.response) {
    return { message: SERVER_ERROR_MSG, code: 500 }
  }

  if (axiosError.response.status === 401) {
    await router.push({ name: 'login', replace: true })
    return { message: UNAUTHORIZED_MSG, code: 401 }
  } else if (axiosError.response.status === 409) {
    return { message: CONFLICTM_MSG, code: 409 }
  }

  return {
    message: UNKNOWN_ERROR_MSG,
    code: axiosError.response?.status || 500,
  }
}

/**
 * Cria um novo usuário enviando as informações do usuário para o servidor.
 *
 * @param {UserRequest} body - O objeto de requisição contendo os detalhes do usuário para criar um novo usuário.
 * @return {Promise<QueryResponse>} Uma Promise que resolve com a resposta contendo o código de status e a mensagem após
 * a operação.
 */
export async function createUser(body: UserRequest): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<CreateResponse, never> = await apiClient.post('/auth/user', body)
    return { message: res.data.message, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e, 'Usuário')
  }
}

/**
 * Cria uma nova categoria associada ao usuário especificado.
 *
 * @param {string} userId - O identificador único do usuário para o qual a categoria está sendo criada.
 * @param {CategRequest} body - Os dados da requisição contendo as informações sobre a categoria a ser criada.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para um objeto contendo uma mensagem e um código de status,
 * indicando o resultado da operação.
 */
export async function createCategory(userId: string, body: CategRequest): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<CreateResponse, never> = await apiClient.post(`/auth/user/${userId}/category`, body)
    return { message: res.data.message, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e, 'Categoria')
  }
}

/**
 * Cria um novo arquivo para o usuário e a categoria especificados.
 *
 * @param {string} userId - O identificador único do usuário.
 * @param {string} categId - O identificador único da categoria.
 * @param {FileRequest} body - Os detalhes do arquivo a ser criado, encapsulados em um objeto FileRequest.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para um objeto QueryResponse contendo a mensagem de resposta
 * do servidor e o código de status.
 */
export async function createFile(userId: string, categId: string, body: FileRequest): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<CreateResponse, never> = await apiClient.post(
      `/auth/user/${userId}/category/${categId}/file`,
      body,
    )
    return { message: res.data.message, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e, 'Arquivo')
  }
}

/**
 * Recupera uma lista de todos os usuários do servidor.
 *
 * @return {Promise<GetAllResponse<UserModel>>} Uma Promise que resolve para um objeto contendo os dados dos usuários,
 * uma mensagem indicando o sucesso ou falha da operação, e um código de resposta.
 */
export async function getAllUsers(): Promise<GetAllResponse<UserModel>> {
  try {
    const res: AxiosResponse<UserModel[], never> = await apiClient.get('/auth/user')
    return {
      data: res.data ?? null,
      message: res.data ? 'Usuários obtidos com sucesso.' : 'Nenhum usuário encontrado.',
      code: res.data ? 200 : 204,
    }
  } catch (e: unknown) {
    const { message, code } = await handleAxiosError(e)
    return { data: null, message, code }
  }
}

export async function getUserById(userId: string): Promise<GetOneResponse<UserModel>> {
  try {
    const res: AxiosResponse<UserModel, never> = await apiClient.get(`/auth/user/${userId}`)
    return {
      data: res.data ?? null,
      message: res.data ? 'Usuário obtido com sucesso.' : 'Usuário não encontrado.',
      code: res.data ? 200 : 204,
    }
  } catch (e: unknown) {
    const { message, code } = await handleAxiosError(e)
    return { data: null, message, code }
  }
}

/**
 * Busca todas as categorias associadas a um usuário específico.
 *
 * @param {string} userId - O identificador único do usuário cujas categorias precisam ser recuperadas.
 * @return {Promise<GetAllResponse<CategModel>>} Uma Promise que resolve para um objeto de resposta contendo a lista de
 * categorias, uma mensagem e um código de status.
 */
export async function getAllCategories(userId: string): Promise<GetAllResponse<CategModel>> {
  try {
    const res: AxiosResponse<CategModel[], never> = await apiClient.get(`/auth/user/${userId}/category`)
    return {
      data: res.data ?? null,
      message: res.data ? 'Categorias obtidas com sucesso.' : 'Nenhuma categoria encontrada.',
      code: res.data ? 200 : 204,
    }
  } catch (e: unknown) {
    const { message, code } = await handleAxiosError(e)
    return { data: null, message, code }
  }
}

/**
 * Busca uma categoria específica pelo seu ID para um determinado usuário na API.
 *
 * @param {string} userId - O ID do usuário ao qual a categoria pertence.
 * @param {string} categId - O ID da categoria a ser recuperada.
 * @return {Promise<GetOneResponse<CategModel>>} Uma Promise que resolve para um objeto contendo os dados da categoria,
 * uma mensagem e um código de status.
 */
export async function getCategoryById(userId: string, categId: string): Promise<GetOneResponse<CategModel>> {
  try {
    const res: AxiosResponse<CategModel, never> = await apiClient.get(`/auth/user/${userId}/category/${categId}`)
    return {
      data: res.data ?? null,
      message: res.data ? 'Categoria obtida com sucesso.' : 'Categoria não encontrada.',
      code: res.data ? 200 : 204,
    }
  } catch (e: unknown) {
    const { message, code } = await handleAxiosError(e)
    return { data: null, message, code }
  }
}

/**
 * Busca todos os arquivos associados a um usuário e a uma categoria específicos no servidor.
 *
 * @param {string} userId - O identificador único do usuário.
 * @param {string} categId - O identificador único da categoria.
 * @return {Promise<GetAllResponse<FileModel>>} Uma Promise que resolve para um objeto contendo a lista de arquivos,
 * código de status e uma mensagem.
 */
export async function getAllFiles(userId: string, categId: string): Promise<GetAllResponse<FileModel>> {
  try {
    const res: AxiosResponse<FileModel[], never> = await apiClient.get(`/auth/user/${userId}/category/${categId}/file`)
    return {
      data: res.data ?? null,
      message: res.data ? 'Arquivos obtidos com sucesso.' : 'Nenhum arquivo encontrado.',
      code: res.data ? 200 : 204,
    }
  } catch (e: unknown) {
    const { message, code } = await handleAxiosError(e)
    return { data: null, message, code }
  }
}

/**
 * Recupera um arquivo pelo seu ID para um determinado usuário e categoria.
 *
 * @param {string} userId - O ID do usuário ao qual o arquivo pertence.
 * @param {string} categId - O ID da categoria onde o arquivo está localizado.
 * @param {string} fileId - O ID do arquivo a ser recuperado.
 * @return {Promise<GetOneResponse<FileModel>>} Uma Promise que resolve para os dados do arquivo, uma mensagem e um
 * código de resposta.
 */
export async function getFileById(userId: string, categId: string, fileId: string): Promise<GetOneResponse<FileModel>> {
  try {
    const res: AxiosResponse<FileModel, never> = await apiClient.get(
      `/auth/user/${userId}/category/${categId}/file/${fileId}`,
    )
    return {
      data: res.data ?? null,
      message: res.data ? 'Arquivo obtido com sucesso.' : 'Arquivo não encontrado.',
      code: res.data ? 200 : 204,
    }
  } catch (e: unknown) {
    const { message, code } = await handleAxiosError(e)
    return { data: null, message, code }
  }
}

/**
 * Atualiza os detalhes de um usuário com o ID especificado.
 *
 * @param {string} userId - O identificador único do usuário a ser atualizado.
 * @param {UserRequest} body - O corpo da requisição contendo as informações atualizadas do usuário.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para a resposta contendo o status da atualização e uma
 * mensagem.
 */
export async function updateUser(userId: string, body: UserRequest): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<string, never> = await apiClient.patch(`/auth/user/${userId}`, body)
    return { message: res.data, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e)
  }
}

/**
 * Atualiza uma categoria existente para um usuário específico usando os dados fornecidos.
 *
 * @param {string} userId - O ID do usuário ao qual a categoria pertence.
 * @param {string} categId - O ID da categoria a ser atualizada.
 * @param {UpdateCategRequest} body - Os dados para atualizar a categoria.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para a resposta contendo a mensagem e o código de status, ou
 * uma resposta de erro caso a requisição falhe.
 */
export async function updateCategory(
  userId: string,
  categId: string,
  body: UpdateCategRequest,
): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<string, never> = await apiClient.patch(`/auth/user/${userId}/category/${categId}`, body)
    return { message: res.data, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e)
  }
}

/**
 * Atualiza um arquivo associado a um usuário e categoria específicos.
 *
 * @param {string} userId - O ID do usuário que possui o arquivo.
 * @param {string} categId - O ID da categoria à qual o arquivo pertence.
 * @param {string} fileId - O ID do arquivo a ser atualizado.
 * @param {UpdateFileRequest} body - O corpo da requisição contendo os dados do arquivo a serem atualizados.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para um objeto contendo o status da operação, mensagem e
 * código de resposta.
 */
export async function updateFile(
  userId: string,
  categId: string,
  fileId: string,
  body: UpdateFileRequest,
): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<string, never> = await apiClient.patch(
      `/auth/user/${userId}/category/${categId}/file/${fileId}`,
      body,
    )
    return { message: res.data, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e)
  }
}

/**
 * Exclui um usuário pelo seu identificador único.
 *
 * @param {string} userId - O identificador único do usuário a ser excluído.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para um objeto QueryResponse contendo o código de status e a
 * mensagem do servidor.
 */
export async function deleteUser(userId: string): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<string, never> = await apiClient.delete(`/auth/user/${userId}`)
    return { message: res.data, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e)
  }
}

/**
 * Exclui uma categoria associada a um usuário específico.
 *
 * @param {string} userId - O identificador único do usuário.
 * @param {string} categId - O identificador único da categoria a ser excluída.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para um objeto contendo a mensagem de resposta e o código de
 * status.
 */
export async function deleteCategory(userId: string, categId: string): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<string, never> = await apiClient.delete(`/auth/user/${userId}/category/${categId}`)
    return { message: res.data, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e)
  }
}

/**
 * Exclui um arquivo correspondente aos identificadores fornecidos de usuário, categoria e arquivo.
 *
 * @param {string} userId - O identificador único do usuário.
 * @param {string} categId - O identificador único da categoria.
 * @param {string} fileId - O identificador único do arquivo a ser excluído.
 * @return {Promise<QueryResponse>} Uma Promise que resolve para um objeto de resposta contendo uma mensagem de status e
 * um código.
 */
export async function deleteFile(userId: string, categId: string, fileId: string): Promise<QueryResponse> {
  try {
    const res: AxiosResponse<string, never> = await apiClient.delete(
      `/auth/user/${userId}/category/${categId}/file/${fileId}`,
    )
    return { message: res.data, code: res.status }
  } catch (e: unknown) {
    return handleAxiosError(e)
  }
}
