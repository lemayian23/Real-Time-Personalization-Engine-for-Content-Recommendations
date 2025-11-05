
#region conda initialize
# !! Contents within this block are managed by 'conda init' !!
If (Test-Path "D:\ANACONDA3\Scripts\conda.exe") {
    (& "D:\ANACONDA3\Scripts\conda.exe" "shell.powershell" "hook") | Out-String | ?{$_} | Invoke-Expression
}
#endregion

