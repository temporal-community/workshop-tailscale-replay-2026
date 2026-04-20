# ABOUTME: Starts tool-calling or agent workflows on the shared Temporal server.
# ABOUTME: Use --agent flag for the agentic loop, default is tool-calling.

import argparse
import asyncio
import logging
import os
import uuid

from agent_workflow import AgentWorkflow
from temporalio.client import Client
from temporalio.contrib.pydantic import pydantic_data_converter
from temporalio.envconfig import ClientConfig
from tool_calling_workflow import ToolCallingWorkflow

USER_ID = os.getenv("WORKSHOP_USER_ID", "unknown")
TOOL_CALLING_TASK_QUEUE = f"{USER_ID}-tool-calling"
AGENT_TASK_QUEUE = f"{USER_ID}-agent"


async def run_tool_calling(query: str) -> None:
    """Start the tool-calling workflow and print the result."""
    config = ClientConfig.load_client_connect_config()
    client = await Client.connect(**config, data_converter=pydantic_data_converter)

    result = await client.execute_workflow(
        ToolCallingWorkflow,
        query,
        id=f"{USER_ID}-tool-calling-{uuid.uuid4()}",
        task_queue=TOOL_CALLING_TASK_QUEUE,
    )
    print(f"Result: {result}")


async def run_agent(query: str) -> None:
    """Start the agentic loop workflow and print the result."""
    config = ClientConfig.load_client_connect_config()
    client = await Client.connect(**config, data_converter=pydantic_data_converter)

    result = await client.execute_workflow(
        AgentWorkflow,
        query,
        id=f"{USER_ID}-agent-{uuid.uuid4()}",
        task_queue=AGENT_TASK_QUEUE,
    )
    print(f"Result: {result}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Run Weather Agent Workflow")
    parser.add_argument("--agent", action="store_true", help="Run the agentic loop")
    parser.add_argument("query", nargs="?", default=None, help="The query to send")
    args = parser.parse_args()

    query = args.query
    if not query:
        if args.agent:
            query = "What's the weather like where I am right now?"
        else:
            query = "What are the weather alerts in California?"

    print(f"Query: {query}")

    if args.agent:
        asyncio.run(run_agent(query))
    else:
        asyncio.run(run_tool_calling(query))


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    main()
