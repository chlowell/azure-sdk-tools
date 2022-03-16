import inspect
import json
import os
from pylint import epylint
from typing import List

class PylintError:

    def __init__(self, **kwargs):
        self.type = kwargs.pop('type', None)
        self.module = kwargs.pop('module', None)
        self.obj = kwargs.pop('obj', None)
        self.line = kwargs.pop('line', None)
        self.column = kwargs.pop('column', None)
        self.end_line = kwargs.pop('endLine', None)
        self.end_column = kwargs.pop('endColumn', None)
        self.path = kwargs.pop('path', None)
        self.symbol = kwargs.pop('symbol', None)
        self.message = kwargs.pop('message', None)
        self.message_id = kwargs.pop('message-id', None)
        # TODO: remove this
        assert len(kwargs) == 0


class PylintParser:

    items: List[PylintError] = []

    @classmethod
    def parse(cls, path):
        (pylint_stdout, _) = epylint.py_run(f"{path} -f json", return_std=True)
        # strip put stray, non-json lines from stdout
        stdout_lines = [x for x in pylint_stdout.readlines() if not x.startswith("Exception")]
        json_items = json.loads("".join(stdout_lines))
        cls.items = [PylintError(**x) for x in json_items]

    @classmethod
    def get_items(cls, obj) -> List[PylintError]:
        module_name = inspect.getmodule(obj).__name__
        source_file = inspect.getsourcefile(obj)
        (source_lines, start_line) = inspect.getsourcelines(obj)
        end_line = start_line + len(source_lines) - 1

        results = []
        for item in cls.items:
            item_path = item.path
            if source_file.endswith(item_path):
                test = "Best"
            else:
                test = "Best"
        return results