# The chocolatey job in .github/workflows/release.yml replaces __URL64__ and
# __CHECKSUM64__ with the release's windows-x86_64 zip and its sha256 before
# packing. Chocolatey has no arm64 slot; the x64 build runs under emulation on
# Windows 11 on ARM.
$ErrorActionPreference = 'Stop';

$packageName = $env:chocolateyPackageName
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
    packageName    = $packageName
    unzipLocation  = $toolsDir
    url64bit       = '__URL64__'
    checksum64     = '__CHECKSUM64__'
    checksumType64 = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
