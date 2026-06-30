; Scenaria Windows installer — Inno Setup (adds scenaria CLI to PATH).
; Build: scripts/build-installer.ps1

#define MyAppName "Scenaria"
#define MyAppPublisher "Scenaria"
#define MyAppURL "https://github.com/bafgion/Scenaria-Golang-"
#define MyAppExe "scenaria-gui.exe"
#define MyCLIExe "scenaria.exe"

[Setup]
AppId={{A4B8C2D1-5E6F-4A7B-9C0D-1E2F3A4B5C6D}
AppName={#MyAppName}
AppVersion={#AppVersion}
AppVerName={#MyAppName} {#AppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}/releases
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
OutputDir=..\dist
OutputBaseFilename=Scenaria-Setup
Compression=lzma2/ultra64
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
ChangesEnvironment=yes
SetupIconFile=..\assets\branding\app.ico
UninstallDisplayIcon={app}\{#MyAppExe}
LicenseFile=
InfoBeforeFile=
MinVersion=10.0

[Languages]
Name: "russian"; MessagesFile: "compiler:Languages\Russian.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked
Name: "addtopath"; Description: "Добавить scenaria CLI в системный PATH"; GroupDescription: "Дополнительно:"; Flags: checkedonce

[Files]
Source: "..\dist\Scenaria\{#MyCLIExe}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\Scenaria\{#MyAppExe}"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\Scenaria\browsers\*"; DestDir: "{app}\browsers"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\dist\Scenaria\examples\*"; DestDir: "{app}\examples"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "..\dist\Scenaria\version.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\Scenaria\README-PORTABLE.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\Scenaria\Start-GUI.bat"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\dist\Scenaria\scenaria-cli.bat"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#MyAppName} IDE"; Filename: "{app}\{#MyAppExe}"; WorkingDir: "{app}"
Name: "{group}\{#MyAppName} CLI"; Filename: "{app}\{#MyCLIExe}"; WorkingDir: "{app}"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExe}"; Tasks: desktopicon; WorkingDir: "{app}"

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Tasks: addtopath; Check: NeedsAddPath(ExpandConstant('{app}'))

[Run]
Filename: "{app}\{#MyAppExe}"; Description: "Запустить {#MyAppName} IDE"; Flags: nowait postinstall skipifsilent unchecked

[Code]
function NeedsAddPath(Param: string): Boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_LOCAL_MACHINE,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

procedure RemoveDirFromPath(const Dir: string);
var
  Path: string;
  P: Integer;
begin
  if not RegQueryStringValue(HKEY_LOCAL_MACHINE,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', Path)
  then
    exit;

  P := Pos(';' + Dir, Path);
  if P > 0 then
    Delete(Path, P, Length(Dir) + 1)
  else if CompareText(Copy(Path, 1, Length(Dir)), Dir) = 0 then begin
    if (Length(Path) > Length(Dir)) and (Path[Length(Dir) + 1] = ';') then
      Delete(Path, 1, Length(Dir) + 1)
    else
      Path := '';
  end;

  RegWriteExpandStringValue(HKEY_LOCAL_MACHINE,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', Path);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RemoveDirFromPath(ExpandConstant('{app}'));
end;
