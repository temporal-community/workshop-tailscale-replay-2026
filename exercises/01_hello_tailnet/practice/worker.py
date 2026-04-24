import asyncio
import concurrent.futures
import logging
import os

from activities import get_ip, get_location_info
from temporalio.client import Client
from temporalio.envconfig import ClientConfig
from temporalio.worker import Worker
from workflow import GetAddressFromIP

USER_ID = os.getenv("WORKSHOP_USER_ID", "unknown")
TASK_QUEUE = f"{USER_ID}-hello-tailnet"


async def main() -> None:
    # TODO: Load the "tailnet" profile from temporal.toml so this connects
    # to the shared Temporal server instead of localhost.
    # Pass profile="tailnet" to load_client_connect_config().
    config = ClientConfig.load_client_connect_config()
    logging.info(f"Connecting to Temporal at {config.get('target_host')}")
    client = await Client.connect(**config)

    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as activity_executor:
        worker = Worker(
            client,
            task_queue=TASK_QUEUE,
            workflows=[GetAddressFromIP],
            activities=[get_ip, get_location_info],
            activity_executor=activity_executor,
        )
        logging.info(f"Starting worker on task queue: {TASK_QUEUE}")
        await worker.run()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(main())
