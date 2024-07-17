from django.urls import path
from lendpage import views

urlpatterns = [
    path('', views.display_form),
]
