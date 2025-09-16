[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)] [string]$Title,
  [Parameter(Mandatory = $true)] [string]$Desc,
  [string]$ConfigPath = './configs/config.yaml'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# Логи
. "$PSScriptRoot/make/utils.ps1"

function Get-AppVersionFromYaml {
  param([string]$Path)
  if (-not (Test-Path -LiteralPath $Path)) { throw ('Config not found: {0}' -f $Path) }
  $lines = Get-Content -LiteralPath $Path -Encoding UTF8
  $inApp = $false
  $base  = $null
  foreach ($line in $lines) {
    if ($line -match '^\s*app\s*:') { $inApp = $true; continue }
    if ($inApp -and $line -match '^\S') { break }
    if ($inApp -and $line -match '^\s*version\s*:\s*(.+?)\s*$') {
      $base = $Matches[1].Trim("'",'"')
      break
    }
  }
  if (-not $base) { throw ('app.version not found in {0}' -f $Path) }
  return $base
}

function Test-BaseVersion {
  param([string]$Base)
  if ($Base -notmatch '^v\d+$') {
    throw ('app.version must look like v1, v2 (got ''{0}'')' -f $Base)
  }
}

function Test-BranchPolicy {
  $branch = (& git rev-parse --abbrev-ref HEAD).Trim()
  if ($LASTEXITCODE -ne 0) { throw 'Not a git repository or no HEAD.' }
  if ($branch -eq 'HEAD') { throw 'Detached HEAD state - checkout a branch first.' }

  $blocked = @('master','main')
  if ($blocked -contains $branch.ToLower()) {
    throw ("Pushing from '{0}' is forbidden by policy." -f $branch)
  }
  return $branch
}

function Get-NextPatchVersion {
  param([string]$Base)
  $pattern = '^{0}\.(\d+)$' -f ([Regex]::Escape($Base))
  $lines = & git ls-remote --tags origin "$Base.*" 2>$null
  $max = 0
  foreach ($ln in $lines) {
    $parts = $ln -split "`t"
    if ($parts.Count -lt 2) { continue }
    $ref = $parts[1]
    if ($ref.EndsWith('^{}')) { $ref = $ref.Substring(0, $ref.Length - 3) }
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
  & git rev-parse -q --verify ("refs/tags/{0}" -f $Candidate) *> $null
  return ($LASTEXITCODE -ne 0)
}

function Push-CurrentBranch {
  param([string]$Branch)
  $remotes = (& git remote) -split '\r?\n' | Where-Object { $_ }
  if ($remotes -notcontains 'origin') { throw "Remote 'origin' is not configured." }

  & git rev-parse --abbrev-ref --symbolic-full-name '@{u}' *> $null
  if ($LASTEXITCODE -ne 0) {
    & git push -u origin ("HEAD:{0}" -f $Branch)
    if ($LASTEXITCODE -ne 0) { throw 'git push (set upstream) failed.' }
  } else {
    & git push
    if ($LASTEXITCODE -ne 0) { throw 'git push failed.' }
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
    & git push origin $CurrentTag
    if ($LASTEXITCODE -eq 0) { return $CurrentTag }

    & git fetch --tags --quiet | Out-Null
    $next   = Get-NextPatchVersion -Base $Base
    $newTag = ('{0}.{1}' -f $Base, $next)

    if ($newTag -eq $CurrentTag) {
      throw 'git push tag failed and a new unique tag could not be determined.'
    }

    & git tag -d $CurrentTag *> $null
    if ($LASTEXITCODE -ne 0) { throw ('failed to delete local tag {0}' -f $CurrentTag) }

    & git tag -a $newTag -m $Title -m $Desc
    if ($LASTEXITCODE -ne 0) { throw ('failed to create local tag {0}' -f $newTag) }

    $CurrentTag = $newTag
  }
  throw 'Failed to push tag after max attempts.'
}

# ------------------ Flow ------------------

Show-Step 'Release: commit & tag'

Show-Info 'Validating branch policy'
$branch = Test-BranchPolicy
Show-Info ('Current branch: {0}' -f $branch)

Show-Info 'Reading base version from YAML'
$base = Get-AppVersionFromYaml -Path $ConfigPath
Test-BaseVersion $base
Show-Info ('Base version (app.version): {0}' -f $base)

Show-Info 'Calculating next patch from remote tags'
$next    = Get-NextPatchVersion -Base $base
$version = ('{0}.{1}' -f $base, $next)

# ensure unique locally too
$guard = 0
while (-not (Test-UniqueTag -Candidate $version)) {
  $next++
  $version = ('{0}.{1}' -f $base, $next)
  $guard++
  if ($guard -gt 100) { throw ('Too many existing tags for base {0}' -f $base) }
}
Show-Info ('Next version candidate: {0}' -f $version)

Show-Info 'Checking working tree'
$changes = & git status --porcelain
if (-not $changes) {
  Show-Info 'No changes to commit — nothing to do'
  Show-OK  'Release skipped'
  exit 0
}

Show-Info 'Committing changes'
& git add -A
& git commit -m $Title -m $Desc
if ($LASTEXITCODE -ne 0) { throw 'git commit failed.' }
Show-OK 'Commit created'

Show-Info ('Creating local tag {0}' -f $version)
& git tag -a $version -m $Title -m $Desc
if ($LASTEXITCODE -ne 0) { throw 'git tag failed.' }
Show-OK ('Tag created: {0}' -f $version)

Show-Step ('Pushing current branch ({0})' -f $branch)
Push-CurrentBranch -Branch $branch
Show-OK 'Branch pushed'

Show-Step 'Pushing tag (with collision retry if needed)'
$finalTag = Push-TagWithRetry -Base $base -CurrentTag $version -Title $Title -Desc $Desc
Show-OK ('Tag pushed: {0}' -f $finalTag)

Show-OK ('Committed and pushed. Tag: {0}' -f $finalTag)
