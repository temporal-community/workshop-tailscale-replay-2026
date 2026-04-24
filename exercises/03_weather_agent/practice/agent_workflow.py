import json
from datetime import timedelta

from temporalio import workflow

with workflow.unsafe.imports_passed_through():
    from activities import create
    from models import HELPFUL_AGENT_SYSTEM_INSTRUCTIONS, OpenAIResponsesRequest
    from tools import get_tools


@workflow.defn
class AgentWorkflow:
    @workflow.run
    async def run(self, input: str) -> str:
        input_list = [{"type": "message", "role": "user", "content": input}]

        # TODO 2: Enable the agentic loop so the LLM can make multiple tool calls.
        # With the condition set to False, the body below never runs and the
        # Workflow returns None. Flip False to True and the LLM can iterate
        # until it has enough information to respond.
        while False:
            workflow.logger.info("=" * 80)

            # Consult the LLM with all available tools
            result = await workflow.execute_activity(
                create,
                OpenAIResponsesRequest(
                    model="gpt-4o-mini",
                    instructions=HELPFUL_AGENT_SYSTEM_INSTRUCTIONS,
                    input=input_list,
                    tools=get_tools(),
                ),
                start_to_close_timeout=timedelta(seconds=30),
            )

            item = result.output[0]

            if item.type == "function_call":
                result = await self._handle_function_call(item, result, input_list)

                input_list.append({
                    "type": "function_call_output",
                    "call_id": item.call_id,
                    "output": result,
                })
            else:
                workflow.logger.info(f"No tools chosen, responding with a message: {result.output_text}")
                return result.output_text

    async def _handle_function_call(self, item, result, input_list):
        """Execute a tool call chosen by the LLM using dynamic activities."""
        serialized_item = result.output[0]
        input_list += [
            serialized_item.model_dump() if hasattr(serialized_item, "model_dump") else serialized_item
        ]

        args = json.loads(item.arguments) if isinstance(item.arguments, str) else item.arguments

        # TODO 3: Execute the tool chosen by the LLM as a Temporal activity.
        # The LLM picked a tool (item.name) and provided arguments (args).
        # Use workflow.execute_activity to run it as a dynamic activity.
        # Hint: await workflow.execute_activity(item.name, args, start_to_close_timeout=timedelta(seconds=30))
        tool_result = ""

        workflow.logger.info(f"Made a tool call to {item.name}")
        return tool_result
