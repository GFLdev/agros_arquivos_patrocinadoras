<div align="center">

[<img src="./public/logo.svg" alt="Agros UFV" width="60%" />](https://www.agros.org.br/)

[<img  src="https://img.shields.io/static/v1?label=license&message=MIT&color=ffffff&labelColor=005cb2" alt="License">](https://github.com/GFLdev/agros_arquivos_patrocinadoras/blob/main/LICENSE)

</div>

# Agros - Arquivos Patrocinadoras

Página para obtenção de arquivos por cada das patrocinadoras do Agros.

## Tecnologias e Frameworks

### Frontend (Node)

- [VueJS](https://vuejs.org/)
- [TailwindCSS](https://tailwindcss.com/)
- [Phosphor Icons](https://phosphoricons.com/)
- [Vite](https://vite.dev/)

### Backend (Go)

- [Echo](https://echo.labstack.com/)
- [Zap](https://github.com/uber-go/zap)

## Setup

```bash
git clone https://github.com/GFLdev/agros_arquivos_patrocinadoras
cd ./agros_arquivos_patrocinadoras
go mod tidy
cd ./web
npm install
```

### Compilição para Desenvolvimento

- Backend:

```bash
go run .
```

- Frontend:

```sh
npm run dev
```

### Compilição para Produção

- Backend:

```bash
go build . -o bin/
```

- Frontend:

```bash
npm run build
```

### Lint com [ESLint](https://eslint.org/)

```bash
npm run lint
```

### Formatação com [Prettier](https://prettier.io/)

```bash
npm run format
```