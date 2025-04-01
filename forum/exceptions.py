from rest_framework.views import exception_handler


def bad_request_handler(exc, context):
    response = exception_handler(exc, context)

    if response is not None and isinstance(response.data, dict):
        response.data = {
            key: value[0] if isinstance(value, list) and len(value) == 1 else value
            for key, value in response.data.items()
        }

    return response
