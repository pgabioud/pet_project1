
# importing the required modules 
import timeit
import numpy

def generate_ca_time():
    SETUP_CODE_1 = ''' 
from your_code import Server
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
server = Server()'''

    SETUP_CODE_2 = ''' 
from your_code import Server
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
server = Server()'''

    SETUP_CODE_3 = ''' 
from your_code import Server
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
server = Server()'''

    TEST_CODE = ''' 
server.generate_ca(valid_attr_str)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Generate ca MIN time 1            : {}'.format(min(times1)/10))
    print('Generate ca MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    # std of 10*X == 10*std of X
    print('Generate ca STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Generate ca MIN time 2            : {}'.format(min(times2)/10))
    print('Generate ca MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Generate ca STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Generate ca MIN time 3            : {}'.format(min(times3)/10))
    print('Generate ca MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Generate ca STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))
 

def prepare_registration_time():
    SETUP_CODE_1 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
revealed_attr_str = "1,1"
username = "oss117"
server = Server()
client = Client()
server_pk, _ = server.generate_ca(valid_attr_str)'''

    SETUP_CODE_2 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(7):
    revealed_attr_str += (",1,0")
print(len(revealed_attr_str.split(',')))
server = Server()
client = Client()
server_pk, _ = server.generate_ca(valid_attr_str)'''

    SETUP_CODE_3 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(22):
    revealed_attr_str += (",1,0")
print(len(revealed_attr_str.split(',')))
server = Server()
client = Client()
server_pk, _ = server.generate_ca(valid_attr_str)'''

    TEST_CODE = '''
client.prepare_registration(server_pk, username, revealed_attr_str)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Prepare registration MIN time 1            : {}'.format(min(times1)/10))
    print('Prepare registration MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    print('Prepare registration STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Prepare registration MIN time 2            : {}'.format(min(times2)/10))
    print('Prepare registration MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Prepare registration STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Prepare registration MIN time 3            : {}'.format(min(times3)/10))
    print('Prepare registration MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Prepare registration STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))


def register_time():
    SETUP_CODE_1 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
revealed_attr_str = "1,1"
username = "oss117"
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)'''

    SETUP_CODE_2 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(7):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)'''

    SETUP_CODE_3 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(22):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)'''

    TEST_CODE = '''
server.register(server_sk, request, username, revealed_attr_str)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Register MIN time 1            : {}'.format(min(times1)/10))
    print('Register MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    print('Register STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Register MIN time 2            : {}'.format(min(times2)/10))
    print('Register MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Register STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Register MIN time 3            : {}'.format(min(times3)/10))
    print('Register MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Register STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))


def proceed_registration_response_time():
    SETUP_CODE_1 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
revealed_attr_str = "1,1"
username = "oss117"
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)'''

    SETUP_CODE_2 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(7):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)'''

    SETUP_CODE_3 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(22):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)'''

    TEST_CODE = '''
client.proceed_registration_response(server_pk, response, prvt_state)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Proceed registration response MIN time 1            : {}'.format(min(times1)/10))
    print('Proceed registration response MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    print('Proceed registration response STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Proceed registration response MIN time 2            : {}'.format(min(times2)/10))
    print('Proceed registration response MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Proceed registration response STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Proceed registration response MIN time 3            : {}'.format(min(times3)/10))
    print('Proceed registration response MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Proceed registration response STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))


def sign_request_time():
    SETUP_CODE_1 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
revealed_attr_str = "1,1"
username = "oss117"
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
message = ("i really enjoy coding the perf tests").encode()'''

    SETUP_CODE_2 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(7):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
message = ("i really enjoy coding the perf tests").encode()'''

    SETUP_CODE_3 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(22):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
message = ("i really enjoy coding the perf tests").encode()'''

    TEST_CODE = '''
client.sign_request(server_pk, sigma, message, revealed_attr_str)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Sign request MIN time 1            : {}'.format(min(times1)/10))
    print('Sign request MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    print('Sign request STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Sign request MIN time 2            : {}'.format(min(times2)/10))
    print('Sign request MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Sign request STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Sign request MIN time 3            : {}'.format(min(times3)/10))
    print('Sign request MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Sign request STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))



def check_request_signature_time():
    SETUP_CODE_1 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
revealed_attr_str = "1,1"
username = "oss117"
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
message = ("i really enjoy coding the perf tests").encode()
sign = client.sign_request(server_pk, sigma, message, revealed_attr_str)'''

    SETUP_CODE_2 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(7):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
