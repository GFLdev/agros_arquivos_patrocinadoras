import { AlertType } from '@/@types/Enumerations.ts'

/**
 * Mapeia um código de status HTTP para um tipo de alerta correspondente.
 *
 * @param {number} code - O código de status HTTP a ser mapeado.
 * @return {AlertType} O tipo de alerta correspondente ao código de status fornecido.
 */
export function codeToAlertType(code: number): AlertType {
  if (code >= 200 && code < 300 && code !== 204) {
    return AlertType.Success
  }
  if (code >= 400 && code < 500) {
    return AlertType.Warning
  }
  if (code >= 500) {
    return AlertType.Error
  }
  return AlertType.Info
}
