from django.test import TestCase
from django.contrib.auth import get_user_model
from .serializers import UserRegistrationSerializer

from rest_framework.test import APIClient
from django.urls import reverse

User = get_user_model()

class AuthTests(TestCase):
    def setUp(self):
        self.client = APIClient()
        self.user = User.objects.create_user(username='testuser', password='testpassword', email='test@example.com')
        self.token_url = reverse('token_obtain_pair')
        self.verify_url = reverse('token_verify')

    def test_auth_flow(self):
        # 1. Get token
        response = self.client.post(self.token_url, {
            'username': 'testuser',
            'password': 'testpassword'
        })
        self.assertEqual(response.status_code, 200)
        self.assertIn('access', response.data)
        access_token = response.data['access']

        # 2. Verify token
        response = self.client.post(self.verify_url, {'token': access_token})
        self.assertEqual(response.status_code, 200)

        # 3. Verify invalid token
        response = self.client.post(self.verify_url, {'token': 'invalidtoken'})
        self.assertEqual(response.status_code, 401)

class SerializerTests(TestCase):
    def test_registration_serializer_success(self):
        data = {
            'username': 'newuser',
            'email': 'new@example.com',
            'password': 'password123',
            'confirm_password': 'password123'
        }
        serializer = UserRegistrationSerializer(data=data)
        self.assertTrue(serializer.is_valid())
        user = serializer.save()
        self.assertTrue(user.check_password('password123'))
        self.assertEqual(user.email, 'new@example.com')

    def test_registration_serializer_mismatch_password(self):
        data = {
            'username': 'newuser',
            'email': 'new@example.com',
            'password': 'password123',
            'confirm_password': 'mismatch'
        }
        serializer = UserRegistrationSerializer(data=data)
        self.assertFalse(serializer.is_valid())
        self.assertIn('password', serializer.errors)

class ViewTests(TestCase):
    def setUp(self):
        self.client = APIClient()
        self.user_data = {
            'username': 'newuser',
            'email': 'new@example.com',
            'password': 'password123',
            'confirm_password': 'password123'
        }
        self.register_url = reverse('register')
        self.profile_url = reverse('profile')
        self.user = User.objects.create_user(username='existing', email='existing@example.com', password='password123')

    def test_registration_view(self):
        response = self.client.post(self.register_url, self.user_data)
        self.assertEqual(response.status_code, 201)
        self.assertTrue(User.objects.filter(email='new@example.com').exists())

    def test_profile_view_unauthenticated(self):
        response = self.client.get(self.profile_url)
        self.assertEqual(response.status_code, 401)

    def test_profile_view_authenticated(self):
        self.client.force_authenticate(user=self.user)
        response = self.client.get(self.profile_url)
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['email'], 'existing@example.com')

    def test_profile_update(self):
        self.client.force_authenticate(user=self.user)
        update_data = {'first_name': 'Updated', 'last_name': 'Name'}
        response = self.client.patch(self.profile_url, update_data)
        self.assertEqual(response.status_code, 200)
        self.user.refresh_from_db()
        self.assertEqual(self.user.first_name, 'Updated')
