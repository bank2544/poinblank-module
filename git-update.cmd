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

REM ดึงการเปลี่ยนแปลงจาก remote ก่อน
git pull --rebase

REM เพิ่มไฟล์ทั้งหมด ยกเว้น .gitignore
for /f "delims=" %%f in ('git ls-files --modified --others --exclude-standard') do (
    if /i not "%%f"==".gitignore" (
        git add "%%f"
    )
)

REM Commit พร้อมข้อความ
git commit -m "!MSG!"

REM ดึงชื่อ branch ปัจจุบัน
for /f "delims=" %%b in ('git rev-parse --abbrev-ref HEAD') do set BRANCH=%%b

REM Push ขึ้น remote
git push -u origin !BRANCH!

echo ✅ Push เสร็จแล้วที่ branch "!BRANCH!"

REM --- FORCE UPDATE TAG v0.3.0 ---
git tag -d v0.3.0
git tag v0.3.0
git push origin :refs/tags/v0.3.0
git push origin v0.3.0

echo ✅ อัปเดต tag v0.3.0 สำเร็จ
pause
