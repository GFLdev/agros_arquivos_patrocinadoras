New-Item .\bin\web -ItemType Directory -ea 0

# Buid do Backend
go mod tidy

go build -o .\bin\reset_passwd.exe .\cmd\admin
go build -o .\bin\agros_patrocinadoras.exe .

# Build do frontend
Set-Location .\web
npm install
npm run build
Move-Item -Path .\dist\ -Destination ..\bin\web\

Set-Location ..\