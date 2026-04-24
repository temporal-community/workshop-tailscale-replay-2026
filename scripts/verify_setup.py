import os
import socket
import sys

PASS = "\033[92m PASS \033[0m"
FAIL = "\033[91m FAIL \033[0m"
WARN = "\033[93m WARN \033[0m"


def check_python_version() -> bool:
    version = sys.version_info
    ok = version >= (3, 13)
    status = PASS if ok else FAIL
    print(f"  [{status}] Python >= 3.13 (found {version.major}.{version.minor}.{version.micro})")
    return ok


def check_import(module: str) -> bool:
    try:
        __import__(module)
        print(f"  [{PASS}] import {module}")
        return True
    except ImportError:
        print(f"  [{FAIL}] import {module} — not installed, run 'uv sync'")
        return False


def check_temporal_connection() -> bool:
    address = os.getenv("TEMPORAL_ADDRESS", "temporal-dev:7233")
    host, _, port = address.rpartition(":")
    if not host:
        host = address
        port = "7233"

    try:
        sock = socket.create_connection((host, int(port)), timeout=5)
        sock.close()
        print(f"  [{PASS}] Temporal server at {address}")
        return True
    except (socket.timeout, ConnectionRefusedError, OSError) as e:
        print(f"  [{FAIL}] Temporal server at {address} — {e}")
        print(f"          Is Tailscale connected? Try: tailscale status")
        return False


def check_aperture_endpoint() -> bool:
    base = os.getenv("APERTURE_URL")
    if not base:
        print(f"  [{WARN}] APERTURE_URL not set — Aperture endpoint unknown")
        return True  # Not a hard failure, might be set later

    try:
        import httpx
        response = httpx.get(f"{base}/v1/models", timeout=5)
        print(f"  [{PASS}] Aperture endpoint at {base} (HTTP {response.status_code})")
        return True
    except Exception as e:
        print(f"  [{FAIL}] Aperture endpoint at {base} — {e}")
        return False


def check_env_var(name: str) -> bool:
    value = os.getenv(name)
    if value and not value.startswith("<"):
        print(f"  [{PASS}] {name} = {value}")
        return True
    else:
        print(f"  [{WARN}] {name} not set")
        return True  # Warnings, not failures


def main() -> None:
    print("\n  Workshop Environment Check")
    print("  " + "=" * 40)

    all_ok = True

    print("\n  Python & Dependencies:")
    all_ok &= check_python_version()
    all_ok &= check_import("temporalio")
    all_ok &= check_import("openai")
    all_ok &= check_import("httpx")
    all_ok &= check_import("pydantic")

    print("\n  Environment Variables:")
    check_env_var("TEMPORAL_ADDRESS")
    check_env_var("APERTURE_URL")
    check_env_var("WORKSHOP_USER_ID")

    print("\n  Connectivity:")
    all_ok &= check_temporal_connection()
    all_ok &= check_aperture_endpoint()

    print("\n  " + "=" * 40)
    if all_ok:
        print(f"  [{PASS}] All checks passed. You're ready!\n")
    else:
        print(f"  [{FAIL}] Some checks failed. See above for details.\n")
        sys.exit(1)


if __name__ == "__main__":
    main()
