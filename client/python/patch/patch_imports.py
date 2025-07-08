import os
import re

GENERATED_ROOT = "src/sdk/gen"
FROM_PATTERN = r'^from storage\.v2\b'
REPLACEMENT = 'from sdk.gen.storage.v2'

def patch_file(path):
    with open(path, "r", encoding="utf-8") as f:
        lines = f.readlines()

    changed = False
    with open(path, "w", encoding="utf-8") as f:
        for line in lines:
            if re.match(FROM_PATTERN, line):
                line = re.sub(FROM_PATTERN, REPLACEMENT, line)
                changed = True
            f.write(line)

    if changed:
        print(f"Patched: {path}")

def patch_all():
    for root, _, files in os.walk(GENERATED_ROOT):
        for file in files:
            if file.endswith(".py"):
                patch_file(os.path.join(root, file))

if __name__ == "__main__":
    patch_all()
