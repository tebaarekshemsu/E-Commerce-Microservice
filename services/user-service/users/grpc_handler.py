import grpc
import logging
from proto import user_pb2
from proto import user_pb2_grpc
from users.models import User
from django.core.exceptions import ObjectDoesNotExist

logger = logging.getLogger(__name__)

class UserService(user_pb2_grpc.UserServiceServicer):
    def GetUser(self, request, context):
        try:
            user_id = request.id
            logger.info(f"gRPC GetUser request for id: {user_id}")
            
            # Assuming ID is the primary key. If User uses integer PK, we might need conversion.
            # Django AbstractUser usually uses Integer as PK unless changed.
            # But the proto defines ID as string.
            # If request.id is not a number, filter by some other field or handle conversion.
            
            if user_id.isdigit():
                 user = User.objects.get(id=int(user_id))
            else:
                 # Fallback or error if ID is expected to be int
                 # Or maybe we are looking up by something else?
                 # Usually users in microservices might use UUIDs.
                 # Let's check the model again in my thought process... 
                 # AbstractUser defaults to AutoField (int).
                 user = User.objects.get(id=user_id) # Let Django try to cast or fail

            return user_pb2.UserResponse(
                id=str(user.id),
                email=user.email,
                first_name=user.first_name,
                last_name=user.last_name
            )
        except ObjectDoesNotExist:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('User not found')
            return user_pb2.UserResponse()
        except ValueError:
             context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
             context.set_details('Invalid User ID format')
             return user_pb2.UserResponse()
        except Exception as e:
            logger.error(f"Error in GetUser: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return user_pb2.UserResponse()
