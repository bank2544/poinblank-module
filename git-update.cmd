@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM ตั้งข้อความ commit
if "%~1"=="" (
    set MSG=อัปเดตอัตโนมัติ %DATE% %TIME%
) else (
    set MSG=%*
)

REM ตรวจสอบว่าเป็น Git repo ไหม
if not exist ".git" (
    echo ❌ ไม่พบ Git repository
    exit /b 1
)

git pull --rebase
git add .
git commit -m "!MSG!"

for /f "delims=" %%b in ('git rev-parse --abbrev-ref HEAD') do set BRANCH=%%b

git push -u origin !BRANCH!

echo ✅ Push เสร็จแล้วที่ branch "!BRANCH!"
pause
