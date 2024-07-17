from django.db import models

class User(models.Model):
    name = models.TextField()
    email = models.EmailField()
    message = models.TextField(blank=True)
