import grpc
import time
from concurrent import futures
import logging

import api_pb2_grpc
import api_pb2

g_compute = [{'title': 'Compute Engine', 'url': 'https://cloud.google.com/compute/'}, {'title': 'App Engine', 'url': 'https://cloud.google.com/appengine/'}, {'title': 'Cloud Functions', 'url': 'https://cloud.google.com/functions/'}, {'title': 'Cloud Run', 'url': 'https://cloud.google.com/run/'}, {'title': 'GKE', 'url': 'https://cloud.google.com/kubernetes-engine/'}]
g_storage = [{'title': 'Cloud Bigtable' , 'url': 'https://cloud.google.com/bigtable/'}, {'title': 'Cloud Spanner', 'url': 'https://cloud.google.com/spanner/'}, {'title': 'Cloud SQL', 'url': 'https://cloud.google.com/sql/'}, {'title': 'Cloud Firestore', 'url': 'https://cloud.google.com/firestore/'},{'title': 'Cloud Memorystore', 'url': 'https://cloud.google.com/memorystore/'},{'title': 'BigQuery', 'url': 'https://cloud.google.com/bigquery/'}]
a_compute = [{'title': 'Amazon EC2', 'url': 'https://aws.amazon.com/ec2/'}, {'title': 'Amazon Lightsail', 'url': 'https://aws.amazon.com/lightsail/'}, {'title': 'Amazon ECS', 'url': 'https://aws.amazon.com/ecs/'}, {'title': 'Amazon EKS', 'url': 'https://aws.amazon.com/eks/'}, {'title': 'AWS Fargate', 'url': 'https://aws.amazon.com/fargate/'},{'title': 'AWS Lambda', 'url': 'https://aws.amazon.com/lambda/'}]
a_storage = [{'title': 'Amazon Aurora', 'url': 'https://aws.amazon.com/rds/aurora/'}, {'title': 'Amazon RDS', 'url': 'https://aws.amazon.com/rds/'}, {'title': 'Amazon Redshift', 'url': 'https://aws.amazon.com/redshift/'}, {'title': 'Amazon DynamoDB', 'url': 'https://aws.amazon.com/dynamodb/'},{'title': 'Amazon ElastiCache for Memcached', 'url': 'https://aws.amazon.com/elasticache/memcached/'}, {'title': 'Amazon ElastiCache for Redis', 'url': 'https://aws.amazon.com/elasticache/redis/'}, {'title': 'Amazon Neptune', 'url': 'https://aws.amazon.com/neptune/'}]
o_compute = [{'title': 'Bare Metal Compute', 'url': 'https://www.oracle.com/cloud/compute/bare-metal.html'}, {'title': 'Container Engine for Kubernetes', 'url': 'https://www.oracle.com/cloud/compute/container-engine-kubernetes.html'}, {'title': 'Virtual Machines', 'url': 'https://www.oracle.com/cloud/compute/virtual-machines.html'}, {'title': 'Oracle Functions', 'url': 'https://www.oracle.com/ro/cloud/cloud-native/functions/'}]
o_storage = [{'title': 'Autonomous Data Warehouse', 'url': 'https://www.oracle.com/database/adw-cloud.html'}, {'title': 'Autonomous Transaction Processing', 'url': 'https://www.oracle.com/database/atp-cloud.html'}, {'title': 'Database Cloud Service: Virtual Machine', 'url': 'https://www.oracle.com/database/vm-cloud.html'}, {'title': 'Exadata Cloud Service', 'url': 'https://www.oracle.com/database/exadata-cloud-service.html'},{'title': 'NoSQL Database', 'url': 'https://www.oracle.com/database/nosql-cloud.html'}]

vendors = {'google': {'compute': g_compute, 'storage': g_storage}, 'aws': {'compute': a_compute, 'storage': a_storage}, 'oracle': {'compute': o_compute, 'storage': o_storage}}

class Storage(api_pb2_grpc.StorageServiceServicer):      

    def GetProdsDetail(self, request, context):
        
        # Retrieve vendor and prodType from client
        vendor = request.vendor.lower()
        product_type = request.productType.lower()

        prod_type_list = get_prods(vendor, product_type)

        products = []
        
        for prod in prod_type_list:
            product = api_pb2.Product()
            product.title = prod["title"]
            product.url = prod["url"]
            products.append(product)


        return api_pb2.StorageResponse(prodDetail=products)

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
    api_pb2_grpc.add_StorageServiceServicer_to_server(Storage(), server)
    server.add_insecure_port('[::]:' + str(port))
    server.start()
    print("Listening on port {}..".format(port))
    try:
        while True:
            time.sleep(10000)
    except KeyboardInterrupt:
        server.stop(0)

if __name__== "__main__":
    logging.basicConfig()
    serve(6000)
