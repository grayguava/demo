!include "MUI2.nsh"
!include "EnVar.nsh"

Name "formseal-fetch"
OutFile "formseal-fetch-setup.exe"
InstallDir "$PROGRAMFILES64\formseal-fetch"
InstallDirRegKey HKLM "Software\formseal-fetch" "InstallDir"
RequestExecutionLevel admin

!define MUI_ABORTWARNING
!define MUI_WELCOMEPAGE_TITLE "Welcome to formseal-fetch Setup"
!define MUI_WELCOMEPAGE_TEXT "This will install formseal-fetch (fsf) on your computer.$\r$\n$\r$\nAfter installation, open any terminal and type 'fsf' to get started."
!define MUI_FINISHPAGE_TITLE "Installation Complete"
!define MUI_FINISHPAGE_TEXT "formseal-fetch has been installed.$\r$\n$\r$\nOpen a new terminal window and type 'fsf status' to verify."

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "..\..\LICENSE.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "English"

Section "Install"
    SetOutPath "$INSTDIR"
    File "..\..\dist\fsf.exe"

    EnVar::SetHKLM
    EnVar::AddValue "PATH" "$INSTDIR"

    WriteUninstaller "$INSTDIR\uninstall.exe"

    WriteRegStr HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "DisplayName" "formseal-fetch"
    WriteRegStr HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "UninstallString" '"$INSTDIR\uninstall.exe"'
    WriteRegStr HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "QuietUninstallString" '"$INSTDIR\uninstall.exe" /S'
    WriteRegStr HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "InstallLocation" "$INSTDIR"
    WriteRegStr HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "DisplayVersion" "INSTALLER_VERSION"
    WriteRegStr HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "Publisher" "YourName"
    WriteRegDWORD HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "NoModify" 1
    WriteRegDWORD HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch" \
        "NoRepair" 1
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\fsf.exe"
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"

    EnVar::SetHKLM
    EnVar::DeleteValue "PATH" "$INSTDIR"

    DeleteRegKey HKLM \
        "Software\Microsoft\Windows\CurrentVersion\Uninstall\formseal-fetch"
    DeleteRegKey HKLM "Software\formseal-fetch"
SectionEnd