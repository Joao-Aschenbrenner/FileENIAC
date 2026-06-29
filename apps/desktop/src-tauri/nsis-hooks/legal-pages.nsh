; SPDX-License-Identifier: MIT
; FileENIAC NSIS Legal Pages Hook
; This file is included by the custom NSIS template to add legal acceptance pages
; Run scripts/patch-installer-nsis.ps1 after building to inject these pages

!include "MUI2.nsh"
!include "nsDialogs.nsh"
!include "LogicLib.nsh"

; Checkbox states
Var LegalTermsAccepted
Var LegalPrivacyAccepted

; ============================================================
; TERMS OF USE PAGE
; ============================================================
Function legalTermsPage
  nsDialogs::Create 1018
  Pop $R0

  ${IfThen} $R0 == 0 ${|} Abort ${|}

  ${NSD_CreateLabel} 0 0 100% 24u "Terms of Use"
  Pop $R0
  SendMessage $R0 ${WM_SETFONT} [SYSFONT] 0

  ${NSD_CreateBrowser} 0 30u 100% 200u "file://$PLUGINSDIR/../../docs/legal/TERMS_OF_USE.md"
  Pop $R0

  ${NSD_CreateCheckbox} 0 240u 100% 10u "I have read and agree to the Terms of Use"
  Pop $R1
  ${NSD_OnClick} $R1 legalTermsCheckbox

  GetFunctionAddress $0 legalTermsPageLeave
  nsDialogs::Show
FunctionEnd

Function legalTermsCheckbox
  Pop $R0
  ${NSD_GetState} $R0 $R1
  ${If} $R1 == ${BST_CHECKED}
    StrCpy $LegalTermsAccepted 1
  ${Else}
    StrCpy $LegalTermsAccepted 0
  ${EndIf}
FunctionEnd

Function legalTermsPageLeave
  ${If} $LegalTermsAccepted == 0
    MessageBox MB_ICONEXCLAMATION "You must accept the Terms of Use to continue installation."
    Abort
  ${EndIf}
FunctionEnd

; ============================================================
; PRIVACY POLICY PAGE
; ============================================================
Function legalPrivacyPage
  nsDialogs::Create 1018
  Pop $R0

  ${IfThen} $R0 == 0 ${|} Abort ${|}

  ${NSD_CreateLabel} 0 0 100% 24u "Privacy Policy"
  Pop $R0
  SendMessage $R0 ${WM_SETFONT} [SYSFONT] 0

  ${NSD_CreateBrowser} 0 30u 100% 200u "file://$PLUGINSDIR/../../docs/legal/PRIVACY_POLICY.md"
  Pop $R0

  ${NSD_CreateCheckbox} 0 240u 100% 10u "I have read and agree to the Privacy Policy"
  Pop $R1
  ${NSD_OnClick} $R1 legalPrivacyCheckbox

  GetFunctionAddress $0 legalPrivacyPageLeave
  nsDialogs::Show
FunctionEnd

Function legalPrivacyCheckbox
  Pop $R0
  ${NSD_GetState} $R0 $R1
  ${If} $R1 == ${BST_CHECKED}
    StrCpy $LegalPrivacyAccepted 1
  ${Else}
    StrCpy $LegalPrivacyAccepted 0
  ${EndIf}
FunctionEnd

Function legalPrivacyPageLeave
  ${If} $LegalPrivacyAccepted == 0
    MessageBox MB_ICONEXCLAMATION "You must accept the Privacy Policy to continue installation."
    Abort
  ${EndIf}
FunctionEnd