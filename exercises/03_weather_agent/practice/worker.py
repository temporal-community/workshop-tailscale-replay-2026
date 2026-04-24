import argparse
import asyncio
import concurrent.futures
import logging
import os
import warnings

from activities import create, dynamic_tool_activity, get_weather_alerts
from agent_workflow import AgentWorkflow
from temporalio.client import Client
from temporalio.contrib.pydantic import pydantic_data_converter
from temporalio.envconfig import ClientConfig
from temporalio.worker import Worker
from tool_calling_workflow import ToolCallingWorkflow

USER_ID = os.getenv("WORKSHOP_USER_ID")
if not USER_ID:
    raise SystemExit(
        "WORKSHOP_USER_ID is not set. Open a new terminal or run `source ~/.bashrc`. "
        "Instruqt sets this automatically for all workshop shells."
    )
TOOL_CALLING_TASK_QUEUE = f"{USER_ID}-tool-calling"
AGENT_TASK_QUEUE = f"{USER_ID}-agent"


async def run_tool_calling_worker() -> None:
    """Run the worker for the simple tool-calling workflow."""
    config = ClientConfig.load_client_connect_config(profile="tailnet")
    logging.info(f"Connecting to Temporal at {config.get('target_host')}")
    client = await Client.connect(**config, data_converter=pydantic_data_converter)

    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as activity_executor:
        worker = Worker(
            client,
            task_queue=TOOL_CALLING_TASK_QUEUE,
            workflows=[ToolCallingWorkflow],
            activities=[create, get_weather_alerts],
            activity_executor=activity_executor,
        )
        logging.info(f"Starting tool-calling worker on task queue: {TOOL_CALLING_TASK_QUEUE}")
        await worker.run()


async def run_agent_worker() -> None:
    """Run the worker for the agentic loop workflow with dynamic activities."""
    config = ClientConfig.load_client_connect_config(profile="tailnet")
    logging.info(f"Connecting to Temporal at {config.get('target_host')}")
    client = await Client.connect(**config, data_converter=pydantic_data_converter)

    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as activity_executor:
        worker = Worker(
            client,
            task_queue=AGENT_TASK_QUEUE,
            workflows=[AgentWorkflow],
            activities=[create, dynamic_tool_activity],
            activity_executor=activity_executor,
        )
        logging.info(f"Starting agent worker on task queue: {AGENT_TASK_QUEUE}")
        await worker.run()


def main() -> None:
    logging.basicConfig(level=logging.INFO)
    logging.getLogger("temporalio").setLevel(logging.WARNING)
    logging.getLogger("temporalio.workflow").setLevel(logging.INFO)
    logging.getLogger("openai").setLevel(logging.WARNING)
    logging.getLogger("httpx").setLevel(logging.WARNING)
    warnings.filterwarnings("ignore", category=UserWarning, module="temporalio.converter")

    parser = argparse.ArgumentParser(description="Run Weather Agent Worker")
    parser.add_argument("--agent", action="store_true", help="Run the agentic loop worker")
    args = parser.parse_args()

    if args.agent:
        asyncio.run(run_agent_worker())
    else:
        asyncio.run(run_tool_calling_worker())


if __name__ == "__main__":
    main()
