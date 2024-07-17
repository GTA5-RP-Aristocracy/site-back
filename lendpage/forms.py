from django import forms
from .models import User

class AddPostForm(forms.Form):
    name = forms.CharField(max_length=50, label="Ваш ник")
    email = forms.EmailField(label="Наш email")
    message = forms.CharField(widget=forms.Textarea(), label="Ваш текст сообщения")
