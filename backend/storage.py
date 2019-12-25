import grpc
import time
from concurrent import futures

import api_pb2_grpc
import api_pb2

google_prods = {'compute': ['Compute Engine', 'App Engine', 'Cloud Functions', 'Cloud Run', 'GKE'], 'database': ['Cloud Bigtable', 'Cloud Spanner', 'Cloud SQL', 'Cloud Firestore', 'Cloud Memorystore', 'BigQuery']}
aws_prods = {'compute': ['Amazon EC2', 'Amazon Lightsail', 'Amazon ECS', 'Amazon EKS', 'AWS Fargate', 'AWS Lambda'], 'database': ['Amazon Aurora', 'Amazon RDS ', 'Amazon Redshift', 'Amazon DynamoDB', 'Amazon ElastiCache for Memcached', 'Amazon ElastiCache for Redis', 'Amazon Neptune']}
oracle_prods = {'compute': ['Bare Metal Compute', 'Container Engine for Kubernetes', 'Virtual Machines', 'Oracle Functions'], 'database': ['Autonomous Data Warehouse', 'Autonomous Transaction Processing', 'Database Cloud Service: Virtual Machine', 'Exadata Cloud Service', 'NoSQL Database']}

vendors = {'google': google_prods, 'aws': aws_prods, 'oracle': oracle_prods}

class BackendService(api_pb2_grpc.BackendServiceServicer):      

    def RetrieveItems(self, request, context):
        
        # Retrieve vendor and prodType from client
        vendor = request.vendor.lower()
        product_type = request.prodType.lower()

        prods = get_prods(vendor, product_type)

        response = api_pb2.ApiResponse()
        response.prods.extend(prods)

        return response

def get_prods(vendor, product_type):

    if vendor in vendors.keys():
        if product_type in vendors[vendor].keys():
            return vendors[vendor][product_type]
        else:
            return "No Product Type"
    else:
        return "No valid Vendor"

def serve(port):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    api_pb2_grpc.add_BackendServiceServicer_to_server(BackendService(), server)
    server.add_insecure_port('[::]:' + str(port))
    server.start()
    print("Listening on port {}..".format(port))
    try:
        while True:
            time.sleep(10000)
    except KeyboardInterrupt:
        server.stop(0)

if __name__== "__main__":
    serve(6000)


