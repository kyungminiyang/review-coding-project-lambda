# Test script for receipt-go Lambda function
# Tests S3 upload functionality with review-receipt-uploads bucket

$url = "https://guvajbup2o7srsqgluydpe5tqq0aukrv.lambda-url.us-east-1.on.aws/"

Write-Host "================================" -ForegroundColor Cyan
Write-Host "receipt-go Lambda 함수 테스트" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: 첫 번째 파일 업로드
Write-Host "Test 1: 첫 번째 파일 업로드 (receipt.jpg)" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor Yellow

$testContent1 = "Test receipt image content #1 - $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
$bytes1 = [System.Text.Encoding]::UTF8.GetBytes($testContent1)
$base64_1 = [Convert]::ToBase64String($bytes1)

$body1 = @{
    filename = "receipt.jpg"
    file_content = $base64_1
    content_type = "image/jpeg"
} | ConvertTo-Json

try {
    $response1 = Invoke-RestMethod -Uri $url -Method POST -ContentType "application/json" -Body $body1
    Write-Host "✅ 성공!" -ForegroundColor Green
    $response1 | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ 실패!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "상세: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Start-Sleep -Seconds 2

# Test 2: 같은 이름으로 두 번째 업로드 (중복 파일명 처리 테스트)
Write-Host "Test 2: 중복 파일명 테스트 (receipt.jpg)" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor Yellow

$testContent2 = "Test receipt image content #2 - $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
$bytes2 = [System.Text.Encoding]::UTF8.GetBytes($testContent2)
$base64_2 = [Convert]::ToBase64String($bytes2)

$body2 = @{
    filename = "receipt.jpg"
    file_content = $base64_2
    content_type = "image/jpeg"
} | ConvertTo-Json

try {
    $response2 = Invoke-RestMethod -Uri $url -Method POST -ContentType "application/json" -Body $body2
    Write-Host "✅ 성공!" -ForegroundColor Green
    $response2 | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ 실패!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "상세: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Start-Sleep -Seconds 2

# Test 3: PNG 파일 업로드
Write-Host "Test 3: PNG 파일 업로드 (document-2025.png)" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor Yellow

$testContent3 = "Test PNG document content - $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
$bytes3 = [System.Text.Encoding]::UTF8.GetBytes($testContent3)
$base64_3 = [Convert]::ToBase64String($bytes3)

$body3 = @{
    filename = "document-2025.png"
    file_content = $base64_3
    content_type = "image/png"
} | ConvertTo-Json

try {
    $response3 = Invoke-RestMethod -Uri $url -Method POST -ContentType "application/json" -Body $body3
    Write-Host "✅ 성공!" -ForegroundColor Green
    $response3 | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ 실패!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "상세: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Start-Sleep -Seconds 2

# Test 4: PDF 파일 업로드
Write-Host "Test 4: PDF 파일 업로드 (invoice-november.pdf)" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor Yellow

$testContent4 = "Test PDF content - Invoice for November - $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
$bytes4 = [System.Text.Encoding]::UTF8.GetBytes($testContent4)
$base64_4 = [Convert]::ToBase64String($bytes4)

$body4 = @{
    filename = "invoice-november.pdf"
    file_content = $base64_4
    content_type = "application/pdf"
} | ConvertTo-Json

try {
    $response4 = Invoke-RestMethod -Uri $url -Method POST -ContentType "application/json" -Body $body4
    Write-Host "✅ 성공!" -ForegroundColor Green
    $response4 | ConvertTo-Json -Depth 10
} catch {
    Write-Host "❌ 실패!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host "상세: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "Test Complete" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Verification Checklist:" -ForegroundColor Yellow
Write-Host "  1. All files uploaded to review-receipt-uploads bucket" -ForegroundColor White
Write-Host "  2. Folder structure is YYYY-MM-DD format" -ForegroundColor White
Write-Host "  3. Duplicate filenames saved with different names (timestamp + random)" -ForegroundColor White
Write-Host "  4. S3 URLs returned successfully" -ForegroundColor White

