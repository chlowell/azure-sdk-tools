from asyncio.log import logger
from multiprocessing.sharedctypes import Value


_HELP_URL_STRING = "https://github.com/Azure/azure-sdk-tools/blob/main/doc/python_apiview_errors.md#{0}"


_supported_diagnostic_codes = {
    "missing-type",
    "missing-source-link",
    "missing-kwargs",
    "name-mismatch",
    "decorator-parse",
    "missing-return-type",
    "missing-typehint",
    "return-type-mismatch",
    "list-return-type"
}


class Diagnostic:
    id_counter = 1

    def __init__(self, *, code, target_id, message):
        if code not in _supported_diagnostic_codes:
            raise ValueError(f"APIView found unexpected code '{code}'.")
        self.diagnostic_id = "AZ_PY_{}".format(Diagnostic.id_counter)
        Diagnostic.id_counter += 1
        self.text = f"[{code}] {message}"
        self.help_link_uri = _HELP_URL_STRING.format(code)
        self.target_id = target_id
