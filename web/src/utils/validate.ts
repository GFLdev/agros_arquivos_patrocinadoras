/**
 * Valida um nome de usuário para garantir que ele atenda a critérios específicos.
 *
 * @param {string | null | undefined} username - O nome de usuário a ser validado. Deve ser uma string não nula, não
 * indefinida, com um tamanho entre 4 e 16 caracteres.
 * @return {boolean} Retorna true se o nome de usuário for válido; caso contrário, retorna false.
 */
export function validateUsername(username: string | null | undefined): boolean {
  return !!username && username.length >= 4 && username.length <= 16
}

/**
 * Valida se a senha fornecida atende aos critérios exigidos.
 * A senha é considerada válida se não for nula ou indefinida e tiver pelo menos 4 caracteres de comprimento.
 *
 * @param {string | null | undefined} password - A string da senha a ser validada.
 * @return {boolean} Retorna true se a senha atender aos critérios, caso contrário, retorna false.
 */
export function validatePassword(password: string | null | undefined): boolean {
  return !!password && password.length >= 4
}

/**
 * Compara duas strings de senha para verificar se são iguais.
 *
 * @param {string} passwd - A string da senha original.
 * @param {string} confirm - A string de confirmação da senha para comparar com a original.
 * @return {boolean} Retorna true se ambas as strings forem iguais, caso contrário retorna false.
 */
export function checkPasswords(passwd: string, confirm: string): boolean {
  return passwd === confirm
}

/**
 * Verifica se o arquivo fornecido está vazio. Um arquivo é considerado vazio se seu tamanho for 0 e não tiver nome.
 *
 * @param {File} file - O arquivo a ser verificado se está vazio.
 * @return {boolean} Retorna true se o arquivo estiver vazio, caso contrário, false.
 */
export function isFileEmpty(file: File): boolean {
  return file.size === 0 && file.name === ''
}
