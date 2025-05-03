@echo off
setlocal enabledelayedexpansion

REM ถ้าผู้ใช้ใส่ข้อความ commit เป็น argument
if "%~1"=="" (
    set MSG=อัปเดตอัตโนมัติ %DATE% %TIME%
) else (
    set MSG=%*
)

REM ตรวจสอบว่าอยู่ใน Git repo หรือไม่
if not exist ".git" (
    echo ❌ ไม่พบ Git repository ในโฟลเดอร์นี้
    exit /b 1
)

REM ดึงอัปเดตจาก remote ก่อน
git pull --rebase

REM เพิ่มไฟล์ทั้งหมด
git add .

REM commit พร้อมข้อความ
git commit -m "!MSG!"

REM หา branch ปัจจุบัน
for /f "delims=" %%b in ('git rev-parse --abbrev-ref HEAD') do set BRANCH=%%b

REM push ขึ้น remote
git push -u origin !BRANCH!

echo ✅ Push สำเร็จไปยัง branch "!BRANCH!"
pause
