import asyncio
import logging
import os
import uuid

from shared import WorkflowInput
from temporalio.client import Client
from temporalio.envconfig import ClientConfig
from workflow import GetAddressFromIP

USER_ID = os.getenv("WORKSHOP_USER_ID", "unknown")
TASK_QUEUE = f"{USER_ID}-hello-tailnet"


async def main() -> None:
    # TODO: Load the "tailnet" profile from temporal.toml so this connects
    # to the shared Temporal server instead of localhost.
    # Pass profile="tailnet" to load_client_connect_config().
    config = ClientConfig.load_client_connect_config()
    client = await Client.connect(**config)

    result = await client.execute_workflow(
        GetAddressFromIP.run,
        WorkflowInput(name=USER_ID),
        # TODO: Add your USER_ID to the Workflow ID to help identify
        # your Workflow Execution
        id=f"geo-ip-{uuid.uuid4()}",
        task_queue=TASK_QUEUE,
    )

    print(f"Your IP address: {result.ip_addr}")
    print(f"Your location:   {result.location}")


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(main())
