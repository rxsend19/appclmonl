# https://stackoverflow.com/a/66512832
Set-Item -Path env:CGO_ENABLED -Value 1;
Write-Output "Build with CGO_ENABLED:$env:CGO_ENABLED";
Set-Location app;
$out_file_name = "appmoncl.exe"
Try { Remove-Item $out_file_name -ErrorAction Ignore; }Catch {}
Write-Output "Deleted old $out_file_name";
(go build -o $out_file_name main)
Write-Output "builded $out_file_name" 
Start-Process -FilePath ((Get-Location).Path + "\" + $out_file_name)
