from dataclasses import dataclass


@dataclass
class WorkflowInput:
    name: str
    seconds: int = 0


@dataclass
class WorkflowOutput:
    ip_addr: str
    location: str
