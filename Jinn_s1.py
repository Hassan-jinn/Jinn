#!/data/data/com.termux/files/usr/bin/python3
# -*- coding: utf-8 -*-
import os
import sys
import time
import base64
import pycurl
from io import BytesIO
import subprocess
import psutil
import fcntl

# Security Constants
GITHUB_TOKEN = "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"
APPROVAL_URL = "https://github.com/1-NALLA/Jinn_App/blob/main/App.txt"
SCRIPT_URL = "https://github.com/1-NALLA/File1/raw/main/JINN8_enc.py"
LOCK_FILE = "/data/data/com.termux/files/usr/tmp/JINN8_enc.py.lock"

# Colors
R = "\033[1;31m"  # Red
G = "\033[1;32m"  # Green
Y = "\033[1;33m"  # Yellow
B = "\033[1;34m"  # Blue
P = "\033[1;35m"  # Purple
C = "\033[1;36m"  # Cyan
W = "\033[1;37m"  # White
A4 = "\033[1;34m"  # Light Blue
G4 = "\033[1;32m"  # Light Green

def clear():
    os.system('clear')

def is_termux():
    return "com.termux" in os.getcwd()

def download_script():
    if not os.path.exists("JINN8_enc.py"):
        print(f"{Y}[!] Downloading JINN Script...{W}")
        os.system(f"curl -sL {SCRIPT_URL} -o JINN8_enc.py")
        os.system("chmod 777 JINN8_enc.py")
        if not os.path.exists("JINN8_enc.py"):
            print(f"{R}[X] Failed to download script!{W}")
            sys.exit(1)

def check_approval(key):
    try:
        buffer = BytesIO()
        c = pycurl.Curl()
        c.setopt(c.URL, APPROVAL_URL)
        c.setopt(c.HTTPHEADER, [f"Authorization: token {GITHUB_TOKEN}"])
        c.setopt(c.WRITEDATA, buffer)
        c.perform()
        response = buffer.getvalue().decode()
        return key in response
    except Exception as e:
        print(f"{R}[!] Connection Error: {str(e)}{W}")
        return False
    finally:
        c.close()

def kill_suspicious_apps():
    for proc in psutil.process_iter(['name']):
        if proc.info['name'] and any(x in proc.info['name'].lower() for x in ['httpcanary', 'wireshark', 'fiddler']):
            try:
                proc.kill()
            except:
                pass

def install_packages():
    try:
        import pycurl, psutil
    except:
        print(f"{Y}[!] Installing required packages...{W}")
        os.system("pip install pycurl psutil --quiet >/dev/null 2>&1")

def create_lock():
    try:
        lock = open(LOCK_FILE, "w")
        fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
        return True
    except:
        print(f"{R}[!] Script already running! Exiting...{W}")
        sys.exit(1)

def show_payment_options(key):
    clear()
    print('''\033[1;94m   ╔════════════════════════════════════════════╗
\033[1;97m   ║\033[1;93m\033[1;37m\033[1;41m                 H A S S A N                \033[0m\033[1;37m\033[1;97m║
\033[1;97m   ║                         \033[1;97m                   ║
\033[1;97m   ║      \033[1;96m     ██╗██╗███╗   ██╗███╗   ██╗    \033[1;97m   ║
\033[1;97m   ║      \033[1;96m     ██║██║████╗  ██║████╗  ██║     \033[1;97m  ║
\033[1;97m   ║           ██║██║██╔██╗ ██║██╔██╗ ██║      \033[1;97m ║
\033[1;97m   ║      \033[1;96m██   ██║██║██║╚██╗██║██║╚██╗██║   \033[1;97m    ║
\033[1;97m   ║\033[1;96m      ╚█████╔╝██║██║ ╚████║██║ ╚████║  \033[1;97m     ║
\033[1;97m   ║       ╚════╝ ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝ \033[1;97m      ║
\033[1;96m   ╚════════════════════════════════════════════╝
''')
    print('\033[1;97m   ══════════════════════════════════════════════')
    print(f"\033[1;31m  YOUR KEY NOT APPROVED CONTACT ADMIN")
    print(f"\033[1;34m  7 DAY's APPROVE {A4}RS 50")
    print(f"\033[1;34m  15 DAY's APPROVE {A4}RS 100")
    print(f"\033[1;34m  30 DAY's APPROVE {A4}RS 200")
    print("\033[1;33m  OTHER COUNTRY")
    print(f"\033[1;34m  15 DAY's APPROVE {A4}1$")
    print(f"\033[1;34m  30 DAY's APPROVE {A4}2$")
    print('\033[1;97m   ══════════════════════════════════════════════')    
    print("\33[37;33m\t WELCOME TO JINN TOOL AND ENJOY \33[0;m")
    print('\033[1;97m   ══════════════════════════════════════════════')
    print(f"\033[1;33m Naya Pay IBAN : {A4} PK03NAYA1234503276129154")
    print(f"\033[1;33m JAZZ CASH  : {A4} 0326 6189817 : M HASSAN")
    print(f"\033[1;33m EASYPAISA  : {A4} 0318 9713740 : M HASSAN")
    print(f"\033[1;31m Note : {G4}SEND PAYMENT PROOF ON WHATSAPP")
    print('\033[1;97m   ══════════════════════════════════════════════')
    print(f" \n\033[1;34m Your Login Key is  :{G4} "+key)
    print(f"\n\033[1;33m [1] CONTACT WITH ME ON WHATSAPP")
    print(f"\033[1;33m [2] CONTACT WITH ME ON FACEBOOK")
    
    adi = input(f" \033[1;32m[•] CHOICE : ")
    if adi in ['1','01']:
        handle_whatsapp_contact(key)
    else:
        os.system('xdg-open https://www.facebook.com/profile.php?id=1623021375&mibextid=ZbWKwL')

def handle_whatsapp_contact(key):
    nm = input(f"\033[1;32m ENTER YOUR NAME : ")
    wp = input(f"\033[1;32m ENTER YOUR WHATSAPP NUMBER :")
    url_wa = "https://api.whatsapp.com/send?phone=+923189713740&text="
    tk = (f'Hello%20Sir%20!%20Please%20Approve%20My%20Tool%20Login%20Key%20:%20My%20Name%20is%20'
          f'{nm}%20and%20whatsapp%20number%20is%20%20:{wp}%20Here%20my%20key%20\n{key}')
    subprocess.run(["am", "start", "--user", "0", "-a", "android.intent.action.VIEW", "-d", url_wa+tk])
    time.sleep(2)

def run_jinn8_script():
    try:
        # First try direct execution
        os.system("python JINN8_enc.py")
    except:
        # Fallback method
        os.system("python JINN8_enc.py")

def key():
    clear()
    uid = str(os.geteuid())
    user = str(os.getlogin()).replace('u0_a', '')
    key = f"JINN_{uid}{user}"
    
    print(f" {G}TOOL IS PAID YOU NEED APPROVAL{W}")
    
    if check_approval(key):
        print(f"{G}[✓] TOOL LOGIN SUCCESSFULLY{W}")
        download_script()
        install_packages()
        create_lock()
        kill_suspicious_apps()
        run_jinn8_script()
    else:
        show_payment_options(key)

if __name__ == "__main__":
    key()