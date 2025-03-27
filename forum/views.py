from django.http import JsonResponse


def home_view(request):
    if request.user.is_authenticated:
        return JsonResponse({"message": "Welcome back!", "user": request.user.username})
    else:
        return JsonResponse({"message": "Hello, guest!"})
