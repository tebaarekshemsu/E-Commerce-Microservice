from django.core.management.base import BaseCommand
import grpc
from concurrent import futures
import logging
import time
import os
from proto import user_pb2_grpc
from users.grpc_handler import UserService

class Command(BaseCommand):
    help = 'Runs the gRPC server'

    def handle(self, *args, **options):
        logger = logging.getLogger(__name__)
        port = os.environ.get('GRPC_PORT', '50051')
        
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        user_pb2_grpc.add_UserServiceServicer_to_server(UserService(), server)
        server.add_insecure_port(f'[::]:{port}')
        
        self.stdout.write(self.style.SUCCESS(f'Starting User gRPC server on port {port}...'))
        server.start()
        
        try:
            while True:
                time.sleep(86400)
        except KeyboardInterrupt:
            server.stop(0)
