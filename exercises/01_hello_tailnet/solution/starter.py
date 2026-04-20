# ABOUTME: Starts the GetAddressFromIP workflow on the shared Temporal server.
# ABOUTME: Uses your WORKSHOP_USER_ID to create a unique workflow ID.

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
    config = ClientConfig.load_client_connect_config()
    client = await Client.connect(**config)

    result = await client.execute_workflow(
        GetAddressFromIP.run,
        WorkflowInput(name=USER_ID),
        id=f"{USER_ID}-geo-ip-{uuid.uuid4()}",
        task_queue=TASK_QUEUE,
    )

    print(f"Your IP address: {result.ip_addr}")
    print(f"Your location:   {result.location}")


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(main())
