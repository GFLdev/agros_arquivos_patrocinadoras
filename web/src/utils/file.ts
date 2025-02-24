import { getFileById } from '@/services/queries.ts'
import type { FileModel, GetOneResponse } from '@/@types/Responses.ts'
import type { Component } from 'vue'
import {
  PhFile,
  PhFileArchive,
  PhFileCsv,
  PhFilePdf,
  PhFileSvg,
  PhFileText,
  PhGif,
  PhImage,
  PhMicrosoftExcelLogo,
  PhMicrosoftPowerpointLogo,
  PhMicrosoftWordLogo,
  PhMusicNote,
  PhVideo,
} from '@phosphor-icons/vue'

const mimeTypeMap: Record<string, Component> = {
  // Imagens
  'image/gif': PhGif,
  'image/svg+xml': PhFileSvg,
  'image/png': PhImage,
  'image/jpeg': PhImage,
  'image/jpg': PhImage,
  'image/webp': PhImage,

  // Vídeos
  'video/mp4': PhVideo,
  'video/x-m4v': PhVideo,
  'video/quicktime': PhVideo,
  'video/x-msvideo': PhVideo,
  'video/x-ms-wmv': PhVideo,
  'video/x-ms-asf': PhVideo,
  'video/x-flv': PhVideo,
  'video/x-matroska': PhVideo,
  'video/x-ms-wmx': PhVideo,
  'video/x-ms-wvx': PhVideo,
  'video/x-ms-wm': PhVideo,

  // Áudios
  'audio/mpeg': PhMusicNote,
  'audio/x-mpeg': PhMusicNote,
  'audio/x-mp3': PhMusicNote,
  'audio/x-wav': PhMusicNote,
  'audio/x-aiff': PhMusicNote,
  'audio/x-aac': PhMusicNote,
  'audio/x-ac3': PhMusicNote,
  'audio/x-caf': PhMusicNote,
  'audio/x-flac': PhMusicNote,
  'audio/x-m4a': PhMusicNote,
  'audio/x-matroska': PhMusicNote,
  'audio/x-ms-wma': PhMusicNote,
  'audio/x-ms-wax': PhMusicNote,

  // PDF
  'application/pdf': PhFilePdf,

  // Arquivos de compressão
  'application/zip': PhFileArchive,
  'application/x-zip-compressed': PhFileArchive,
  'application/x-rar-compressed': PhFileArchive,
  'application/x-7z-compressed': PhFileArchive,
  'application/x-tar': PhFileArchive,
  'application/x-gzip': PhFileArchive,
  'application/x-bzip': PhFileArchive,
  'application/x-bzip2': PhFileArchive,
  'application/x-xz': PhFileArchive,
  'application/x-lzma': PhFileArchive,
  'application/x-lz4': PhFileArchive,
  'application/x-zstd': PhFileArchive,
  'application/vnd.rar': PhFileArchive,

  // Documentos
  'application/msword': PhMicrosoftWordLogo,
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document': PhMicrosoftWordLogo,
  'application/vnd.ms-word': PhMicrosoftWordLogo,
  'application/vnd.oasis.opendocument.text': PhMicrosoftWordLogo,
  'application/vnd.ms-word.document.macroenabled.12': PhMicrosoftWordLogo,
  'application/vnd.ms-word.template.macroenabled.12': PhMicrosoftWordLogo,

  // Planilhas
  'application/vnd.ms-excel': PhMicrosoftExcelLogo,
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': PhMicrosoftExcelLogo,
  'application/vnd.oasis.opendocument.spreadsheet': PhMicrosoftExcelLogo,
  'application/vnd.ms-excel.sheet.macroenabled.12': PhMicrosoftExcelLogo,
  'application/vnd.ms-excel.template.macroenabled.12': PhMicrosoftExcelLogo,
  'application/vnd.ms-excel.addin.macroenabled.12': PhMicrosoftExcelLogo,
  'application/vnd.ms-excel.sheet.binary.macroenabled.12': PhMicrosoftExcelLogo,

  // Apresentações
  'application/vnd.ms-powerpoint': PhMicrosoftPowerpointLogo,
  'application/vnd.openxmlformats-officedocument.presentationml.presentation': PhMicrosoftPowerpointLogo,
  'application/vnd.oasis.opendocument.presentation': PhMicrosoftPowerpointLogo,
  'application/vnd.ms-powerpoint.presentation.macroenabled.12': PhMicrosoftPowerpointLogo,
  'application/vnd.ms-powerpoint.template.macroenabled.12': PhMicrosoftPowerpointLogo,
  'application/vnd.ms-powerpoint.slideshow.macroenabled.12': PhMicrosoftPowerpointLogo,

  // CSV
  'text/csv': PhFileCsv,

  // Texto
  'text/plain': PhFileText,
}

export function toBase64(file: File): Promise<string> {
  return new Promise((resolve, reject): void => {
    const reader: FileReader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = (e: ProgressEvent<FileReader>) => reject(e)
  })
}

export function toBlob(b64Data: string, contentType: string = '', sliceSize: number = 512): Blob {
  // Base64 para binário
  const byteCharacters: string = atob(b64Data)
  // Array com os chunks/slices
  const byteArrays: Uint8Array[] = []

  // Leitura em chunks
  for (let offset: number = 0; offset < byteCharacters.length; offset += sliceSize) {
    const slice: string = byteCharacters.slice(offset, offset + sliceSize)
    const byteNumbers: number[] = new Array(slice.length)

    for (let i: number = 0; i < slice.length; i++) {
      byteNumbers[i] = slice.charCodeAt(i)
    }

    const byteArray: Uint8Array = new Uint8Array(byteNumbers)
    byteArrays.push(byteArray)
  }
  return new Blob(byteArrays, { type: contentType })
}

export async function downloadFile(userId: string, categId: string, fileId: string): Promise<void> {
  try {
    const res: GetOneResponse<FileModel> = await getFileById(userId, categId, fileId)
    if (res.code === 204 || !res.data) {
      return
    }
    const blob: Blob = toBlob(res.data.blob, res.data.mimetype)
    const fileURL: string = URL.createObjectURL(blob)
    const link: HTMLAnchorElement = document.createElement('a')

    link.href = fileURL
    link.download = res.data.name.replace(' ', '_') + res.data.extension
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(fileURL)
  } catch (e: unknown) {
    console.error(e)
  }
}

export function fileIcon(mimetype: string): Component {
  return mimeTypeMap[mimetype] || PhFile
}
