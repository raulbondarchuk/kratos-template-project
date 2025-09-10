[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Title,
  [Parameter(Mandatory = $true)] [string]$Desc,
  [string]$ConfigPath = "./configs/config.yaml"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Get-AppVersionFromYaml {
  param([string]$Path)
  if (-not (Test-Path $Path)) { throw "Config not found: $Path" }
  $lines = Get-Content -LiteralPath $Path -Encoding UTF8
  $inApp = $false
  $base  = $null
  foreach ($line in $lines) {
    if ($line -match '^\s*app\s*:') { $inApp = $true; continue }
    if ($inApp -and $line -match '^\S') { break }
    if ($inApp -and $line -match '^\s*version\s*:\s*(.+?)\s*$') {
      $base = $Matches[1].Trim("'`"")
      break
    }
  }
  if (-not $base) { throw "app.version not found in $Path" }
  return $base
}

function Test-BaseVersion {
  param([string]$Base)
  if ($Base -notmatch '^v\d+$') {
    throw "app.version must look like v1, v2 (got '$Base')"
  }
}

function Test-BranchPolicy {
  $branch = (& git rev-parse --abbrev-ref HEAD).Trim()
  if ($LASTEXITCODE -ne 0) { throw "Not a git repository or no HEAD." }
  if ($branch -eq 'HEAD') { throw "Detached HEAD state - checkout a branch first." }

  $blocked = @('master', 'main')
  if ($blocked -contains $branch.ToLower()) {
    throw "Pushing from '$branch' is forbidden by policy."
  }
  return $branch
}

function Get-NextPatchVersion {
  param([string]$Base)
  $pattern = '^{0}\.(\d+)$' -f ([Regex]::Escape($Base))
  # Look only at remote tags so everyone shares the same source of truth
  $lines = & git ls-remote --tags origin "$Base.*" 2>$null
  $max = 0
  foreach ($ln in $lines) {
    $parts = $ln -split "`t"
    if ($parts.Count -lt 2) { continue }
    $ref = $parts[1]
    if ($ref.EndsWith("^{}")) { $ref = $ref.Substring(0, $ref.Length - 3) }
    if ($ref -match 'refs/tags/(.+)$') {
      $tag = $Matches[1]
      if ($tag -match $pattern) {
        $n = [int]$Matches[1]
        if ($n -gt $max) { $max = $n }
      }
    }
  }
  if ($max -eq 0) { return 1 } else { return ($max + 1) }
}

function Test-UniqueTag {
  param([string]$Candidate)
  & git rev-parse -q --verify "refs/tags/$Candidate" *> $null
  return ($LASTEXITCODE -ne 0)
}

function Push-CurrentBranch {
  param([string]$Branch)
  $remotes = (& git remote) -split '\r?\n' | Where-Object { $_ }
  if ($remotes -notcontains 'origin') { throw "Remote 'origin' is not configured." }

  & git rev-parse --abbrev-ref --symbolic-full-name '@{u}' *> $null
  if ($LASTEXITCODE -ne 0) {
    & git push -u origin "HEAD:$Branch"
    if ($LASTEXITCODE -ne 0) { throw "git push (set upstream) failed." }
  } else {
    & git push
    if ($LASTEXITCODE -ne 0) { throw "git push failed." }
  }
}

function Push-TagWithRetry {
  param(
    [string]$Base,
    [string]$CurrentTag,
    [string]$Title,
    [string]$Desc,
    [int]$MaxAttempts = 5
  )
  for ($i = 1; $i -le $MaxAttempts; $i++) {
    & git push origin "$CurrentTag"
    if ($LASTEXITCODE -eq 0) {
      return $CurrentTag
    }

    # Collision: someone pushed the same tag name
    & git fetch --tags --quiet | Out-Null
    $next = Get-NextPatchVersion -Base $Base
    $newTag = "$Base.$next"

    if ($newTag -eq $CurrentTag) {
      throw "git push tag failed and a new unique tag could not be determined."
    }

    # Retag locally: delete old tag and create a new one pointing to HEAD
    & git tag -d "$CurrentTag" *> $null
    if ($LASTEXITCODE -ne 0) { throw "failed to delete local tag $CurrentTag" }

    # Tag annotation equals full commit message (Title + Desc)
    & git tag -a "$newTag" -m "$Title" -m "$Desc"
    if ($LASTEXITCODE -ne 0) { throw "failed to create local tag $newTag" }

    $CurrentTag = $newTag
  }
  throw "Failed to push tag after $MaxAttempts attempts."
}

# --- branch policy
$branch = Test-BranchPolicy
Write-Host "Current branch: $branch"

# --- base version from YAML
$base = Get-AppVersionFromYaml -Path $ConfigPath
Test-BaseVersion $base
Write-Host "Base version from YAML: $base"

# --- next patch
$next = Get-NextPatchVersion -Base $base
$version = "$base.$next"

# ensure unique locally (handles manual local tags)
$guard = 0
while (-not (Test-UniqueTag -Candidate $version)) {
  $next++
  $version = "$base.$next"
  $guard++
  if ($guard -gt 100) { throw "Too many existing tags for base $base" }
}
Write-Host "Next version: $version"

# --- anything to commit?
$changes = & git status --porcelain
if (-not $changes) {
  Write-Host "No changes to commit. Nothing to do."
  exit 0
}

# --- commit + local tag
& git add -A
& git commit -m "$Title" -m "$Desc"
if ($LASTEXITCODE -ne 0) { throw "git commit failed." }

# Tag annotation equals full commit message (Title + Desc)
& git tag -a "$version" -m "$Title" -m "$Desc"
if ($LASTEXITCODE -ne 0) { throw "git tag failed." }

# --- push branch, then push tag with collision retry
Push-CurrentBranch -Branch $branch
$finalTag = Push-TagWithRetry -Base $base -CurrentTag $version -Title $Title -Desc $Desc

Write-Host "Committed and pushed. Tag: $finalTag"