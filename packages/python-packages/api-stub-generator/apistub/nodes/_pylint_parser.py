import inspect
import json
import os
from pylint import epylint
from typing import List

class PylintError:

    def __init__(self, pkg_name, **kwargs):
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
        if self.path.startswith(pkg_name):
            self.path = self.path[(len(f"{pkg_name}\\\\") - 1):]
        code = self.symbol[0]
        self.level = "ERROR" if code in "EF" else "WARNING"

    def generate_tokens(self, apiview, line_id):
        text = f"{self.message} [{self.symbol}]"
        apiview.add_diagnostic(text, line_id)


class PylintParser:

    ALLOWED_PYLINT_SYMBOLS = {
        "connection-string-should-not-be-constructor-param",
        "missing-client-constructor-parameter-credential",
        "missing-client-constructor-parameter-kwargs",
        "client-method-has-more-than-5-positional-arguments",
        "client-method-missing-type-annotations",
        "incorrect-naming-convention",
        "client-method-missing-kwargs",
        "config-missing-kwargs-in-policy",
        "async-client-bad-name",
        "file-needs-copyright-header",
        "client-method-name-no-double-underscore",
        "specify-parameter-names-in-call",
        "package-name-incorrect",
        "client-suffix-needed",
        "docstring-admonition-needs-newline"
    }

    items: List[PylintError] = []

    @classmethod
    def parse(cls, path):
        pkg_name = os.path.split(path)[-1]
        (pylint_stdout, _) = epylint.py_run(f"{path} -f json", return_std=True)
        # strip put stray, non-json lines from stdout
        stdout_lines = [x for x in pylint_stdout.readlines() if not x.startswith("Exception")]
        json_items = json.loads("".join(stdout_lines))
        cls.items = [PylintError(pkg_name, **x) for x in json_items]

    @classmethod
    def get_items(cls, obj) -> List[PylintError]:
        results = []
        try:
            source_file = inspect.getsourcefile(obj)
            (source_lines, start_line) = inspect.getsourcelines(obj)
            end_line = start_line + len(source_lines) - 1
        except:
            return results

        for item in cls.items:
            item_path = item.path
            if source_file.endswith(item_path):
                # only include linter warnings associated with the first line of a code block
                if item.line == start_line:
                    results.append(item)
        return results