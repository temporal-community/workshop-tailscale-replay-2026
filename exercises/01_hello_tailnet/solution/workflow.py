# ABOUTME: Temporal Workflow that gets your public IP and geolocates it.
# ABOUTME: Demonstrates basic workflow structure with two sequential activities.

import asyncio
from datetime import timedelta

from shared import WorkflowInput, WorkflowOutput
from temporalio import workflow

with workflow.unsafe.imports_passed_through():
    from activities import get_ip, get_location_info


@workflow.defn
class GetAddressFromIP:
    @workflow.run
    async def run(self, input: WorkflowInput) -> WorkflowOutput:
        ip_address = await workflow.execute_activity(
            get_ip,
            start_to_close_timeout=timedelta(seconds=5),
        )

        await asyncio.sleep(input.seconds)

        location = await workflow.execute_activity(
            get_location_info,
            ip_address,
            start_to_close_timeout=timedelta(seconds=5),
        )

        return WorkflowOutput(ip_addr=ip_address, location=location)