message = ("i really enjoy coding the perf tests").encode()
sign = client.sign_request(server_pk, sigma, message, revealed_attr_str)'''

    SETUP_CODE_3 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(22):
    revealed_attr_str += (",1,0")
server = Server()
client = Client()
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
message = ("i really enjoy coding the perf tests").encode()
sign = client.sign_request(server_pk, sigma, message, revealed_attr_str)'''

    TEST_CODE = '''
server.check_request_signature(server_pk, message, revealed_attr_str, sign)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Check request signature MIN time 1            : {}'.format(min(times1)/10))
    print('Check request signature MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    print('Check request signature STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Check request signature MIN time 2            : {}'.format(min(times2)/10))
    print('Check request signature MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Check request signature STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Check request signature MIN time 3            : {}'.format(min(times3)/10))
    print('Check request signature MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Check request signature STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))


def complete_process_time():
    SETUP_CODE_1 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "major,gender,postal,tel,social_sec,sk"
revealed_attr_str = "1,1"
username = "oss117"
message = ("i really enjoy coding the perf tests").encode()'''

    SETUP_CODE_2 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(19):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(7):
    revealed_attr_str += (",1,0")
message = ("i really enjoy coding the perf tests").encode()'''

    SETUP_CODE_3 = ''' 
from your_code import Server, Client
import credential
valid_attr_str = "0"
for i in range(49):
    valid_attr_str += ("," + str(i+1))
username = "oss117"
revealed_attr_str = "1,0"
for _ in range(22):
    revealed_attr_str += (",1,0")
message = ("i really enjoy coding the perf tests").encode()'''

    TEST_CODE = '''
server = Server()
client = Client()    
server_pk, server_sk = server.generate_ca(valid_attr_str)
request, prvt_state = client.prepare_registration(server_pk, username, revealed_attr_str)
response = server.register(server_sk, request, username, revealed_attr_str)
sigma = client.proceed_registration_response(server_pk, response, prvt_state)
sign = client.sign_request(server_pk, sigma, message, revealed_attr_str)
server.check_request_signature(server_pk, message, revealed_attr_str, sign)'''

    times1 = timeit.repeat(setup = SETUP_CODE_1, stmt = TEST_CODE, repeat = 50, number = 10) 
    times2 = timeit.repeat(setup = SETUP_CODE_2, stmt = TEST_CODE, repeat = 50, number = 10)
    times3 = timeit.repeat(setup = SETUP_CODE_3, stmt = TEST_CODE, repeat = 50, number = 10)

    print("")
    print('Complete process MIN time 1                      : {}'.format(min(times1)/10))
    print('Complete process signature MEAN time 1           : {}'.format(numpy.mean(times1)/10))
    print('Complete process signature STANDARD ERROR time 1 : {}\n'.format(numpy.nanstd(times1)/10))

    print('Complete process signature MIN time 2            : {}'.format(min(times2)/10))
    print('Complete process signature MEAN time 2           : {}'.format(numpy.mean(times2)/10))
    print('Complete process signature STANDARD ERROR time 2 : {}\n'.format(numpy.nanstd(times2)/10))

    print('Complete process signature MIN time 3            : {}'.format(min(times3)/10))
    print('Complete process signature MEAN time 3           : {}'.format(numpy.mean(times3)/10))
    print('Complete process signature STANDARD ERROR time 3 : {}\n'.format(numpy.nanstd(times3)/10))



if __name__ == "__main__": 
    #generate_ca_time()
    prepare_registration_time()
    #register_time()
    #proceed_registration_response_time()
    #sign_request_time()
    #check_request_signature_time()
    #complete_process_time()
