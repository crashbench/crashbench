Write-Host "Installing CrashBench (The Open Chaos & Safety Benchmark for AI Agents)..." -ForegroundColor Cyan

if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "Compiling via go install..." -ForegroundColor Yellow
    go install github.com/crashbench/crashbench/cmd/crashbench@latest
    Write-Host "✅ CrashBench installed successfully! Run 'crashbench help' to get started." -ForegroundColor Green
} else {
    Write-Host "Go compiler not found. Please install Go or download precompiled binaries from https://github.com/crashbench/crashbench/releases" -ForegroundColor Red
}
