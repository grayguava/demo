# cli/ui.py
# Terminal output primitives for fscl.

import os
import sys

if os.name == "nt":
    try:
        os.system("chcp 65001 >nul")
    except:
        pass

try:
    sys.stdout.reconfigure(encoding="utf-8")
    sys.stderr.reconfigure(encoding="utf-8")
except:
    pass

RESET  = "\x1b[0m"
BOLD   = "\x1b[1m"
DIM    = "\x1b[2m"

RED    = "\x1b[31m"
GREEN  = "\x1b[32m"
YELLOW = "\x1b[33m"
BLUE   = "\x1b[34m"
MAGENTA= "\x1b[35m"
CYAN   = "\x1b[36m"
WHITE  = "\x1b[37m"
GRAY   = "\x1b[90m"

# Orange accent
O = "\x1b[38;5;208m"   # Orange
S = "\x1b[38;5;112m"   # Seafoam
G = "\x1b[38;5;244m"   # Gray
C = "\x1b[38;5;117m"   # Light blue
Y = "\x1b[38;5;103m"   # Yellow
W = WHITE + BOLD
D = DIM
R = RESET


def br():
    print()


def rule():
    print(G + " " + "\u2500" * 52 + R)


def badge(label, color):
    return f"{color}{BOLD} {label} {R}"


def fail(msg):
    br()
    print(f"{badge('ERR', RED)} {msg}")
    br()
    raise SystemExit(1)


def row(icon, label, value):
    pad   = 12
    label = (label + " " * pad)[:pad]
    print(f"{S}{icon}{R}  {D}{label}{R}  {W}{value}{R}")


def cmd_line(command, desc):
    pad     = 34
    command = (command + " " * pad)[:pad]
    print(f"  {W}{command}{R}{G}{desc}{R}")


def code(msg):
    print(f"  {O}{msg}{R}")


def link(msg):
    print(f"  {O}{msg}{R}")


def ok(msg):
    print(f"  {G}\u2713{R} {msg}")


def info(msg):
    print(f"  {O}{msg}{R}")


def warn(msg):
    print(f"  {Y}Warning:{R} {msg}")