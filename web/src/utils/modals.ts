import { AlertType } from '@/@types/Enumerations.ts'

/**
 * Representa a duração de uma transição ou efeito de animação em milissegundos.
 *
 * Valor: 200 milissegundos
 */
export const TRANSITION_DURATION = 200

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

// Representa um alerta com texto, tipo e duração personalizáveis.
export class Alert {
  // Texto do alerta.
  public text: string
  // Tipo do alerta.
  public type: AlertType
  // Duração do alerta.
  public duration: number
  // Estado de visibilidade do alerta.
  public show: boolean

  /**
   * Cria uma instância da classe Alert com o tipo e a duração especificados.
   *
   * @param {AlertType} [type=AlertType.Info] - O tipo do alerta. O padrão é AlertType.Info.
   * @param {number} [duration=3000] - A duração do alerta em milissegundos. O padrão é 3000.
   * @return {Alert} Uma instância da classe Alert.
   */
  constructor(type: AlertType = AlertType.Info, duration: number = 3000) {
    this.text = ''
    this.type = type
    this.duration = duration
    this.show = false
  }

  /**
   * Exibe um alerta com o texto, tipo e duração especificados.
   *
   * @param {string} text - A mensagem a ser exibida no alerta.
   * @param {AlertType} [type=AlertType.Info] - O tipo do alerta (ex.: Info, Aviso, Erro).
   * @param {number} [duration=3000] - A duração para exibir o alerta, em milissegundos.
   * @return {void} Não retorna nenhum valor.
   */
  handleAlert(text: string, type: AlertType = AlertType.Info, duration: number = 3000): void {
    this.text = text
    this.type = type
    this.duration = duration
    this.show = true
  }
}
