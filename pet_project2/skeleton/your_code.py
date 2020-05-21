"""
Classes that you need to complete.
"""

# Optional import
from serialization import jsonpickle
import credential
import petrelic
import hashlib
import random as rd
import time

class Server:
    """Server"""
    
    @staticmethod
    def generate_ca(valid_attributes):
        """Initializes the credential system. Runs exactly once in the
        beginning. Decides on schemes public parameters and choses a secret key
        for the server.

        Args:
            valid_attributes (string): a list of all valid attributes. Users cannot
            get a credential with a attribute which is not included here.

            Note: You can use JSON to encode valid_attributes in the string.

        Returns:
            (tuple): tuple containing:
                byte[] : server's pubic information
                byte[] : server's secret key
            You are free to design this as you see fit, but all commuincations
            needs to be encoded as byte arrays.
        """

        Issuer = credential.Issuer()
        Issuer.setup(valid_attributes.split(","))
        pub = Issuer.get_serialized_public_key()
        sec = Issuer.get_serialized_secret_key()

        print("CA generated")
        return pub, sec


    @staticmethod
    def register(server_sk, issuance_request, username, attributes):
        """ Registers a new account on the server.

        Args:
            server_sk (byte []): the server's secret key (serialized) and other parameters
            issuance_request (bytes[]): The issuance request (serialized)
            username (string): username
            attributes (string): attributes

            Note: You can use JSON to encode attributes in the string.

        Return:
            response (bytes[]): the client should be able to build a credential
            with this response.
        """
        # time to compute elapsed communication time
        #print(time.time_ns())

        if(attributes):
            attributes_list = attributes.split(',')
            pub_attr_len = len(attributes_list)
        else:
            attributes_list = []
            pub_attr_len = 0
        server_sk = credential.Signature.deserialize(server_sk)
        server_pb_params = server_sk.get("public_params")
        issuance_request = credential.Signature.deserialize(issuance_request)

        # Zero Knowledge Proof verification on client commit
        C = issuance_request.get("C")
        r = issuance_request.get("r")
        gamma = issuance_request.get("gamma")
        g = server_pb_params.get("g")
        Y = server_pb_params.get("pk")

        to_hash = C.to_binary() + g.to_binary()
        for i in range(len(gamma)-1):
            to_hash += Y[i+pub_attr_len].to_binary()
        for gammai in gamma:
            to_hash += gammai.to_binary()

        c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)
        gamma_prod = server_pb_params.get("G1").neutral_element()
        for gammai in gamma:
            gamma_prod *= gammai

        check_prod = (C**((-c_hash) % server_pb_params.get("G1").order()))*(g**r[0])
        for i in range(len(r)-1):
            check_prod *= Y[i+pub_attr_len]**r[i+1]
        
        proof_correct = (gamma_prod == check_prod)

        # if proof incorrect, return empty byte array
        if not proof_correct:
            print("ZKP incorrect: cancel registration")
            return bytearray()
        # if proof correct, return issuer credentials
        else:
            print("ZKP correct: proceed registration")
            sigma = credential.Issuer.issue(issuance_request.get("C"), username, attributes_list, server_sk)
            return bytearray(jsonpickle.encode(sigma), 'utf-8')


    def check_request_signature(self, server_pk, message, revealed_attributes_str, signature):
        """
        Args:
            server_pk (byte[]): the server's public key (serialized) and other parameters
            message (byte[]): The message to sign
            revealed_attributes (string): revealed attributes
            signature (bytes[]): user's autorization (serialized)

            Note: You can use JSON to encode revealed_attributes in the string.

        Returns:
            valid (boolean): is signature valid
        """
        # time to compute elapsed communication time
        #print(time.time_ns())

        if(revealed_attributes_str):
            revealed_attributes = revealed_attributes.split(",")
        else:
            revealed_attributes = []
        server_pk = credential.Signature.deserialize(server_pk)
        signature = credential.Signature.deserialize(signature)
        
        print("Signature verified")
        return signature.verify(server_pk, revealed_attributes, message)


class Client:
    """Client"""

    def prepare_registration(self, server_pk, username, attributes):
        """Prepare a request to register a new account on the server.

        Args:
            server_pk (byte[]): a server's public key (serialized) and other parameters
            username (string): username
            attributes (string): user's attributes

            Note: You can use JSON to encode attributes in the string.

        Return:
            tuple:
                byte[]: an issuance request
                (private_state): You can use state to store and transfer information
                from prepare_registration to proceed_registration_response.
                You need to design the state yourself.
        """

        server_pk = credential.Signature.deserialize(server_pk)
        g = server_pk.get("g")
        Y = server_pk.get("pk")
        G1 = server_pk.get("G1")
        t = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        postal = 1028
        tel = 791234567
        social_sec = 178051120
        sk = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        self.private_attr = [postal, tel, social_sec, sk]
        if(attributes):
            attributes_list = attributes.split(",")
            pub_attr_len = len(attributes_list)
        else:
            pub_attr_len = 0
        AnonCredential = credential.AnonCredential()
        C, gamma, r = AnonCredential.create_issue_request(self.private_attr, G1, Y, g, t, pub_attr_len)
        request = (jsonpickle.encode({"C": C, "gamma": gamma, "r": r})).encode('utf-8')
        private_state = {"C": C, "t": t}
        
        print("Registration request created")
        # time to compute elapsed communication time
        #print(time.time_ns())
        return (request, private_state) 


    def proceed_registration_response(self, server_pk, server_response, private_state):
        """Process the response from the server.

        Args:
            server_pk (byte[]): a server's public key (serialized) and other parameters
            server_response (byte[]): the response from the server (serialized)
            private_state (private_state): state from the prepare_registration
            request corresponding to this response

        Return:
            credential (byte []): create an attribute-based credential for the user
        """

        AnonCredential = credential.AnonCredential()
        sigma_prime = credential.Signature.deserialize(server_response)
        sigma = AnonCredential.receive_issue_response(sigma_prime, private_state.get("t"))
        return bytearray(jsonpickle.encode({"sigma": sigma, "private_attr": self.private_attr}), 'utf-8')


    def sign_request(self, server_pk, cred, message, revealed_info_str):
        """Signs the request with the clients credential.
        
        Arg:
            server_pk (byte[]): a server's public key (serialized) and other parameters
            cred (byte[]): client's credential (serialized)
            message (byte[]): message to sign
            revealed_info (string): attributes which need to be authorized

            Note: You can use JSON to encode revealed_info.

        Returns:
            byte []: message's signature (serialized)
        """
        
        cred = credential.Signature.deserialize(cred)
        sigma = cred.get("sigma")
        private_attr = cred.get("private_attr")
        server_pk = credential.Signature.deserialize(server_pk)
        if revealed_info_str:
            revealed_info = revealed_info_str.split(",")
        else:
            revealed_info = []
        
        signature = credential.Signature()
        signature.create_sign_request(server_pk, sigma, message, revealed_info, private_attr)
        
        print("Sign request sent")
        # time to compute elapsed communication time
        #print(time.time_ns())
        return signature.serialize()
        
