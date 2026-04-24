import requests
from temporalio import activity


@activity.defn
def get_ip() -> str:
    """Get the public IP address of the machine running this activity."""
    response = requests.get("https://icanhazip.com")
    response.raise_for_status()
    return response.text.strip()


@activity.defn
def get_location_info(ip: str) -> str:
    """Get city, region, and country for an IP address."""
    response = requests.get(f"http://ip-api.com/json/{ip}")
    response.raise_for_status()
    result = response.json()
    return f"{result['city']}, {result['regionName']}, {result['country']}"
