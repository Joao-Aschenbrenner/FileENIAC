<#
Script para desinstalar e reinstalar o ENIAC Workspace

1. Desinstala completamente o ENIAC Workspace
2. Remove todos os arquivos residuais
3. Instala a versão mais recente
#>

# Desinstalar completamente
Write-Host "Desinstalando ENIAC Workspace..."

# Remover do registro
if (Test-Path "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\ENIAC Workspace") {
    Remove-Item "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\ENIAC Workspace" -Recurse -Force
}

# Remover arquivos do programa
$programFiles = "$env:ProgramFiles\ENIAC Workspace"
if (Test-Path $programFiles) {
    Remove-Item $programFiles -Recurse -Force
}

# Remover arquivos do usuário
$localAppData = "$env:LOCALAPPDATA\ENIAC Workspace"
if (Test-Path $localAppData) {
    Remove-Item $localAppData -Recurse -Force
}

# Remover atalhos
$startMenu = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\ENIAC Workspace"
if (Test-Path $startMenu) {
    Remove-Item $startMenu -Recurse -Force
}

# Remover arquivos temporários
$tempFiles = "$env:TEMP\ENIAC Workspace*"
if (Test-Path $tempFiles) {
    Remove-Item $tempFiles -Force
}

# Instalar versão mais recente
Write-Host "Instalando ENIAC Workspace..."

# Caminho para o instalador (substitua pelo caminho real)
$installerPath = "C:\Users\USUARIO\Documents\GitWrkspc\eniac-workspace\build\ENIAC_Workspace_Setup.exe"

if (Test-Path $installerPath) {
    Start-Process -FilePath $installerPath -ArgumentList "/VERYSILENT" -Wait
    Write-Host "Instalação concluída!"
} else {
    Write-Host "Erro: Instalador não encontrado em $installerPath"
}

Write-Host "Processo concluído!"
