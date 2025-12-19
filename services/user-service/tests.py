from django.test import TestCase
from django.contrib.auth.models import User
from rest_framework.test import APIClient
from django.urls import reverse

class AuthTests(TestCase):
    def setUp(self):
        self.client = APIClient()
        self.user = User.objects.create_user(username='testuser', password='testpassword')
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
