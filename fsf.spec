from PyInstaller.building.build_main import Analysis, PYZ, EXE

a = Analysis(
    ['cli/fsf.py'],
    pathex=['.'],
    binaries=[],
    datas=[
        ('version.txt', '.'),
    ],
    hiddenimports=[
        'keyring.backends.Windows',
        'keyring.backends.macOS',
        'keyring.backends.SecretService',
        'keyring.backends.kwallet',
        'keyring.backends.fail',
        'keyring.backends.null',
        'cli.commands',
        'cli.commands.config',
        'cli.commands.fetch',
        'cli.commands.providers',
        'cli.commands.setup',
        'cli.providers',
        'cli.providers.cloudflare',
        'cli.providers.cloudflare.account',
        'cli.providers.cloudflare.storage',
        'cli.security',
        'cli.security.tokens',
    ],
    hookspath=[],
    noarchive=False,
)

pyz = PYZ(a.pure)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.datas,
    name='fsf',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=True,
    onefile=True,
)