from enum import Enum

_HELP_URL_STRING = "https://github.com/Azure/azure-sdk-tools/blob/main/doc/python_apiview_errors.md#{0}"


class DiagnosticLevel(int, Enum):
    DEFAULT = 0
    INFO = 1
    WARNING = 2
    ERROR = 3

    def __str__(self):
        if self == 0:
            return "DEFAULT"
        elif self == 1:
            return "INFO"
        elif self == 2:
            return "WARNING"
        elif self == 3:
            return "ERROR"


class Diagnostic:
    id_counter = 1

    def __init__(self, *, symbol: str, target_id: str, message: str, level: DiagnosticLevel):
        self.diagnostic_id = "AZ_PY_{}".format(Diagnostic.id_counter)
        Diagnostic.id_counter += 1
        self.text = f"{message} [{symbol}]"
        self.help_link_uri = _HELP_URL_STRING.format(symbol)
        self.target_id = target_id
        self.level = level

    def log(self, log_func):
        log_func(f"{str(self.level)}: {self.target_id}: {self.text}")
