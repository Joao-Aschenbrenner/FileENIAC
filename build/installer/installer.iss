; FileENIAC Installer - Inno Setup Script
; Copyright (C) 2026 ENIAC Systems

#define MyAppName "FileENIAC"
#define MyAppVersion "1.0.0-rc1"
#define MyAppPublisher "ENIAC Systems"
#define MyAppURL "https://github.com/ENIACSystems/FileENIAC"
#define MyAppExeName "eniac.exe"
#define AppIdValue "{E8A1B2C3-D4F5-6789-0ABC-DEF012345678}"

[Setup]
AppId={{#AppIdValue}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
DefaultDirName={localappdata}\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
DisableDirPage=yes
UsePreviousAppDir=yes
OutputDir=.
OutputBaseFilename=FileENIAC_Setup
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=none
UninstallDisplayIcon={app}\icon.ico
ChangesEnvironment=yes
CloseApplications=yes
AppModifyPath="{app}\{#MyAppExeName}"

[Languages]
Name: "brazilianportuguese"; MessagesFile: "compiler:Languages\BrazilianPortuguese.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"

[Messages]
brazilianportuguese.SetupAppTitle=Instalar %1
brazilianportuguese.SetupWindowTitle=%1 - Instalador

[Tasks]
Name: "desktopicon"; Description: "Criar atalho na &Area de Trabalho"; GroupDescription: "Atalhos:"; Flags: checkedonce
Name: "addtopath"; Description: "Adicionar ao &PATH do sistema"; GroupDescription: "Configuracao:"; Flags: checkedonce

[Files]
Source: "..\bin\eniac.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\bin\FileENIAC.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\runtime\WebView2Loader.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "icon.ico"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Parameters: "native --app ""{app}\FileENIAC.exe"""; WorkingDir: "{app}"; IconFilename: "{app}\icon.ico"; Comment: "Iniciar FileENIAC (modo nativo)"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Parameters: "native --app ""{app}\FileENIAC.exe"""; WorkingDir: "{app}"; IconFilename: "{app}\icon.ico"; Comment: "Iniciar FileENIAC (modo nativo)"; Tasks: desktopicon

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Tasks: addtopath; Check: NeedsAddPath

[Run]
Filename: "{app}\{#MyAppExeName}"; Parameters: "native --app ""{app}\FileENIAC.exe"""; Description: "Iniciar {#MyAppName}"; Flags: postinstall nowait skipifsilent shellexec

[Code]
const
  ACTION_INSTALL  = 0;
  ACTION_REPAIR   = 1;
  ACTION_UNINSTALL = 2;

var
  SelectedAction: Integer;

function IsAppInstalled: Boolean;
var
  UninstallString: string;
begin
  Result := RegQueryStringValue(HKEY_CURRENT_USER,
    'Software\Microsoft\Windows\CurrentVersion\Uninstall\{#AppIdValue}_is1',
    'UninstallString', UninstallString);
end;

function GetInstalledAppPath: string;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER,
    'Software\Microsoft\Windows\CurrentVersion\Uninstall\{#AppIdValue}_is1',
    'Inno Setup: App Path', Result) then
    Result := '';
end;

function RunUninstaller: Boolean;
var
  UninstallString: string;
  ResultCode: Integer;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER,
    'Software\Microsoft\Windows\CurrentVersion\Uninstall\{#AppIdValue}_is1',
    'UninstallString', UninstallString) then
  begin
    Result := False;
    Exit;
  end;

  UninstallString := RemoveQuotes(UninstallString);
  Result := Exec(UninstallString, '/SILENT /NORESTART', '', SW_SHOW, ewWaitUntilTerminated, ResultCode);
end;

procedure InitializeWizard;
var
  Form: TSetupForm;
  InfoLabel: TLabel;
  RepairBtn, UninstallBtn, CancelBtn: TNewButton;
  InstalledPath: string;
begin
  SelectedAction := ACTION_INSTALL;

  if not IsAppInstalled then
    Exit;

  InstalledPath := GetInstalledAppPath;

  Form := CreateCustomForm(ScaleX(460), ScaleY(220), False, False);
  try
    Form.Caption := '{#MyAppName} - Instalador';
    Form.CenterOnShow := True;
    Form.BorderStyle := bsDialog;
    Form.BorderIcons := [biSystemMenu];

    InfoLabel := TLabel.Create(Form);
    InfoLabel.Parent := Form;
    InfoLabel.Left := ScaleX(24);
    InfoLabel.Top := ScaleY(20);
    InfoLabel.Width := ScaleX(412);
    InfoLabel.Height := ScaleY(60);
    InfoLabel.WordWrap := True;
    InfoLabel.Font.Size := 9;
    if InstalledPath <> '' then
      InfoLabel.Caption := '{#MyAppName} ja esta instalado em:' + #13#10 + InstalledPath + #13#10 + #13#10 + 'O que voce deseja fazer?'
    else
      InfoLabel.Caption := '{#MyAppName} ja esta instalado neste computador.' + #13#10 + #13#10 + 'O que voce deseja fazer?';

    RepairBtn := TNewButton.Create(Form);
    RepairBtn.Parent := Form;
    RepairBtn.Left := ScaleX(24);
    RepairBtn.Top := ScaleY(120);
    RepairBtn.Width := ScaleX(130);
    RepairBtn.Height := ScaleY(36);
    RepairBtn.Caption := 'Reparar';
    RepairBtn.ModalResult := mrOk;
    RepairBtn.Default := True;

    UninstallBtn := TNewButton.Create(Form);
    UninstallBtn.Parent := Form;
    UninstallBtn.Left := ScaleX(165);
    UninstallBtn.Top := ScaleY(120);
    UninstallBtn.Width := ScaleX(130);
    UninstallBtn.Height := ScaleY(36);
    UninstallBtn.Caption := 'Desinstalar';
    UninstallBtn.ModalResult := mrYes;

    CancelBtn := TNewButton.Create(Form);
    CancelBtn.Parent := Form;
    CancelBtn.Left := ScaleX(306);
    CancelBtn.Top := ScaleY(120);
    CancelBtn.Width := ScaleX(130);
    CancelBtn.Height := ScaleY(36);
    CancelBtn.Caption := 'Sair';
    CancelBtn.ModalResult := mrCancel;

    case Form.ShowModal of
      mrOk:
        begin
          SelectedAction := ACTION_REPAIR;
          WizardForm.Caption := '{#MyAppName} - Reparar';
        end;
      mrYes:
        begin
          SelectedAction := ACTION_UNINSTALL;
        end;
      mrCancel:
        begin
          SelectedAction := ACTION_INSTALL;
          Abort;
        end;
    else
      Abort;
    end;
  finally
    Form.Free;
  end;
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;
end;

function InitializeSetup: Boolean;
begin
  if (SelectedAction = ACTION_UNINSTALL) then
  begin
    RunUninstaller;
    Abort;
    Result := False;
  end
  else
    Result := True;
end;

function NeedsAddPath: Boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER,
    'Environment',
    'Path', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  Result := Pos(LowerCase('{app}'), LowerCase(OrigPath)) = 0;
end;
