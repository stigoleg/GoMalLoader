import argparse
from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes

parser = argparse.ArgumentParser(description="Encrypt shellcode/DLL with AES-CBC.")
parser.add_argument("--in", required=True, help="Input binary file (raw shellcode or DLL)")
parser.add_argument("--out", required=True, help="Output file for AES-encrypted payload")
parser.add_argument("--key", required=True, help="16-byte AES key (string)")

args = parser.parse_args()

key = args.key.encode("utf-8")
assert len(key) == 16, "AES key must be exactly 16 bytes"

with open(args.in, "rb") as f:
    data = f.read()

pad = AES.block_size - len(data) % AES.block_size
data += bytes([pad] * pad)

iv = get_random_bytes(16)
cipher = AES.new(key, AES.MODE_CBC, iv)
encrypted = iv + cipher.encrypt(data)

with open(args.out, "wb") as f:
    f.write(encrypted)

print(f"[+] Encrypted payload saved to {args.out}")
