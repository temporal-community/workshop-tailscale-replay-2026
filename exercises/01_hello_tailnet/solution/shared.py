# ABOUTME: Shared data models for the Hello Tailnet exercise.
# ABOUTME: Defines workflow input/output dataclasses.

from dataclasses import dataclass


@dataclass
class WorkflowInput:
    name: str
    seconds: int = 0


@dataclass
class WorkflowOutput:
    ip_addr: str
    location: str
