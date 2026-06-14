<#
Script para desinstalar e reinstalar o FileENIAC

1. Desinstala completamente o FileENIAC
2. Remove todos os arquivos residuais
3. Instala a versão mais recente
#>

# Desinstalar completamente
Write-Host "Desinstalando FileENIAC..."

# Remover do registro
if (Test-Path "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\FileENIAC") {
    Remove-Item "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\FileENIAC" -Recurse -Force
}

# Remover arquivos do programa
$programFiles = "$env:ProgramFiles\FileENIAC"
if (Test-Path $programFiles) {
    Remove-Item $programFiles -Recurse -Force
}

# Remover arquivos do usuário
$localAppData = "$env:LOCALAPPDATA\FileENIAC"
if (Test-Path $localAppData) {
    Remove-Item $localAppData -Recurse -Force
}

# Remover atalhos
$startMenu = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\FileENIAC"
if (Test-Path $startMenu) {
    Remove-Item $startMenu -Recurse -Force
}

# Remover arquivos temporários
$tempFiles = "$env:TEMP\FileENIAC*"
if (Test-Path $tempFiles) {
    Remove-Item $tempFiles -Force
}

# Instalar versão mais recente
Write-Host "Instalando FileENIAC..."

# Caminho para o instalador (substitua pelo caminho real)
$installerPath = "C:\Users\USUARIO\Documents\GitWrkspc\FileENIAC\build\installer\FileENIAC_Setup.exe"

if (Test-Path $installerPath) {
    Start-Process -FilePath $installerPath -ArgumentList "/VERYSILENT" -Wait
    Write-Host "Instalação concluída!"
} else {
    Write-Host "Erro: Instalador não encontrado em $installerPath"
}

Write-Host "Processo concluído!"
