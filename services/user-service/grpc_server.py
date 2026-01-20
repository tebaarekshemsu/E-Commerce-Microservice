import grpc
from concurrent import futures
import logging
import os
import sys
import django
from django.conf import settings

# Add the project directory to Python path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

# Setup Django
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
django.setup()

from users.models import User
from proto import user_pb2, user_pb2_grpc

class UserServiceServicer(user_pb2_grpc.UserServiceServicer):
    def GetUser(self, request, context):
        try:
            logging.info(f"gRPC GetUser called with ID: {request.id}")
            
            # Get user from database
            user = User.objects.get(id=request.id)
            
            # Create response
            response = user_pb2.UserResponse(
                id=str(user.id),
                email=user.email,
                first_name=user.first_name,
                last_name=user.last_name
            )
            
            logging.info(f"✅ Retrieved user {request.id} via gRPC")
            return response
            
        except User.DoesNotExist:
            logging.error(f"User {request.id} not found")
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f'User with id {request.id} not found')
            return user_pb2.UserResponse()
        except Exception as e:
            logging.error(f"Error getting user: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f'Internal server error: {str(e)}')
            return user_pb2.UserResponse()

def serve():
    port = os.getenv('GRPC_PORT', '9003')
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    user_pb2_grpc.add_UserServiceServicer_to_server(UserServiceServicer(), server)
    
    listen_addr = f'[::]:{port}'
    server.add_insecure_port(listen_addr)
    
    logging.info(f"🚀 User gRPC server starting on port {port}")
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logging.info("Shutting down gRPC server...")
        server.stop(0)

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    serve()