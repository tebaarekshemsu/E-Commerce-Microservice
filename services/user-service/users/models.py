from django.db import models
from django.contrib.auth.models import AbstractUser

class User(AbstractUser):
    class Role(models.TextChoices):
        ADMIN = 'ADMIN', 'Admin'
        USER = 'USER', 'User'

    email = models.EmailField(unique=True)
    phone = models.CharField(max_length=50, blank=True, null=True)
    image_url = models.CharField(max_length=255, blank=True, null=True)
    role = models.CharField(
        max_length=20,
        choices=Role.choices,
        default=Role.USER
    )

    REQUIRED_FIELDS = ['email', 'first_name', 'last_name']
    
    def __str__(self):
        return self.username

class Address(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='addresses')
    full_address = models.CharField(max_length=255)
    postal_code = models.CharField(max_length=20, blank=True, null=True)
    city = models.CharField(max_length=100)

    def __str__(self):
        return f"{self.full_address}, {self.city}"
