import sys, traceback
from pathlib import Path
p = Path(r"E:\1tebeFile\Development\E-commerce\.github\workflows\notify-telegram.yml")
try:
    import yaml
except Exception as e:
    print("MISSING_PYYAML", e)
    sys.exit(2)
try:
    with p.open(encoding='utf-8') as f:
        yaml.safe_load(f)
    print("YAML_OK")
    sys.exit(0)
except Exception:
    print("PARSE_ERROR")
    traceback.print_exc()
    sys.exit(3)
